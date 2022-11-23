package ngsild

import (
	"context"
	"fmt"

	"bitbucket.org/phoops/otp-stops-extractor/internal/core/entities"
	"github.com/philiphil/geojson"
	"github.com/phoops/ngsi-gold/client"
	"github.com/phoops/ngsi-gold/ldcontext"
	"github.com/phoops/ngsi-gold/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Client struct {
	logger       *zap.SugaredLogger
	baseURL      string
	ngsiLdClient *client.NgsiLdClient
}

func NewClient(logger *zap.SugaredLogger, baseURL string) (*Client, error) {
	if logger == nil {
		return nil, errors.New("all parameters must be non-nil")
	}
	logger = logger.With("component", "NGSI-LD client")
	ngsiLdClient, err := client.New(
		client.SetURL(baseURL),
	)
	if err != nil {
		return nil, errors.Wrap(err, "can't instantiate ngsi-ld client")
	}

	return &Client{
		logger,
		baseURL,
		ngsiLdClient,
	}, nil
}

func (c *Client) WriteStopsBatch(ctx context.Context, stops []*entities.Stop) error {
	payload := []*client.EntityWithContext{}
	for _, s := range stops {
		e := stopToBrokerEntity(s)
		payload = append(payload, &client.EntityWithContext{
			LdCtx:  &ldcontext.DefaultContext,
			Entity: e,
		})
		c.logger.Debugw("added entity to batch", "id", e.ID, "gtfs IDs", s.GtfsIDs, "entity GTFS IDs", e.Properties["gtfsIDs"].Value)
	}

	err := c.ngsiLdClient.BatchUpsertEntities(ctx, payload, client.UpsertSetUpdateMode)
	if err != nil {
		c.logger.Errorw("can't update entities", "err", err)
		return errors.Wrap(err, "can't update entities")
	}
	return nil
}

func stopToBrokerEntity(e *entities.Stop) *model.Entity {
	id := fmt.Sprintf("urn:ngsi-ld:GtfsStop:%s", e.Code)
	location := geojson.NewPointGeometry([]float64{float64(e.Lon), float64(e.Lat)})

	return &model.Entity{
		ID:   id,
		Type: "GtfsStop",
		Properties: model.Properties{
			"name": model.Property{
				Value: e.Name,
			},
			"code": model.Property{
				Value: e.Code,
			},
			"gtfsIDs": model.Property{
				Value: e.GtfsIDs,
			},
		},
		// Build multi-value relationship from e.Agencies
		// Need library to support multi-value relationship
		// Relationships: model.Relationships{
		// 	"operatedBy": []
		// },
		Location: &model.GeoProperty{
			Value: location,
		},
	}

}
