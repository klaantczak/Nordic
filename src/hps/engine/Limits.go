package engine

// Limits sets the simulation limits, such as number of iterations, total
// simulation time, duration of the simulation.
type Limits struct {
	// Limit simulation by the number of events
	Events int

	// Limit simulation by the environment time
	Time float64

	// Limit simlation by the running time
	Duration int

	// Limit by custom predicate
	Predicate func() bool
}

// NewLimits creates simulation limits with the default values (no limits).
func NewLimits() *Limits {
	return &Limits{}
}
