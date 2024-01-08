package types

type Job struct {
	Name      string   `json:"name"`
	Link      string   `json:"link"`
	Locations []string `json:"locations"`
	Process   string   `json:"process"`
}

type Set map[string]struct{}

func NewSet() Set {
	return make(Set)
}

func (s Set) Add(element string) {
	s[element] = struct{}{}
}

func (s Set) Exists(element string) bool {
	_, ok := s[element]
	return ok
}

func (s Set) Remove(element string) {
	delete(s, element)
}
