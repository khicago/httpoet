package httpoet

type Poet struct {
	hosts []string
	baseH IHeader
}

func New(host string, hosts ...string) *Poet {
	if len(hosts) == 0 {
		hosts = []string{host}
	} else {
		hosts = append([]string{host}, hosts...)
	}

	poet := &Poet{
		hosts: hosts,
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
