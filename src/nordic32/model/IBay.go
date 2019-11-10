package model

import "hps"

type IBay interface {
	Breaker() (hps.IMachine, bool)
}
