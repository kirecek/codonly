package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kirecek/codonly/pkg/codonly"
	"github.com/kirecek/codonly/pkg/providers/gcp"
)

func main() {
	codonly, err := codonly.NewCodonly()
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()

	provider := gcp.NewGoogleProvider(os.Getenv("GOOGLE_PROJECT"))

	resources, err := provider.ListResources(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	for _, resource := range resources {
		if !codonly.Contains(&resource) {
			fmt.Printf("'%s' is not present in provided state\n", resource.IDValue)
		}
	}
}
