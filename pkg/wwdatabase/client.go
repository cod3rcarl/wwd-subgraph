package client

import (
	"fmt"

	pb "github.com/cod3rcarl/wwd-protorepo-wwdatabase/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceInterface interface {
	Close()
	WWdatabaseInterface
}

type WWdatabaseInterface interface{}

type Client struct {
	logger           *zap.Logger
	cc               *grpc.ClientConn
	wwdatabaseServer pb.WwdatabaseClient
}

func NewClient(l *zap.Logger, cfg Config) (*Client, error) {
	ccon, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, ErrUnavailable
	}
	return &Client{
		cc:               ccon,
		logger:           l,
		wwdatabaseServer: pb.NewWwdatabaseClient(ccon),
	}, nil
}

func (c *Client) Close() {
	if err := c.cc.Close(); err != nil {
		c.logger.Error("failed to close client", zap.Error(err))
	}
}
