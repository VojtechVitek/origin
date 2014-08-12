package generator

import (
	"fmt"
	"math/rand"
	"regexp"
	"sync"
)

// GeneratorInterface is an abstract interface for generating
// random values from an input expression
type GeneratorInterface interface {
	GenerateValue(expression string) (interface{}, error)
}

// Generator implements GeneratorInterface
type Generator struct {
	seed *rand.Rand
}

func New(seed *rand.Rand) (Generator, error) {
	return Generator{seed: seed}, nil
}

// MatchesGeneratorExpression loops over registered generators
// and searches for an expression matching the given string.
// Returns true if match is found, false otherwise.
func MatchesGeneratorExpression(value string) bool {
	for regexp, _ := range generators {
		if regexp.MatchString(value) {
			return true
		}
	}
	return false
}

// GenerateValue loops over registered generators and tries to find the
// one matching the input expression. If match is found, it then generates
// value using that matching generator
func (g *Generator) GenerateValue(expression string) (interface{}, error) {
	for regexp, generatorFactory := range generators {
		if regexp.FindStringIndex(expression) != nil {
			generator, _ := generatorFactory(g.seed)
			return generator.GenerateValue(expression)
		}
	}

	return expression, fmt.Errorf("No matching generators found.")
}

// GeneratorFactory is an abstract factory for creating generators
// (objects that implement GeneratorInterface interface)
type GeneratorFactory func(*rand.Rand) (GeneratorInterface, error)

// generators stores registered generators
var generators = make(map[*regexp.Regexp]GeneratorFactory)
var generatorsMutex sync.Mutex

// RegisterGenerator registers new generator for a certain input expression
func RegisterGenerator(r *regexp.Regexp, f GeneratorFactory) {
	generatorsMutex.Lock()
	defer generatorsMutex.Unlock()
	generators[r] = f
}
