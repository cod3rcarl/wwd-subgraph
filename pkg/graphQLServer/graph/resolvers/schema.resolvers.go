package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/cod3rcarl/wwd-subgraph/pkg/graphQLServer/graph/generated"
	"github.com/cod3rcarl/wwd-subgraph/pkg/graphQLServer/graph/model"
)

func (r *mutationResolver) CreateChampion(ctx context.Context, input model.CreateChampionInput) (*model.CreateChampionPayload, error) {
	pbChamp, err := r.Server.CreateChampion(ctx, &input)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating champion")
	}

	return pbChamp, nil
}

func (r *mutationResolver) DeleteChampion(ctx context.Context, input model.DeleteChampionInput) (*model.DeleteChampionPayload, error) {
	pbChamp, err := r.Server.DeleteChampion(ctx, &input)
	if err != nil {
		return nil, errors.Wrap(err, "Error deleting champion")
	}

	return pbChamp, nil
}

func (r *queryResolver) Champion(ctx context.Context, filter *model.ChampionFilterInput) (*model.ChampionPayload, error) {
	if filter == nil {
		return &model.ChampionPayload{
			Errors: []model.NewError{
				model.ChampionNoResultsReturned{
					Message: "No Champs Returned",
					Path: []string{
						"Champion Query",
					},
				},
			},
		}, nil
	}

	if filter.PreviousChampion != nil {

		pbChamp, err := r.Server.GetChampionByOrderNumber(ctx, int32(*filter.PreviousChampion-1))
		if err != nil {
			return &model.ChampionPayload{
				Errors: []model.NewError{
					model.ChampionNoResultsReturned{
						Message: "No Champs Returned",
						Path: []string{
							"Champion Query",
						},
					},
				},
			}, nil
		}

		return &model.ChampionPayload{
			Champion: pbChamp,
			Errors:   nil,
		}, nil
	}

	if filter.CurrentChampion != nil {
		pbChamp, err := r.Server.GetCurrentChampion(ctx, *filter.CurrentChampion)
		if err != nil {
			return nil, errors.Wrap(err, "Error retrieving champion")
		}
		return &model.ChampionPayload{
			Champion: pbChamp,
			Errors:   nil,
		}, nil
	}

	if filter.Date != nil {
		d, err := time.Parse("2006-01-02", *filter.Date)
		if err != nil {
			return nil, errors.Wrap(err, "Error retrieving champion")
		}

		today := time.Now()

		if d.Unix() > today.Unix() {
			return nil, errors.New("Either I haven't updated the list, or I can't tell the future")
		}
		firstDate, _ := time.Parse("2006-01-02", "1963-04-10")

		if d.Unix() < firstDate.Unix() {
			return nil, errors.New("First champion on record won 10th April 1963")
		}

		pbChamp, err := r.Server.GetChampionByDate(ctx, d)
		if err != nil {
			return nil, errors.Wrap(err, "Error retrieving champion")
		}

		return &model.ChampionPayload{
			Champion: pbChamp,
			Errors:   nil,
		}, nil
	}

	return nil, errors.New("Please provide an input to search")
}

func (r *queryResolver) Champions(ctx context.Context, filter *model.ChampionsFilterInput) (*model.ChampionsPayload, error) {
	if filter == nil {
		pbChamps, err := r.Server.GetChampions(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "Error getting champions")
		}

		return pbChamps, nil
	}

	if filter.Year != nil {
		sd, err := time.Parse("2006-01-02", *filter.Year.Start)
		if err != nil {
			return nil, errors.Wrap(err, "Error retrieving champion")
		}
		ed, err := time.Parse("2006-01-02", *filter.Year.End)
		if err != nil {
			return nil, errors.Wrap(err, "Error retrieving champion")
		}

		today := time.Now()
		firstDate, _ := time.Parse("2006-01-02", "1963-04-10")
		if sd.Unix() > today.Unix() {
			return nil, errors.New("Either I haven't updated the list, or I can't tell the future")
		}

		if sd.Unix() < firstDate.Unix() && ed.Unix() < firstDate.Unix() {
			return nil, errors.New("First champion on record won 10th April 1963")
		}

		if ed.Unix() > today.Unix() {
			ed = today
		}

		pbChamps, err := r.Server.GetChampionsByYear(ctx, sd, ed)
		if err != nil {
			return nil, errors.Wrap(err, "Error retrieving champion")
		}

		return pbChamps, nil
	}

	if filter.Show != nil {
		pbChamps, err := r.Server.GetChampionsByShow(ctx, *filter.Show)
		if err != nil {
			return nil, errors.Wrap(err, "Error retrieving champion")
		}

		return pbChamps, nil
	}

	if filter.TitleHolder != nil {
		pbChamps, err := r.Server.GetChampionReignsByName(ctx, *filter.TitleHolder)
		if err != nil {
			return nil, errors.Wrap(err, "Error retrieving champion")
		}

		return pbChamps, nil
	}

	return nil, errors.New("No idea, check the code")
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
)
