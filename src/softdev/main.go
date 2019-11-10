package main

import (
	"fmt"
	"hps"
	"hps/engine"
	"hps/loggers"
	sm "hps/statemachine"
	"hps/statemachine/triggers"
	"hps/tools"
	"hps/tools/rnd"
	"time"
)

type Fragment struct {
}

const IssueStatusOnHold = 0
const IssueStatusInProcess = 1
const IssueStatusDone = 2

type Issue struct {
	Status int
}

type Artefacts struct {
	Code   []*Fragment
	Issues []*Issue
}

func NewArtefacts() *Artefacts {
	return &Artefacts{
		[]*Fragment{},
		[]*Issue{},
	}
}

func (a *Artefacts) AddIssue() {
	a.Issues = append(a.Issues, &Issue{IssueStatusOnHold})
}

func (a *Artefacts) FetchIssue() *Issue {
	for _, i := range a.Issues {
		if i.Status == IssueStatusOnHold {
			i.Status = IssueStatusInProcess
			return i
		}
	}
	return nil
}

func (a *Artefacts) AddCodeFragment() {
	a.Code = append(a.Code, &Fragment{})
}

func (a *Artefacts) ModifyCodeFragment() {
}

func makeDeveloper(artefacts *Artefacts, r hps.IRnd) hps.IMachine {
	p := func(d float64) sm.ITrigger {
		return sm.ITrigger(triggers.NewProbabilisticTrigger(1 / d))
	}

	m, _ := sm.NewMachine(
		"Developer",
		"Developer",
		[]string{"select", "analyse", "add-code", "modify-code", "commit-change"},
		"select")
	// Developer selects the issue to work on,
	selectSt, _ := m.GetState("select")
	// then analyses the issue,
	analyse, _ := m.GetState("analyse")
	// then either adds or modifies the code,
	addCode, _ := m.GetState("add-code")
	modifyCode, _ := m.GetState("modify-code")
	// then commits change and move back to choosing the next issue to work on.
	commitChange, _ := m.GetState("commit-change")

	tenMinutes, _ := tools.ParseDuration("10 minutes")
	fourHours, _ := tools.ParseDuration("4 hours")

	selectSt.AddTransition("analyse", p(tenMinutes))
	analyse.AddTransition("add-code", p(fourHours))
	analyse.AddTransition("modify-code", p(fourHours))
	addCode.AddTransition("commit-change", p(tenMinutes))
	modifyCode.AddTransition("commit-change", p(tenMinutes))
	commitChange.AddTransition("select", triggers.NewDeterministicTrigger(0))

	var issue *Issue

	selectSt.Leaving(func(p *sm.State, time float64, duration float64) {
		issue = artefacts.FetchIssue()
	})

	addCode.Leaving(func(p *sm.State, time float64, duration float64) {
		artefacts.AddCodeFragment()
	})

	modifyCode.Leaving(func(p *sm.State, time float64, duration float64) {
		artefacts.ModifyCodeFragment()
	})

	commitChange.Leaving(func(p *sm.State, time float64, duration float64) {
		if issue == nil {
			return
		}

		issue.Status = IssueStatusDone
	})

	return hps.IMachine(m)
}

// Random number generator wrapper
type Rnd struct {
	mt *rnd.MT19937
}

func (r *Rnd) Next() float64 {
	return r.mt.Float3()
}

func main() {
	r := &Rnd{}
	r.mt = rnd.MT19937New(uint32(time.Now().UnixNano()))

	artefacts := NewArtefacts()
	for i := 0; i < 100; i++ {
		artefacts.AddIssue()
	}

	ir := hps.IRnd(r)

	e := engine.NewEnvironment(loggers.NewConsoleLogger())
	e.AddMachine(makeDeveloper(artefacts, ir))
	e.AddMachine(makeDeveloper(artefacts, ir))
	e.AddMachine(makeDeveloper(artefacts, ir))

	limits := engine.NewLimits()
	limits.Predicate = func() bool {
		for _, i := range artefacts.Issues {
			if i.Status != IssueStatusDone {
				return true
			}
		}
		return false
	}

	result := e.Run(limits)

	second, _ := tools.ParseDuration("second")
	fmt.Println("100 issues, 3 developers")
	fmt.Println(time.Duration(result.Time/second) * time.Second)
	fmt.Println(result.Time/second/(60*60*24), "working days")
	fmt.Printf("code: %d\n", len(artefacts.Code))
}
