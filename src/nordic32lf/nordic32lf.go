// package name: nordic32lf
package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"hps/engine"
	jf "hps/jsonfactory"
	nm "hps/networkmachine"
	"hps/statemachine/triggers"
	"io/ioutil"
	ctrl "nordic32/plugins/control"
	lf "nordic32/plugins/loadflow"
	q "nordic32/query"
	"unsafe"
)

var modelFilePath string
var modelLength int

func read(filename string) *nm.Machine {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}

	f := jf.NewFactory()
	err = f.LoadJSON(fileContent)
	if err != nil {
		panic(err.Error())
	}

	machines, err := f.CreateNetwork("baseline")
	if err != nil {
		panic(err.Error())
	}

	env := engine.NewEnvironment()
	for _, m := range machines {
		env.AddMachine(m)
	}

	substations, _ := q.FindSubstationsNetwork(env)

	return substations
}

//export Init
func Init(filename string) (C.size_t, *C.double) {
	modelFilePath = filename
	substations := read(modelFilePath)
	result := []float64{}
	for _, s := range q.FindSubstations(substations) {
		for _, gb := range q.FindGeneratorBays(s) {
			g, _ := q.FindGenerator(gb)
			ok, _ := g.GetState("ok")
			toFail, _ := q.GetTransitionTo(ok, "fail")
			result = append(result, toFail.Trigger().(*triggers.ProbabilisticTrigger).GetParameter())
		}
		for _, ldb := range q.FindLoadBays(s) {
			ld, _ := q.FindLoad(ldb)
			ok, _ := ld.GetState("ok")
			toFail, _ := q.GetTransitionTo(ok, "fail")
			result = append(result, toFail.Trigger().(*triggers.ProbabilisticTrigger).GetParameter())
		}
	}
	for _, lnk := range q.FindLinks(substations) {
		ok, _ := lnk.GetState("ok")
		toFail, _ := q.GetTransitionTo(ok, "fail")
		propertyName := toFail.Trigger().(*triggers.PropertyTrigger).GetProperty()
		property, _ := lnk.Property(propertyName)
		parameter := property.Value().(*triggers.ProbabilisticTrigger).GetParameter()
		result = append(result, parameter)
	}

	rlen := len(result)

	modelLength = rlen

	// allocate the *C.double array
	p := C.malloc(C.size_t(rlen) * C.size_t(unsafe.Sizeof(C.double(0))))

	// convert the pointer to a go slice so we can index it
	r := (*[1<<30 - 1]C.double)(p)[:rlen:rlen]
	for i, v := range result {
		r[i] = C.double(v)
	}

	return C.size_t(rlen), (*C.double)(p)
}

//export Run
func Run(input *C.int) (C.size_t, *C.int) {
	values := (*[1<<30 - 1]C.int)(unsafe.Pointer(input))[:modelLength:modelLength]
	states := make([]int, modelLength)
	for i := 0; i < modelLength; i++ {
		states[i] = int(values[i])
	}

	substations := read(modelFilePath)

	i := 0
	for _, s := range q.FindSubstations(substations) {
		for _, gb := range q.FindGeneratorBays(s) {
			g, _ := q.FindGenerator(gb)
			if states[i] == 1 {
				g.SetState("fail")
			}
			i++
		}
		for _, ldb := range q.FindLoadBays(s) {
			ld, _ := q.FindLoad(ldb)
			if states[i] == 1 {
				ld.SetState("fail")
			}
			i++
		}
	}
	for _, lnk := range q.FindLinks(substations) {
		if states[i] == 1 {
			lnk.SetState("fail")
		}
		i++
	}

	model := lf.BuildModel(substations)

	for i := 0; i < len(model.Nodes); i++ {
		status := lf.StatusEnabled
		if states[i] == 1 {
			status = lf.StatusDisabled
		}
		model.Nodes[i].Status = status
	}
	for i, j := len(model.Nodes), 0; j < len(model.Links); i, j = i+1, j+1 {
		model.Links[j].Enabled = states[i] != 1
	}

	for _, link := range lf.GetTestedListOfOverloadedLines(substations, model) {
		overloaded, _ := link.Machine.Property("overloaded")
		value, _ := overloaded.GetBool()
		if !value {
			overloaded.SetValue(true)
		}
	}

	actions := ctrl.BuildActions(substations)
	ctrl.ApplyActions(actions, ctrl.BuildGraphFromNetwork(substations))

	result := []int{}
	for _, s := range q.FindSubstations(substations) {
		for _, gb := range q.FindGeneratorBays(s) {
			g, _ := q.FindGenerator(gb)
			connected := q.GetBoolProp(gb, "connected")
			ok := q.IsOk(g)
			state := 0
			if !ok {
				state = 1
			}
			if !connected {
				state = 2
			}
			result = append(result, state)
		}
		for _, ldb := range q.FindLoadBays(s) {
			ld, _ := q.FindLoad(ldb)
			connected := q.GetBoolProp(ldb, "connected")
			ok := q.IsOk(ld)
			state := 0
			if !ok {
				state = 1
			}
			if !connected {
				state = 2
			}
			result = append(result, state)
		}
	}
	for _, lnk := range q.FindLinks(substations) {
		connected := q.GetBoolProp(lnk, "connected")
		ok := q.IsOk(lnk)
		state := 0
		if !ok {
			state = 1
		}
		if !connected {
			state = 2
		}
		result = append(result, state)
	}

	rlen := len(result)
	// allocate the *C.double array
	p := C.malloc(C.size_t(rlen) * C.size_t(unsafe.Sizeof(C.int(0))))

	// convert the pointer to a go slice so we can index it
	r := (*[1<<30 - 1]C.int)(p)[:rlen:rlen]
	for i, v := range result {
		r[i] = C.int(v)
	}

	return C.size_t(rlen), (*C.int)(p)
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C shared library
}
