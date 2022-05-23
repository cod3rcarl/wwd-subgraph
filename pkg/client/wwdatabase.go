package client

import (
	"context"
	"time"

	pb "github.com/cod3rcarl/wwd-protorepo-wwdatabase/v1"
	"github.com/cod3rcarl/wwd-subgraph/pkg/graphQLServer/graph/model"
	"github.com/pkg/errors"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Client) CreateChampion(ctx context.Context, input *model.CreateChampionInput) (*model.CreateChampionPayload, error) {
	dateWon, err := time.Parse("2006-01-02", *input.DateWon)
	if err != nil {
		return nil, ErrDateIncorrect
	}
	pbChamp, err := s.wwdatabaseServer.AddChampion(ctx, &pb.NewChampionData{
		TitleHolder: input.TitleHolder,
		DateWon:     &timestamppb.Timestamp{Seconds: dateWon.Unix()},
		Show:        input.Show,
	})
	if err != nil {
		return nil, s.handleErr(err)
	}

	return &model.CreateChampionPayload{
		Success:  pbChamp.Success,
		Champion: pbChampionToModel(pbChamp.Champion),
	}, nil
}

func (s *Client) DeleteChampion(ctx context.Context, input *model.DeleteChampionInput) (
	*model.DeleteChampionPayload, error,
) {
	pbChamp, err := s.wwdatabaseServer.DeleteChampion(ctx, &pb.DeleteChampionRequest{
		Id: input.ID,
	})
	if err != nil {
		return nil, errors.Errorf("error in Deletehampion(): %v", err)
	}

	return &model.DeleteChampionPayload{
		Success: true,
		ID:      pbChamp.Id,
	}, nil
}

func (s *Client) GetChampions(ctx context.Context) (
	*model.ChampionsPayload, error,
) {
	pbChamps, err := s.wwdatabaseServer.GetChampions(ctx, &pb.GetChampionsRequest{})
	if err != nil {
		return nil, errors.Errorf("error in GetAllChampions(): %v", err)
	}

	return pbChampionsToModel(pbChamps), nil
}

func (s *Client) GetChampionsByShow(ctx context.Context, show string) (
	*model.ChampionsPayload, error,
) {
	pbChamps, err := s.wwdatabaseServer.GetChampionsByShow(ctx, &pb.GetChampionsByShowRequest{
		Show: show,
	})
	if err != nil {
		return nil, s.handleErr(err)
	}

	return pbChampionsToModel(pbChamps), nil
}

func (s *Client) GetChampionsByYear(ctx context.Context, start time.Time, end time.Time) (
	*model.ChampionsPayload, error,
) {
	pbChamps, err := s.wwdatabaseServer.GetChampionsByYear(ctx, &pb.GetChampionsByYearRequest{
		StartDate: &timestamppb.Timestamp{Seconds: start.Unix()},
		EndDate:   &timestamppb.Timestamp{Seconds: end.Unix()},
	})
	if err != nil {
		return nil, s.handleErr(err)
	}

	return pbChampionsToModel(pbChamps), nil
}

func (s *Client) GetChampionReignsByName(ctx context.Context, name string) (
	*model.ChampionsPayload, error,
) {
	pbChamps, err := s.wwdatabaseServer.GetChampionByName(ctx, &pb.GetChampionByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, s.handleErr(err)
	}

	return pbChampionsToModel(pbChamps), nil
}

func (s *Client) GetChampionByOrderNumber(ctx context.Context, tn int32) (
	*model.Champion, error,
) {
	pbChamp, err := s.wwdatabaseServer.GetChampionByOrderNumber(ctx, &pb.ChampionNumber{
		TitleHolderOrderNumber: tn,
	})
	if err != nil {
		return nil, s.handleErr(err)
	}

	return pbChampionToModel(pbChamp.Champion), nil
}

func (s *Client) GetCurrentChampion(ctx context.Context, cc bool) (
	*model.Champion, error,
) {
	pbChamp, err := s.wwdatabaseServer.GetCurrentChampion(ctx, &pb.GetCurrentChampionRequest{
		CurrentChampion: cc,
	})
	if err != nil {
		return nil, s.handleErr(err)
	}

	return pbChampionToModel(pbChamp.Champion), nil
}

func (s *Client) GetChampionByDate(ctx context.Context, date time.Time) (
	*model.Champion, error,
) {
	pbChamp, err := s.wwdatabaseServer.GetChampionByDate(ctx, &pb.GetChampionByDateRequest{
		Date: &timestamppb.Timestamp{Seconds: date.Unix()},
	})
	if err != nil {
		return nil, s.handleErr(err)
	}

	return pbChampionToModel(pbChamp.Champion), nil
}

func pbChampionsToModel(pbC *pb.ChampionsList) *model.ChampionsPayload {
	champions := []*model.Champion{}

	for _, champion := range pbC.Champions {
		titlehn := int(champion.TitleHolderNumber)
		titlehnon := int(champion.TitleHolderOrderNumber)

		dl := champion.DateLost.AsTime().String()

		c := &model.Champion{
			ID:                     champion.Id,
			TitleHolder:            champion.TitleHolder,
			TitleHolderNumber:      &titlehn,
			TitleHolderOrderNumber: &titlehnon,
			DateWon:                champion.DateWon.AsTime().String(),
			DateLost:               &dl,
			Show:                   champion.Show,
			CurrentChampion:        &champion.CurrentChampion,
			WrestlerID:             int(champion.WrestlerId),
		}
		champions = append(champions, c)
	}
	count := len(champions)

	return &model.ChampionsPayload{
		Champions:  champions,
		Errors:     nil,
		TotalCount: count,
	}
}

func pbChampionToModel(pbC *pb.Champion) *model.Champion {
	titlehn := int(pbC.TitleHolderNumber)
	titlehnon := int(pbC.TitleHolderOrderNumber)

	dl := pbC.DateLost.AsTime().String()

	return &model.Champion{
		ID:                     pbC.Id,
		TitleHolder:            pbC.TitleHolder,
		TitleHolderNumber:      &titlehn,
		TitleHolderOrderNumber: &titlehnon,
		DateWon:                pbC.DateWon.AsTime().String(),
		DateLost:               &dl,
		Show:                   pbC.Show,
		CurrentChampion:        &pbC.CurrentChampion,
		PreviousChampion:       &pbC.PreviousChampion,
		WrestlerID:             int(pbC.WrestlerId),
	}
}
