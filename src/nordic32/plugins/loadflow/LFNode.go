package loadflow

type LFNode struct {
	ID         int
	Name       string
	Capacity   float64
	Type       string
	Generators []*LFGenerator
	Loads      []*LFLoad

	// Calculated field
	Power float64

	// Calculated field
	Status int
}
