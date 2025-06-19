package vault

type stringSet map[string]struct{}

func (s stringSet) Has(val string) bool {
	_, ok := s[val]
	return ok
}

func (s stringSet) Add(val string) {
	s[val] = struct{}{}
}

func (s stringSet) Delete(val string) {
	delete(s, val)
}
