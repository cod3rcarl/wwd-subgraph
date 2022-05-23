package client

import (
	"errors"

	"github.com/cod3rcarl/wwd-subgraph/pkg/graphQLServer/graph/model"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternal            = errors.New("service internal error")
	ErrUnavailable         = errors.New("service unavailable")
	ErrUnknown             = errors.New("service error")
	ErrNoNameOrID          = errors.New("must provide name or ID")
	ErrDateIncorrect       = errors.New("date incorrectly formatted")
	ErrNoChampionsReturned = errors.New("no champions returned from query")
	ErrTableDoesNotExist   = errors.New("error in query, the table you are trying to access does not exist")
	ErrInvalidNullTime     = errors.New("timestamp provided is not valid")
)

func (c *Client) handleErr(err error) error {
	if s, ok := status.FromError(err); ok {
		if s.Code() == codes.NotFound {
			return err
		}
		if s.Code() == codes.Unavailable {
			c.logger.Error("service unavailable", zap.Error(err))

			return ErrUnavailable
		}
	}

	c.logger.Error("service error", zap.Error(err))

	return ErrUnknown
}

func CreateUserError(errCode string, message string, path string) model.NewError {
	var userError model.NewError
	var pathArray []string
	pathArray = append(pathArray, path)

	switch errCode {
	case `NoResultsReturned`:
		err := model.ChampionNoResultsReturned{
			Message: message,
			Path:    pathArray,
		}
		userError = err
	}

	return userError
}
