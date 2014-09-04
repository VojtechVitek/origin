package master

import (
	"net/http"
	"os"
	"path"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/apiserver"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/controller"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/kubelet"
	kconfig "github.com/GoogleCloudPlatform/kubernetes/pkg/kubelet/config"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/master"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/proxy"
	pconfig "github.com/GoogleCloudPlatform/kubernetes/pkg/proxy/config"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/tools"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/GoogleCloudPlatform/kubernetes/plugin/pkg/scheduler"
	"github.com/GoogleCloudPlatform/kubernetes/plugin/pkg/scheduler/factory"
	etcdconfig "github.com/coreos/etcd/config"
	"github.com/coreos/etcd/etcd"
	etcdclient "github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"
	"github.com/google/cadvisor/client"
	"github.com/openshift/origin/pkg/api"
	"github.com/openshift/origin/pkg/service"
	"github.com/openshift/origin/pkg/template"
	"github.com/spf13/cobra"
)

func NewCommandStartAllInOne(name string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: "Launch in all-in-one mode",
		Run: func(c *cobra.Command, args []string) {
			startAllInOne()
		},
	}
}

func startAllInOne() {
	dockerAddr := getDockerEndpoint("")
	dockerClient, err := docker.NewClient(dockerAddr)
	if err != nil {
		glog.Fatal("Couldn't connect to docker.")
	}
	if err := dockerClient.Ping(); err != nil {
		glog.Errorf("WARNING: Docker could not be reached at %s.  Docker must be installed and running to start containers.\n%v", dockerAddr, err)
	}

	cadvisorClient, err := cadvisor.NewClient("http://127.0.0.1:4194")
	if err != nil {
		glog.Errorf("Error on creating cadvisor client: %v", err)
	}

	// initialize etcd
	etcdAddr := "127.0.0.1:4001"
	etcdServers := []string{} // default
	etcdConfig := etcdconfig.New()
	etcdConfig.BindAddr = etcdAddr
	etcdConfig.DataDir = "openshift.local.etcd"
	etcdConfig.NameFromHostname()
	etcdServer := etcd.New(etcdConfig)
	go util.Forever(func() {
		glog.Infof("Started etcd at http://%s", etcdAddr)
		etcdServer.Run()
	}, 0)

	etcdClient := etcdclient.NewClient(etcdServers)
	for i := 0; ; i += 1 {
		_, err := etcdClient.Get("/", false, false)
		if err == nil || tools.IsEtcdNotFound(err) {
			break
		}
		if i > 100 {
			glog.Fatal("Could not reach etcd: %v", err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	// initialize Kubelet
	minionHost := "127.0.0.1"
	minionPort := 10250
	rootDirectory := path.Clean("/var/lib/openshift")
	os.MkdirAll(rootDirectory, 0750)
	cfg := kconfig.NewPodConfig(kconfig.PodConfigNotificationSnapshotAndUpdates)
	kconfig.NewSourceEtcd(kconfig.EtcdKeyForHost(minionHost), etcdClient, cfg.Channel("etcd"))
	k := kubelet.NewMainKubelet(
		minionHost,
		dockerClient,
		cadvisorClient,
		etcdClient,
		rootDirectory,
		30*time.Second)
	go util.Forever(func() { k.Run(cfg.Updates()) }, 0)
	go util.Forever(func() {
		kubelet.ListenAndServeKubeletServer(k, cfg.Channel("http"), minionHost, uint(minionPort))
	}, 0)

	// initialize OpenShift API
	storage := map[string]apiserver.RESTStorage{
		"services":        service.NewRESTStorage(service.MakeMemoryRegistry()),
		"templateConfigs": template.NewStorage(),
	}
	osAddr := "127.0.0.1:8080"
	osPrefix := "/osapi/v1beta1"

	osMux := http.NewServeMux()

	// initialize Kubernetes API
	kubePrefix := "/api/v1beta1"
	kubeClient := client.New("http://"+osAddr, nil)
	podInfoGetter := &client.HTTPPodInfoGetter{
		Client: http.DefaultClient,
		Port:   uint(minionPort),
	}
	masterConfig := &master.Config{
		Client:             kubeClient,
		EtcdServers:        etcdServers,
		HealthCheckMinions: true,
		Minions:            []string{minionHost},
		PodInfoGetter:      podInfoGetter,
	}
	m := master.New(masterConfig)

	apiserver.NewAPIGroup(m.API_v1beta1()).InstallREST(osMux, kubePrefix)
	apiserver.NewAPIGroup(storage, api.Codec).InstallREST(osMux, osPrefix)
	apiserver.InstallSupport(osMux)

	osApi := &http.Server{
		Addr:           osAddr,
		Handler:        apiserver.RecoverPanics(osMux),
		ReadTimeout:    5 * time.Minute,
		WriteTimeout:   5 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	go util.Forever(func() {
		glog.Infof("Started Kubernetes API at http://%s%s", osAddr, kubePrefix)
		glog.Infof("Started OpenShift API at http://%s%s", osAddr, osPrefix)
		glog.Fatal(osApi.ListenAndServe())
	}, 0)

	// initialize kube proxy
	serviceConfig := pconfig.NewServiceConfig()
	endpointsConfig := pconfig.NewEndpointsConfig()
	pconfig.NewConfigSourceEtcd(etcdClient,
		serviceConfig.Channel("etcd"),
		endpointsConfig.Channel("etcd"))
	loadBalancer := proxy.NewLoadBalancerRR()
	proxier := proxy.NewProxier(loadBalancer)
	serviceConfig.RegisterHandler(proxier)
	endpointsConfig.RegisterHandler(loadBalancer)
	glog.Infof("Started Kubernetes Proxy")

	// initialize replication manager
	controllerManager := controller.NewReplicationManager(kubeClient)
	controllerManager.Run(10 * time.Second)
	glog.Infof("Started Kubernetes Replication Manager")

	// initialize scheduler
	configFactory := &factory.ConfigFactory{Client: kubeClient}
	config := configFactory.Create()
	s := scheduler.New(config)
	s.Run()
	glog.Infof("Started Kubernetes Scheduler")

	select {}
}

func env(key string, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	} else {
		return val
	}
}

func getDockerEndpoint(dockerEndpoint string) string {
	var endpoint string
	if len(dockerEndpoint) > 0 {
		endpoint = dockerEndpoint
	} else if len(os.Getenv("DOCKER_HOST")) > 0 {
		endpoint = os.Getenv("DOCKER_HOST")
	} else {
		endpoint = "unix:///var/run/docker.sock"
	}
	glog.Infof("Connecting to docker on %s", endpoint)

	return endpoint
}
