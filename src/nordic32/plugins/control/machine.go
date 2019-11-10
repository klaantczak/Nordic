package control

import (
	"hps"
	sm "hps/statemachine"
	"hps/statemachine/triggers"
	"hps/tools"
)

func addControlMachine(e hps.IEnvironment) *sm.Machine {
	m, _ := sm.NewMachine("Control", "Control", []string{"idle", "working"}, "idle")

	idle, _ := m.GetState("idle")
	working, _ := m.GetState("working")

	m.AddProperty("start", "Activation")
	start, _ := m.Property("start")
	t := triggers.NewPropertyTrigger("start")

	idle.AddTransition("working", t)

	d, _ := tools.ParseDuration("1s")
	working.AddTransition("idle", triggers.NewDeterministicTrigger(d))

	start.SetValue(triggers.NewIdleTrigger())

	e.AddMachine(m)

	return m
}
