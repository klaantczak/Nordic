package networkmachine

import "hps"

type INetworkMachine interface {
	Machines() []hps.IMachine
}
