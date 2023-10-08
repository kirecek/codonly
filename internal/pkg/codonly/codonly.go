package codonly

import (
	"context"
	"fmt"

	"github.com/kirecek/codonly/internal/pkg/providers"
	"github.com/kirecek/codonly/internal/pkg/state"
)

type Codonly struct {
	state state.State
	data  providers.DataProvider
}

func New(state state.State, data providers.DataProvider) *Codonly {
	return &Codonly{
		state: state,
		data:  data,
	}
}

func (c *Codonly) Run(ctx context.Context) error {
	resources, err := c.data.ListResources(ctx)
	if err != nil {
		return err
	}

	for _, resource := range resources {
		if !c.state.Contains(&resource) {
			fmt.Printf("'%s/%s' was not found in the provided state\n", resource.Type, resource.DisplayName)
		}
	}

	return nil
}
