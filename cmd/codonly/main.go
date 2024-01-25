package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/kirecek/codonly/internal/pkg/codonly"
	"github.com/kirecek/codonly/internal/pkg/providers"
	"github.com/kirecek/codonly/internal/pkg/providers/gcp"
	"github.com/kirecek/codonly/internal/pkg/state"
	"github.com/kirecek/codonly/internal/pkg/state/terraform"
)

func main() {
	var (
		stateProvider state.State
		dataProvider  providers.DataProvider
		err           error
	)

	statePath := flag.String("state", "", "Path a terraform output")
	flag.Parse()

	if statePath != nil && *statePath != "" {
		stateProvider, err = terraform.NewFromStateOutput(*statePath)
	} else {
		stateProvider, err = terraform.New()
	}

	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	dataProvider = gcp.NewGoogleProvider(os.Getenv("GOOGLE_PROJECT"))

	codonly := codonly.New(stateProvider, dataProvider)
	codonly.Run(ctx)
}
