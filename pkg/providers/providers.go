package providers

import (
	"context"

	"github.com/kirecek/codonly/pkg/state"
)

type Provider interface {
	ListResources(context.Context) ([]state.Resource, error)
}
