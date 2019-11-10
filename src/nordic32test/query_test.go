package nordic32test

import (
	"hps/engine"
	jf "hps/jsonfactory"
	"hps/loggers"
	"nordic32/model"
	q "nordic32/query"
	"os"
	"testing"
)

func substations() *model.Network {
	e := engine.NewEnvironment(loggers.NewConsoleLogger())

	f := jf.NewFactory(model.TypeFactory(), nil)

	err := f.Load(os.Getenv("GOPATH") + "/models/nordic32.json")
	if err != nil {
		panic(err)
	}

	machines, err := f.CreateNetwork("baseline")
	if err != nil {
		panic(err)
	}

	for _, m := range machines {
		e.AddMachine(m)
	}

	substations, _ := q.FindSubstationsNetwork(e)

	return substations
}

func Test_FindSubstations(t *testing.T) {
	list := substations().Substations()
	if len(list) != 32 {
		t.Errorf("Expect 32 substations")
	}
}

func Test_FindLoadBays(t *testing.T) {
	loads := map[string]int{
		"1013": 1, "4043": 1, "1042": 1, "1045": 1, "4062": 1, "4047": 1, "4072": 1,
		"1022": 1, "2031": 1, "4041": 1, "4071": 1, "1012": 1, "4061": 1, "1043": 1,
		"1044": 1, "1041": 1, "1011": 1, "2032": 1, "4042": 1, "4046": 1, "4063": 1,
		"4051": 1}

	for _, s := range substations().Substations() {
		if len(s.LoadBays()) != loads[s.Name()] {
			t.Errorf("Substation %s, expect %d load bays", s.Name(), loads[s.Name()])
		}
	}
}

func Test_FindGeneratorBays(t *testing.T) {
	generators := map[string]int{
		"4041": 1, "1043": 1, "1042": 1, "4071": 1, "1013": 1, "4031": 1, "1012": 1,
		"4072": 1, "2032": 1, "4042": 1, "4062": 1, "4047": 1, "4051": 1, "4011": 1,
		"4012": 1, "1021": 1, "1022": 1, "4063": 1, "1014": 1, "4021": 1}

	for _, s := range substations().Substations() {
		if len(s.GeneratorBays()) != generators[s.Name()] {
			t.Errorf("Substation %s, expect %d generator bays", s.Name(), generators[s.Name()])
		}
	}
}

func Test_FindWorkingLoadBays(t *testing.T) {
	s := substations()
	m, _ := q.FindMachineByPath(s.Machines(), "1013")
	s1013 := m.(*model.Substation)

	bays1 := s1013.WorkingLoadBays()
	if len(bays1) != 1 {
		t.Errorf("Expect one working load bay")
	}

	bays1[0].Disconnect()

	bays2 := s1013.WorkingLoadBays()
	if len(bays2) != 0 {
		t.Errorf("Expect no working load bays")
	}
}
