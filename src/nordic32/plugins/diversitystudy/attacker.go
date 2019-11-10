package diversitystudy

import (
	"hps"
	nm "hps/networkmachine"
	sm "hps/statemachine"
	n32model "nordic32/model"
	"nordic32/plugins/diversitystudy/model"
)

func rndOne(rnd hps.IRnd, n int) int {
	return int(rnd.Next() * float64(n))
}

func attachAttackerHandlers(attacker *model.Attacker, network *n32model.Network, rnd hps.IRnd) {
	act := attacker.StateAttack()
	act.Entering(func(p *sm.State, time float64, duration float64) {
		switch attacker.Target() {
		case "one":
			attack(time, rnd, network.Substations(), selectBays, attackOneBreaker)
		case "all":
			attack(time, rnd, network.Substations(), selectBays, attackAllBreakers)
		case "lines,one":
			attack(time, rnd, network.Substations(), selectLineBays, attackOneBreaker)
		case "lines,all":
			attack(time, rnd, network.Substations(), selectLineBays, attackAllBreakers)
		case "loads,one":
			attack(time, rnd, network.SubstationsWithLoads(), selectLoadBays, attackOneBreaker)
		case "loads,all":
			attack(time, rnd, network.SubstationsWithLoads(), selectLoadBays, attackAllBreakers)
		case "generators,one":
			attack(time, rnd, network.SubstationsWithGenerators(), selectGeneratorBays, attackOneBreaker)
		case "generators,all":
			attack(time, rnd, network.SubstationsWithGenerators(), selectGeneratorBays, attackAllBreakers)
		default:
			panic("attack target is unknown")
		}
	})
}

func attack(time float64, rnd hps.IRnd, substations []*n32model.Substation,
	selectBays func(*n32model.Substation) []n32model.IBay,
	attackComponents func(float64, hps.IRnd, []*model.BreakerComponent)) {

	substation := substations[rndOne(rnd, len(substations))]
	if bays := selectBays(substation); len(bays) > 0 {
		bay := bays[rndOne(rnd, len(bays))]
		machine, _ := bay.Breaker()
		breaker := model.NewBreaker(machine.(*nm.Machine))
		components := breaker.Components()
		attackComponents(time, rnd, components)
	}
}

func selectBays(substation *n32model.Substation) []n32model.IBay {
	result := []n32model.IBay{}
	result = append(result, selectLineBays(substation)...)
	result = append(result, selectLoadBays(substation)...)
	result = append(result, selectGeneratorBays(substation)...)
	return result
}

func selectLoadBays(substation *n32model.Substation) []n32model.IBay {
	result := []n32model.IBay{}
	for _, v := range substation.LoadBays() {
		result = append(result, n32model.IBay(v))
	}
	return result
}

func selectLineBays(substation *n32model.Substation) []n32model.IBay {
	result := []n32model.IBay{}
	for _, v := range substation.LineBays() {
		result = append(result, n32model.IBay(v))
	}
	return result
}

func selectGeneratorBays(substation *n32model.Substation) []n32model.IBay {
	result := []n32model.IBay{}
	for _, v := range substation.GeneratorBays() {
		result = append(result, n32model.IBay(v))
	}
	return result
}

func attackOneBreaker(time float64, rnd hps.IRnd, components []*model.BreakerComponent) {
	component := components[rndOne(rnd, len(components))]
	component.Compromise(time)
}

func attackAllBreakers(time float64, rnd hps.IRnd, components []*model.BreakerComponent) {
	for _, component := range components {
		component.Compromise(time)
	}
}
