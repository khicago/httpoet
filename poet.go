package httpoet

type Poet struct {
	host  string
	baseH IHeader
}

func New(host string) *Poet {
	poet := &Poet{
		host: host,
		baseH: make(Hs),
	}
	return poet
}

func (hp *Poet) SetBaseH(header IHeader) *Poet {
	if header == nil {
		return hp
	}
	if hp.baseH == nil {
		hp.baseH = make(Hs) // cow & immutable
	}
	hp.baseH = hp.baseH.WithH(header) // will override existed keys
	return hp
}
