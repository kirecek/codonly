package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
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
	res := make([]state.Resource, 0)

	instances, err := g.listCloudSQLInstances(ctx)
	if err != nil {
		return nil, err
	}
	res = append(res, instances...)

	buckets, err := g.listBuckets(ctx)
	if err != nil {
		return nil, err
	}
	res = append(res, buckets...)

	return res, nil
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
			IDValue: instance.Name,
			Type:    "google_sql_database_instance",
		})
	}

	return resources, nil
}

func (g *GCP) listBuckets(ctx context.Context) ([]state.Resource, error) {
	svc, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}

	idKey := "name"

	resources := make([]state.Resource, 0)
	it := svc.Buckets(ctx, g.projectID)

	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		resources = append(resources, state.Resource{
			IDKey:   &idKey,
			IDValue: battrs.Name,
			Type:    "google_storage_bucket",
		})
	}

	return resources, nil
}
