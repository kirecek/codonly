package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kirecek/codonly/pkg/providers"
	"github.com/kirecek/codonly/pkg/providers/gcp"
	"github.com/kirecek/codonly/pkg/state"
)

func main() {
	var (
		stateProvider state.StateProvider
		dataProvider  providers.Provider
		err           error
	)

	statePath := flag.String("state", "", "Path a terraform output")
	flag.Parse()

	if statePath != nil && *statePath != "" {
		stateProvider, err = state.NewTerraformProviderFromStateOutput(*statePath)
	} else {
		stateProvider, err = state.NewTerraformProvider()
	}

	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	dataProvider = gcp.NewGoogleProvider(os.Getenv("GOOGLE_PROJECT"))

	resources, err := dataProvider.ListResources(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	for _, resource := range resources {
		if !stateProvider.Contains(&resource) {
			fmt.Printf("'%s/%s' was not found in state\n", resource.Type, resource.DisplayName)
		}
	}
}
