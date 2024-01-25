package state

// StateProvider is simple interrace that interacts on of infrastructure
// management tools like terraform.
type StateProvider interface {
	Contains(*Resource) bool
}

// Resource is a unique representation of a resource present in the state provider.
type Resource struct {
	IDKey       *string
	IDValue     string
	Type        string
	DisplayName string
}
