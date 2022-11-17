package otp

import (
	"context"

	"bitbucket.org/phoops/otp-stops-extractor/internal/core/entities"
	"github.com/hasura/go-graphql-client"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Client struct {
	logger        *zap.SugaredLogger
	baseURL       string
	graphqlClient *graphql.Client
}

func NewClient(logger *zap.SugaredLogger, baseURL string) (*Client, error) {
	if logger == nil || baseURL == "" {
		return nil, errors.New("all parametes must be non-nil")
	}

	graphqlClient := graphql.NewClient(baseURL, nil)
	logger = logger.With("component", "OTP client")

	return &Client{
		logger,
		baseURL,
		graphqlClient,
	}, nil
}

func (c *Client) GetStopsByBoundingBox(ctx context.Context, minLon, maxLon, minLat, maxLat float32) ([]*entities.Stop, error) {
	// query stopsByBB {
	// 	stopsByBbox(
	//    minLon: 11.801471,
	//    maxLon: 11.922997 ,
	//    minLat: 43.350637,
	//    maxLat: 43.503272
	// 	){ id lat lon code name }
	// }

	// Represent
	var stopsByBB struct {
		StopsByBbox []struct {
			ID   string
			Lat  float32
			Lon  float32
			Code string
			Name string
		} `graphql:"stopsByBbox(minLon: $minLon, maxLon: $maxLon, minLat: $minLat, maxLat: $maxLat)"`
	}

	variables := map[string]any{
		"minLon": minLon,
		"maxLon": maxLon,
		"minLat": minLat,
		"maxLat": maxLat,
	}

	c.logger.Debugw("About to query graphql", "endpoint", c.baseURL)
	err := c.graphqlClient.Query(ctx, &stopsByBB, variables)
	if err != nil {
		return nil, errors.Wrap(err, "graphql query to OTP failed")
	}

	c.logger.Debug("About to convert to entities")
	res := []*entities.Stop{}
	for _, stop := range stopsByBB.StopsByBbox {
		res = append(res, &entities.Stop{
			Code: stop.Code,
			Name: stop.Name,
			Lat:  stop.Lat,
			Lon:  stop.Lon,
		})
	}

	return res, nil
}
