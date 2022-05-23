package resolvers

import (
	"time"

	"github.com/cod3rcarl/wwd-subgraph/pkg/client"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

const timeout = time.Minute

type Resolver struct {
	Logger *zap.Logger
	Server *client.Client
}
