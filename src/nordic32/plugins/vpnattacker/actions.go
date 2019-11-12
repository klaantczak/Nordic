package vpnattacker

import (
	"hps"
	sm "hps/statemachine"
	"nordic32/model"
)

func rndOne(rnd hps.IRnd, n int) int {
	return int(rnd.Next() * float64(n))
}

func actions(m *VpnAttacker, network *model.Network, rnd hps.IRnd) {

}
