package gcp

import (
	"context"

	"google.golang.org/api/option"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"

	"github.com/kirecek/codonly/pkg/state"
)

type GCP struct {
	projectID string
}

func NewGoogleProvider(projectID string) *GCP {
	return &GCP{projectID: projectID}
}

func (g *GCP) ListResources(ctx context.Context) ([]state.Resource, error) {
	return g.listCloudSQLInstances(ctx)
}

func (g *GCP) listCloudSQLInstances(ctx context.Context) ([]state.Resource, error) {
	svc, err := sqladmin.NewService(ctx, option.WithScopes(sqladmin.SqlserviceAdminScope))
	if err != nil {
		return nil, err
	}

	resp, err := svc.Instances.List(g.projectID).Do()
	if err != nil {
		return nil, err
	}

	resources := make([]state.Resource, 0)

	for _, instance := range resp.Items {
		resources = append(resources, state.Resource{
			IDKey:   "id",
			IDValue: instance.Name,
			Type:    "google_sql_database_instance",
		})
	}

	return resources, nil

}
