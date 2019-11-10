package hpstest

import (
	sm "hps/statemachine"
	"hps/statemachine/triggers"
	"math"
	"testing"
)

func Test_CreateTransition(t *testing.T) {
	tr := sm.CreateTransition("test", triggers.NewDeterministicTrigger(0))
	if tr.To() != "test" {
		t.Errorf("The event target should be \"test\".")
	}
}

func Test_CheckTransitionFrequency(t *testing.T) {
	experiment := func(parameter float64, frequency float64) {
		tr := sm.CreateTransition("", triggers.NewProbabilisticTrigger(parameter))
		sum := 0.0
		for i := 0; i < 10000; i++ {
			sum += tr.Trigger().Time(0.0)
		}
		if math.Abs(sum/10000.0-frequency) > frequency*0.1 {
			t.Errorf("Expected frequency %v but have got %v", frequency, sum/10000.0)
		}
	}

	// 100 times in a time slot
	experiment(100.0, 0.01)

	// 10 times in a time slot
	experiment(10, 0.1)

	// once in a time slot
	experiment(1, 1)

	// once in 10 time slots
	experiment(0.1, 10)

	// once in 100 time slots
	experiment(0.01, 100)
}
