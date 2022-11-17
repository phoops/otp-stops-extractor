package usecase

import (
	"context"
	"fmt"

	"bitbucket.org/phoops/otp-stops-extractor/internal/core/entities"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type StopsFetcher interface {
	GetStopsByBoundingBox(ctx context.Context, minLon, maxLon, minLat, maxLat float32) ([]*entities.Stop, error)
}

type StopsPersistor interface {
	WriteStopsBatch(ctx context.Context, stops []*entities.Stop) error
}

type SyncStopsUC struct {
	logger    *zap.SugaredLogger
	fetcher   StopsFetcher
	persistor StopsPersistor
}

func NewSyncStops(
	logger *zap.SugaredLogger,
	fetcher StopsFetcher,
	persistor StopsPersistor,
) (*SyncStopsUC, error) {
	if logger == nil || fetcher == nil || persistor == nil {
		return nil, errors.New("all parameters must be non-nil")
	}
	logger = logger.With("usecase", "SyncStops")

	return &SyncStopsUC{
		logger,
		fetcher,
		persistor,
	}, nil
}

func (u *SyncStopsUC) Execute(ctx context.Context, minLon, maxLon, minLat, maxLat float32) error {
	u.logger.Info("Running Stops Synchronization")
	stops, err := u.fetcher.GetStopsByBoundingBox(ctx, minLon, maxLon, minLat, maxLat)
	if err != nil {
		u.logger.Errorw("can't read stops", "error", err, "bounding box", fmt.Sprintf("(%f,%f),(%f,%f)", minLon, minLat, maxLon, maxLat))
		return errors.Wrap(err, "can't read stops")
	}
	u.logger.Debugw("stops read", "fetched", len(stops))

	err = u.persistor.WriteStopsBatch(ctx, stops)
	if err != nil {
		u.logger.Errorw("can't write stops", "error", err)
		return errors.Wrap(err, "can't write stops")
	}
	u.logger.Infow("stops written", "size", len(stops))
	return nil
}
