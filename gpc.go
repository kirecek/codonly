package main

import (
	"context"
	"strings"

	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

type GCP struct {
	projectID string
}

func (g *GCP) listServiceAccounts(ctx context.Context) ([]CodonlyResource, error) {
	iamService, err := iam.NewService(ctx, option.WithScopes(iam.CloudPlatformScope))
	if err != nil {
		return nil, err
	}

	accounts, err := iamService.Projects.ServiceAccounts.List("projects/" + g.projectID).Do()
	if err != nil {
		return nil, err
	}

	resources := make([]CodonlyResource, 0)

	for _, account := range accounts.Accounts {
		if strings.HasSuffix(account.DisplayName, "Engine default service account") {
			continue
		}

		resources = append(resources, CodonlyResource{
			IDKey:   "id",
			IDValue: account.Name,
			Type:    "google_service_account",
		})
	}

	return resources, nil
}
