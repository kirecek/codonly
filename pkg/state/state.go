package state

type StateProvider interface {
	Contains(*Resource) bool
}

type Resource struct {
	IDKey   *string
	IDValue string
	Type    string
}
