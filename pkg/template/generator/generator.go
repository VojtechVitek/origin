package generator

// GeneratorInterface is an interface for generating
// random values from an input expression
type GeneratorInterface interface {
	GenerateValue(expression string) (interface{}, error)
}
