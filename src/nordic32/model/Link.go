package model

import (
	"hps"
	sm "hps/statemachine"
)

type Link struct {
	*sm.Machine
	connected  *hps.Property
	overloaded *hps.Property
	max        *hps.Property
}

func NewLink(machine *sm.Machine) *Link {
	connected, _ := machine.Property("connected")
	overloaded, _ := machine.Property("overloaded")
	max, _ := machine.Property("max")
	return &Link{machine, connected, overloaded, max}
}

func (l *Link) Ok() bool {
	return l.State().Name() == "ok"
}

func (l *Link) StateOk() *sm.State {
	ok, _ := l.GetState("ok")
	return ok
}

func (l *Link) StateFail() *sm.State {
	ok, _ := l.GetState("fail")
	return ok
}

func (l *Link) Connected() bool {
	value, _ := l.connected.GetBool()
	return value
}

func (l *Link) Disconnected() bool {
	value, _ := l.connected.GetBool()
	return !value
}

func (l *Link) PropertyConnected() *hps.Property {
	return l.connected
}

func (l *Link) Connect() {
	l.connected.SetValue(true)
}

func (l *Link) Disconnect() {
	l.connected.SetValue(false)
}

func (l *Link) Max() float64 {
	value, _ := l.max.GetFloat()
	return value
}

func (l *Link) SetMax(value float64) {
	l.max.SetValue(value)
}

func (l *Link) PropertyMax() *hps.Property {
	return l.max
}

func (l *Link) From() string {
	property, _ := l.Property("from")
	value, _ := property.GetString()
	return value
}

func (l *Link) To() string {
	property, _ := l.Property("to")
	value, _ := property.GetString()
	return value
}

func (l *Link) X() float64 {
	property, _ := l.Property("x")
	value, _ := property.GetFloat()
	return value
}

func (l *Link) Overloaded() bool {
	property, _ := l.Property("overloaded")
	value, _ := property.GetBool()
	return value
}

func (l *Link) Overload() {
	property, _ := l.Property("overloaded")
	property.SetValue(true)
}

func (l *Link) NormalLoad() {
	property, _ := l.Property("overloaded")
	property.SetValue(false)
}

func (l *Link) PropertyOverloaded() *hps.Property {
	return l.overloaded
}
