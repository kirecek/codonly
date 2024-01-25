package providers

import (
	"context"

	"github.com/kirecek/codonly/internal/pkg/state"
)

type DataProvider interface {
	ListResources(context.Context) ([]state.Resource, error)
}
