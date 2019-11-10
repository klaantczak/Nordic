package hps

type ITypeFactory func(IMachine) IMachine

// Represents machine factory, object that creates different machine kinds.
type IFactory interface {
	Create(string, string) (IMachine, error)
}

type IMachine interface {
	Name() string
	Kind() string
	Event(float64) (IEvent, bool)
	Property(string) (*Property, bool)
	Properties() []*Property
}

type ILogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

type IEvent interface {
	Time() float64
	Invoke() float64
	Machine() IMachine
	String() string
}

// Event handler type for the simlation events.
type EnvironmentEventHandler func(IEnvironment, IEvent)

type EnvironmentLogHandler func(IEnvironment, float64, string)

type IEnvironment interface {
	StartingIteration(handler EnvironmentEventHandler)
	EndingIteration(handler EnvironmentEventHandler)
	Machines() []IMachine
	Machine(name string) (IMachine, bool)
	AddMachine(machine IMachine)
	RemoveMachine(machine IMachine)
}

type IRnd interface {
	// Returns integer from [0;1).
	Next() float64
}

type IJSON interface {
	ToJSON() map[string]interface{}
}
