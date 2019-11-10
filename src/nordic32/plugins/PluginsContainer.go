package plugins

import (
	"hps"
	"hps/engine"
)

type PluginsContainer struct {
	plugins []IPlugin
}

func CreatePluginsContainer(p ...IPlugin) *PluginsContainer {
	ps := &PluginsContainer{[]IPlugin{}}
	for _, i := range p {
		ps.plugins = append(ps.plugins, i)
	}
	_ = IPlugin(ps)
	return ps
}

func (ps *PluginsContainer) Names() string {
	names := ""
	switch len(ps.plugins) {
	case 0:
		names = "no plugins"
	case 1:
		names = "plugin " + ps.plugins[0].Name()
	default:
		for _, p := range ps.plugins {
			names += " " + p.Name()
		}
		names = "plugins" + names
	}
	return names
}

func (ps *PluginsContainer) Name() string {
	return "plugins"
}

func (ps *PluginsContainer) Init(e hps.IEnvironment) error {
	for _, p := range ps.plugins {
		if err := p.Init(e); err != nil {
			return err
		}
	}
	return nil
}

func (ps *PluginsContainer) Done(r *engine.SimulationResult) {
	for _, p := range ps.plugins {
		p.Done(r)
	}
}
