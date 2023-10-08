package codonly

import "github.com/kirecek/codonly/pkg/state"

type CodonlyResource struct {
	IDKey   string
	IDValue string
	Type    string
}

type CodonlyApp struct {
	state *state.TerraformProvider
}

func NewCodonly() (*CodonlyApp, error) {
	state, err := state.NewTerraformProvider()
	if err != nil {
		return nil, err
	}

	return &CodonlyApp{
		state: state,
	}, nil
}

func (c *CodonlyApp) Contains(r *CodonlyResource) bool {
	for _, resource := range c.state.Data.Values.Module.Resources {
		if resource.Type == r.Type && resource.Values[r.IDKey].(string) == r.IDValue {
			return true
		}
	}

	return false
}
