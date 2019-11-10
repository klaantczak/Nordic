package experiments

import (
	"hps/engine"
	"hps/statemachine"
	"hps/statemachine/triggers"
	"hps/tools"
	"nordic32/query"
)

func init() {
	constructors["modificator.weekly-attacks.weekly-and-daily-inspections"] = func(environment *engine.Environment) *Experiment {
		//plg.Modificator.SetFrequency("weekly")
		//plg.Inspector.SetFrequency("probabilistic", "week")

		(func() {
			d25, _ := tools.ParseDuration("2.5 years")
			d50, _ := tools.ParseDuration("5 years")

			m := statemachine.NewMachine("special", "")
			m.AddState("weekly").Entering(func(p *statemachine.State, time float64, duration float64) {
				inspector, _ := query.FindMachineByName(environment.Machines(), "Inspector")
				idle, _ := inspector.(*statemachine.Machine).GetState("idle")
				frequency, _ := tools.ParseDuration("week")
				idle.Transitions()[0].SetTrigger(statemachine.ITrigger(triggers.NewProbabilisticTrigger(1 / frequency)))
			}).AddTransition("daily", triggers.NewDeterministicTrigger(d25))
			m.AddState("daily").Entering(func(p *statemachine.State, time float64, duration float64) {
				inspector, _ := query.FindMachineByName(environment.Machines(), "Inspector")
				idle, _ := inspector.(*statemachine.Machine).GetState("idle")
				frequency, _ := tools.ParseDuration("day")
				idle.Transitions()[0].SetTrigger(statemachine.ITrigger(triggers.NewProbabilisticTrigger(1 / frequency)))
			}).AddTransition("weekly", triggers.NewDeterministicTrigger(d50))
			m.SetState("weekly")
			environment.AddMachine(m)
		})()

		return &Experiment{environment}
	}
}
