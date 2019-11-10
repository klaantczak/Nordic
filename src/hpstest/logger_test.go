package hpstest

import (
	"hps/engine"
	"hps/loggers"
	sm "hps/statemachine"
	"hps/statemachine/triggers"
	"hps/tools"
	"testing"
)

// TODO this test needs refactoring
func Test_LogPropertyValue1(t *testing.T) {
	e := engine.NewEnvironment(loggers.NewConsoleLogger())

	m, _ := sm.NewMachine("Simple", "Simple", []string{"ok"}, "ok")
	m.AddProperty("trigger", "Trigger")
	trigger, _ := m.Property("trigger")
	trigger.SetValue(triggers.NewProbabilisticTrigger(1.0))

	e.AddMachine(m)

	loggerOutput := loggers.NewMemoryLogger()
	logger := tools.NewEventLog(loggerOutput)
	logger.AttachTraces(e)

	actual := loggerOutput.ToString()

	expected := ``

	if actual != expected {
		t.Errorf("Expect %s, got %s", expected, actual)
	}
}

// TODO this test needs refactoring
func Test_LogPropertyCustomValue(t *testing.T) {
	e := engine.NewEnvironment(loggers.NewConsoleLogger())

	m, _ := sm.NewMachine("Simple", "Simple", []string{"ok"}, "ok")
	m.AddProperty("custom", "")
	custom, _ := m.Property("custom")
	custom.SetValue(struct {
		Value string `json:"value"`
	}{"ok"})

	e.AddMachine(m)

	loggerOutput := loggers.NewMemoryLogger()
	logger := tools.NewEventLog(loggerOutput)
	logger.AttachTraces(e)

	actual := loggerOutput.ToString()

	expected := ``

	if actual != expected {
		t.Errorf("Expect %s, got %s", expected, actual)
	}
}
