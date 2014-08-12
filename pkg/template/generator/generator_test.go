package generator

import (
	"math/rand"
	"testing"
)

func TestCreateGenerator(t *testing.T) {
	_, err := New(rand.New(rand.NewSource(1337)))
	if err != nil {
		t.Errorf("Failed to create generator")
	}
}

func TestExpressionValueGenerator(t *testing.T) {
	sampleGenerator, _ := New(rand.New(rand.NewSource(1337)))

	var tests = []struct {
		Expression    string
		ExpectedValue string
	}{
		{"test[A-Z0-9]{4}template", "testQ3HVtemplate"},
		{"[\\d]{4}", "6841"},
		{"[\\w]{4}", "DVgK"},
		{"[\\a]{10}", "nFWmvmjuaZ"},
		{"admin[0-9]{2}[A-Z]{2}", "admin32VU"},
		{"admin[0-9]{2}test[A-Z]{2}", "admin56testGS"},
		{"[password]", "x4sLZWES"},
	}

	for _, test := range tests {
		value, err := sampleGenerator.GenerateValue(test.Expression)
		if err != nil {
			t.Errorf("Failed to generate value from %s due to error: %v", test.Expression, err)
		}
		if value != test.ExpectedValue {
			t.Errorf("Failed to generate expected value from %s\n. Generated: %s\n. Expected: %s\n", test.Expression, value, test.ExpectedValue)
		}
	}

}

func TestExpressionValueGeneratorErrors(t *testing.T) {
	sampleGenerator, _ := New(rand.New(rand.NewSource(1337)))

	if v, err := sampleGenerator.GenerateValue("[ABC]{3}"); err == nil {
		t.Errorf("Expected [ABC]{3} to produce malformed syntax error (returned: %s)", v)
	}

	if v, err := sampleGenerator.GenerateValue("[GET:http://custom.url.int/new]"); err == nil {
		t.Errorf("Expected dial tcp error, got", v)
	}

	if v, err := sampleGenerator.GenerateValue("test"); err == nil {
		t.Errorf("Expected No matching generators found, got %s", v)
	}

	if v, err := sampleGenerator.GenerateValue("[Z-A]{3}"); err == nil {
		t.Errorf("Expected Invalid range specified error, got %s", v)
	}

	if v, err := sampleGenerator.GenerateValue("[A-Z]{300}"); err == nil {
		t.Errorf("Expected Invalid range specified error, got %s", v)
	}

	if v, err := sampleGenerator.GenerateValue("[A-Z]{0}"); err == nil {
		t.Errorf("Expected Invalid range specified error, got %s", v)
	}
}

func TestPasswordGenerator(t *testing.T) {
	sampleGenerator, _ := New(rand.New(rand.NewSource(1337)))

	value, _ := sampleGenerator.GenerateValue("[password]")
	if value != "4U390O49" {
		t.Errorf("Failed to generate expected password. Generated: %s\n. Expected: %s\n", value, "4U390O49")
	}
}

func TestRemoteValueGenerator(t *testing.T) {
	sampleGenerator, _ := New(rand.New(rand.NewSource(1337)))

	_, err := sampleGenerator.GenerateValue("[GET:http://api.example.com/new]")
	if err == nil {
		t.Errorf("Expected error while fetching non-existent remote.")
	}
}
