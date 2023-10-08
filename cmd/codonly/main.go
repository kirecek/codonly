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
		sp  state.StateProvider
		dp  providers.Provider
		err error
	)

	statePath := flag.String("state", "", "Path a terraform output")
	flag.Parse()

	if statePath != nil && *statePath != "" {
		sp, err = state.NewTerraformProviderFromStateOutput(*statePath)
	} else {
		sp, err = state.NewTerraformProvider()
	}

	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	dp = gcp.NewGoogleProvider(os.Getenv("GOOGLE_PROJECT"))

	resources, err := dp.ListResources(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	for _, resource := range resources {
		if !sp.Contains(&resource) {
			fmt.Printf("'%s' is not present in provided state\n", resource.IDValue)
		}
	}
}
