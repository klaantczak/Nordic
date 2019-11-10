package statemachine

type IStateMachine interface {
	State() *State
	Changed(handler MachineEventHandler)
}
