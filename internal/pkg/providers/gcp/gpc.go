package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/kirecek/codonly/internal/pkg/providers"
	"github.com/kirecek/codonly/internal/pkg/state"
	container "google.golang.org/api/container/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

// GCP represents a provider for google cloud platform resources.
type GCP struct {
	projectID string
}

// NewGoogleProvider returns a GCP struct.
func NewGoogleProvider(projectID string) providers.DataProvider {
	return &GCP{projectID: projectID}
}

// ListResources wraps a data fetched from google cloud into codonly resource objects.
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

	clusters, err := g.listContainerClusters(ctx)
	if err != nil {
		return nil, err
	}
	res = append(res, clusters...)

	return res, nil
}

func (g *GCP) listContainerClusters(ctx context.Context) ([]state.Resource, error) {
	svc, err := container.NewService(ctx, option.WithScopes(container.CloudPlatformScope))
	if err != nil {
		return nil, err
	}

	resp, err := svc.Projects.Locations.Clusters.List(fmt.Sprintf("projects/%s/locations/-", g.projectID)).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	resources := make([]state.Resource, 0)
	for _, cluster := range resp.Clusters {
		resources = append(resources, state.Resource{
			IDValue:     cluster.Id,
			Type:        "google_container_cluster",
			DisplayName: cluster.Name,
		})
	}

	return resources, nil
}

func (g *GCP) listCloudSQLInstances(ctx context.Context) ([]state.Resource, error) {
	svc, err := sqladmin.NewService(ctx, option.WithScopes(sqladmin.SqlserviceAdminScope))
	if err != nil {
		return nil, err
	}

	resp, err := svc.Instances.List(g.projectID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	resources := make([]state.Resource, 0)

	for _, instance := range resp.Items {
		resources = append(resources, state.Resource{
			IDValue:     instance.Name,
			Type:        "google_sql_database_instance",
			DisplayName: instance.Name,
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
			IDKey:       &idKey,
			IDValue:     battrs.Name,
			Type:        "google_storage_bucket",
			DisplayName: battrs.Name,
		})
	}

	return resources, nil
}
