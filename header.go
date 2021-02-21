package httpoet

type (
	headerIdentity struct{}

	IHeader interface {
		foreach(fn func(key string, value ...string)) *headerIdentity

		WithKV(key string, value ...string) IHeader
		WithKVAppend(key string, value ...string) IHeader

		WithH(header IHeader) IHeader
		WithHAppend(header IHeader) IHeader

		CountOf(key string) int
	}
)

type (
	H  map[string]string
	Hs map[string][]string
)

func (h H) foreach(fn func(key string, values ...string)) *headerIdentity {
	for k, v := range h {
		fn(k, v)
	}
	return &headerIdentity{}
}

func (hs Hs) foreach(fn func(key string, values ...string)) *headerIdentity {
	for k, vs := range hs {
		fn(k, vs...)
	}
	return &headerIdentity{}
}

/////////////// Region Basic

func (h H) CountOf(key string) int {
	if _, ok := h[key]; ok {
		return 1
	}
	return 0
}

func (hs Hs) CountOf(key string) int {
	if vs, ok := hs[key]; ok {
		return len(vs)
	}
	return 0
}

/////////////// Region H Private

func (h H) clone() H {
	newH := make(H)
	for k, v := range h {
		newH[k] = v
	}
	return newH
}

func (h H) overrideKey(key string, value ...string) {
	switch len(value) {
	case 0:
		delete(h, key)
	case 1:
		h[key] = value[0]
	default:
		panic("should not enter this branch")
	}
}

/////////////// Region H Public

/////////////// Region H Interface

func (h H) toHs() Hs {
	newHs := make(Hs)
	for k, v := range h {
		newHs[k] = []string{v}
	}
	return newHs
}

func (h H) WithKV(key string, value ...string) IHeader {
	if len(value) <= 1 {
		newH := h.clone()
		newH.overrideKey(key, value...)
		return newH
	}
	newHs := make(Hs)
	return newHs.WithKV(key, value...)
}

func (h H) WithKVAppend(key string, value ...string) IHeader {
	if len(value) == 0 {
		return h.clone()
	}
	if h.CountOf(key) == 0 && len(value) == 1 {
		newH := h.clone()
		newH.overrideKey(key, value...)
		return newH
	}
	newHs := h.toHs()
	return newHs.WithKVAppend(key, value...)
}

func (h H) WithH(headers IHeader) IHeader {
	countMax := 0
	headers.foreach(func(key string, _ ...string) {
		c := headers.CountOf(key)
		if c > countMax {
			countMax = c
		}
	})
	if countMax <= 1 {
		newH := h.clone()
		headers.foreach(func(key string, vs ...string) {
			if len(vs) == 1 {
				newH[key] = vs[0]
			}
		})
		return newH
	}

	return h.toHs().WithH(headers)
}

func (h H) WithHAppend(headers IHeader) IHeader {
	countMax := 0
	headers.foreach(func(key string, _ ...string) {
		c := headers.CountOf(key) + h.CountOf(key)
		if c > countMax {
			countMax = c
		}
	})
	if countMax <= 1 {
		newH := h.clone()
		headers.foreach(func(key string, vs ...string) {
			if len(vs) == 1 {
				newH[key] = vs[0]
			}
		})
		return newH
	}
	return h.toHs().WithHAppend(headers)
}

/////////////// Region HS Private

func copyStrLst(slice []string, values ...[]string) []string {
	count := len(slice)
	total := count
	for _, vs := range values {
		total += len(vs)
	}

	ret := make([]string, total)
	copy(ret, slice)
	for _, vs := range values {
		copy(ret[count:], vs)
		count += len(vs)
	}
	return ret
}

func (hs Hs) cloneValues(key string, appends ...string) []string {
	if vs, ok := hs[key]; ok {
		return copyStrLst(vs, appends)
	}
	return copyStrLst(appends)
}

func (hs Hs) clone() Hs {
	newHs := make(Hs)
	for k := range hs {
		newHs[k] = hs.cloneValues(k)
	}
	return newHs
}

/////////////// Region HS Interface

func (hs Hs) WithKV(key string, value ...string) IHeader {
	newHs := make(Hs)
	if len(value) > 0 {
		newHs[key] = value
	}

	for k := range hs {
		if k == key {
			continue
		}
		newHs[k] = hs.cloneValues(k)
	}
	return newHs
}

func (hs Hs) WithKVAppend(key string, values ...string) IHeader {
	if len(values) <= 0 {
		return hs.clone()
	}
	newHs := make(Hs)
	newHs[key] = hs.cloneValues(key, values...)
	for k := range hs {
		if k == key {
			continue
		}
		newHs[k] = hs.cloneValues(k)
	}
	return newHs
}

func (hs Hs) WithH(headers IHeader) IHeader {
	newHs := make(Hs)
	headers.foreach(func(key string, vsInput ...string) {
		newHs[key] = make([]string, len(vsInput))
		copy(newHs[key], vsInput)
	})
	hs.foreach(func(key string, vs ...string) {
		if _, ok := newHs[key]; ok {
			return
		}
		newHs[key] = copyStrLst(vs)
	})
	return newHs
}

func (hs Hs) WithHAppend(headers IHeader) IHeader {
	newHs := make(Hs)
	headers.foreach(func(key string, vsInput ...string) {
		if vs, ok := hs[key]; ok {
			newHs[key] = copyStrLst(vs, vsInput)
		} else {
			newHs[key] = copyStrLst(vsInput)
		}
	})

	hs.foreach(func(key string, vs ...string) {
		if _, ok := newHs[key]; ok {
			return
		}
		newHs[key] = copyStrLst(vs)
	})
	return newHs
}
