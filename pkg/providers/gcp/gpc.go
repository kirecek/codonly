package gcp

import (
	"context"
	"strings"

	"github.com/kirecek/codonly/pkg/codonly"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

type GCP struct {
	projectID string
}

func NewGoogleProvider(projectID string) *GCP {
	return &GCP{projectID: projectID}
}

func (g *GCP) ListResources(ctx context.Context) ([]codonly.CodonlyResource, error) {
	return g.listServiceAccounts(ctx)
}

func (g *GCP) listServiceAccounts(ctx context.Context) ([]codonly.CodonlyResource, error) {
	iamService, err := iam.NewService(ctx, option.WithScopes(iam.CloudPlatformScope))
	if err != nil {
		return nil, err
	}

	accounts, err := iamService.Projects.ServiceAccounts.List("projects/" + g.projectID).Do()
	if err != nil {
		return nil, err
	}

	resources := make([]codonly.CodonlyResource, 0)

	for _, account := range accounts.Accounts {
		if strings.HasSuffix(account.DisplayName, "Engine default service account") {
			continue
		}

		resources = append(resources, codonly.CodonlyResource{
			IDKey:   "id",
			IDValue: account.Name,
			Type:    "google_service_account",
		})
	}

	return resources, nil
}
