package main

import (
	"context"
	"fmt"
	"log"
)

type CodonlyResource struct {
	IDKey   string
	IDValue string
	Type    string
}

type CodonlyApp struct {
	state *TerraformProvider
}

func NewCodonly() (*CodonlyApp, error) {
	state, err := NewTerraformProvider()
	if err != nil {
		return nil, err
	}

	return &CodonlyApp{
		state: state,
	}, nil
}

func (c *CodonlyApp) Contains(r *CodonlyResource) bool {
	for _, resource := range c.state.state.Values.Module.Resources {
		if resource.Type == r.Type && resource.Values[r.IDKey].(string) == r.IDValue {
			return true
		}
	}

	return false
}

func main() {
	codonly, err := NewCodonly()
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()

	gcp := GCP{projectID: "kirecek"}

	resources, err := gcp.listServiceAccounts(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	for _, resource := range resources {
		if !codonly.Contains(&resource) {
			fmt.Printf("'%s' is not present in provided state\n", resource.IDValue)
		}
	}

}
