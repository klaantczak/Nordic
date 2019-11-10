package hpstest

import (
	"hps/tools"
	"testing"
)

func Test_ParseDuration(t *testing.T) {
	if v, err := tools.ParseDuration("1"); err != nil || v != 1.0 {
		t.Errorf("Parsed '1' to %v, expect 1.0", v)
	}

	if v, err := tools.ParseDuration("1.5"); err != nil || v != 1.5 {
		t.Errorf("Parsed '1.5' to %v, expect 1.5", v)
	}

	if v, err := tools.ParseDuration("1year"); err != nil || v != 1 {
		t.Errorf("Parsed '1 year' to %v, expect 1.0", v)
	}

	if v, err := tools.ParseDuration("52 weeks"); err != nil || v != 1 {
		t.Errorf("Parsed '52 weeks' to %v, expect 1.0", v)
	}
	if v, err := tools.ParseDuration("10 min"); err != nil || v != 1.9025875190258754e-05 {
		t.Errorf("Parsed '10 min' to %v, expect 1.9...e-5", v)
	}
}
