package main

import (
	"flag"
	"fmt"
	"hps"
	"hps/engine"
	jf "hps/jsonfactory"
	"hps/loggers"
	"hps/tools"
	"hps/tools/rnd"
	"nordic32/model"
	"nordic32/plugins"
	"nordic32/plugins/attacker"
	"nordic32/plugins/control"
	"nordic32/plugins/counters"
	"nordic32/plugins/diversitystudy"
	"nordic32/plugins/inspector"
	"nordic32/plugins/loadflow"
	"nordic32/plugins/modificator"
	"nordic32/plugins/protection"
	"strings"
	"time"
)

// Flags is the command line configuration options.
type Flags struct {
	Duration int
	File     string
	Network  string
	Seed     uint64
	Time     float64
	Events   int
	Version  bool
}

func getFlags() Flags {
	duration := flag.Int("duration", 0, "simulate model for the specified number of seconds")
	file := flag.String("file", "", "path to the model file")
	network := flag.String("network", "", "model file network to run")
	seed := flag.Uint64("seed", 0, "random number generator seed")
	time := flag.Float64("time", 0, "simulate model to the specified time")
	events := flag.Int("events", 0, "simulate specific number of events")
	version := flag.Bool("version", false, "print version info and exit")

	flag.Parse()

	return Flags{*duration, *file, *network, *seed, *time, *events, *version}
}

func getPlugins(logger hps.ILogger, rnd hps.IRnd) *plugins.PluginsContainer {
	// Plugins order is important
	loadflowPlugin := loadflow.NewPlugin(logger).(*loadflow.Plugin)
	return plugins.CreatePluginsContainer(
		loadflowPlugin,
		diversitystudy.NewPlugin(logger, rnd),
		protection.NewPlugin(logger),
		attacker.NewPlugin(logger, rnd),
		modificator.NewPlugin(logger, rnd),
		inspector.NewPlugin(logger, rnd),
		control.NewPlugin(logger),
		counters.NewPlugin(logger, loadflowPlugin),
		plugins.NewCustomBreakerPlugin(logger),
	)
}

// Rnd is the random number generator used for simulations.
type Rnd struct {
	mt *rnd.MT19937
}

// Next returns the next random floating point number from
// the MT19937 random number generator.
func (r *Rnd) Next() float64 {
	return r.mt.Float3()
}

func main() {
	r := &Rnd{}
	hps.SetRnd(hps.IRnd(r))

	flags := getFlags()

	if flags.Version {
		fmt.Printf("build %s created on %s", hps.BUILD_ID, hps.BUILD_DATE)
		return
	}

	loggerOutput := loggers.NewConsoleLogger()

	logger := tools.NewEventLog(loggerOutput)

	if flags.File == "" {
		logger.Print("model file is not specified")
		return
	}

	f := jf.NewFactory(model.TypeFactory(), hps.IRnd(r))
	err := f.Load(flags.File)
	if err != nil {
		logger.Print(err)
		return
	}

	networkNames := f.GetNetworkNames()

	if len(networkNames) == 0 {
		logger.Print("no networks defined in the model file")
		return
	}

	if len(networkNames) > 1 && flags.Network == "" {
		logger.Printf(
			"network should be one of %s",
			strings.Join(networkNames, ", "))
		return
	}

	networkName := ""
	if len(networkNames) == 1 && flags.Network == "" {
		networkName = networkNames[0]
	} else {
		networkName = flags.Network
	}

	machines, err := f.CreateNetwork(networkName)
	if err != nil {
		logger.Print(err)
		return
	}

	env := engine.NewEnvironment(logger)

	for _, m := range machines {
		env.AddMachine(m)
	}

	logger.AttachTraces(env)

	seed := uint64(0)
	if flags.Seed == 0 {
		seed = uint64(time.Now().UnixNano())
	} else {
		seed = flags.Seed
	}

	r.mt = rnd.MT19937New(uint32(seed))

	limits := engine.NewLimits()

	if flags.Duration == 0 && flags.Time == 0 && flags.Events == 0 {
		limits.Duration = 3
		limits.Events = 10
	} else {
		limits.Duration = flags.Duration
		limits.Time = flags.Time
		limits.Events = flags.Events
	}

	logger.Printf("file = %v", flags.File)
	logger.Printf("network = %v", networkName)
	logger.Printf("seed = %v", seed)

	plugins := getPlugins(logger, hps.IRnd(r))

	if err := plugins.Init(env); err != nil {
		logger.Print(err)
		return
	}

	result := env.Run(limits)

	plugins.Done(result)
}
