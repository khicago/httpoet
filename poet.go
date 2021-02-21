package httpoet

import (
	"github.com/khicago/irr"
)

type Poet struct {
	host  string
	baseH Hs
}

func New(host string, header Hs) (*Poet, irr.IRR) {
	poet := &Poet{
		host: host,
	}
	if err := poet.AddBaseH(header); err != nil {
		return nil, err
	}
	return poet, nil
}

func (hp *Poet) AddBaseH(header Hs) irr.IRR {
	if header == nil || len(header) == 0 {
		return nil
	}
	newH := make(Hs) // cow & immutable
	for k, v := range header {
		newH[k] = v
	}

	if hp.baseH != nil {
		for k, v := range hp.baseH {
			if _, ok := newH[k]; ok {
				return irr.TraceSkip(1, "cannot override exist key= %s", k)
			}
			newH[k] = v
		}
	}

	hp.baseH = newH
	return nil
}

func (hp *Poet) OverrideBaseH(header Hs) irr.IRR {
	if header == nil || len(header) == 0 {
		return nil
	}

	newH := make(Hs) // cow & immutable
	if hp.baseH != nil {
		for k, v := range hp.baseH {
			newH[k] = v
		}
	}
	for k, v := range header {
		newH[k] = v
	}

	hp.baseH = newH
	return nil
}
