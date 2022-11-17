package main

import (
	"context"

	"bitbucket.org/phoops/otp-stops-extractor/internal/core/usecase"
	"bitbucket.org/phoops/otp-stops-extractor/internal/infrastructure/config"
	ngsild "bitbucket.org/phoops/otp-stops-extractor/internal/infrastructure/ngsi-ld"
	"bitbucket.org/phoops/otp-stops-extractor/internal/infrastructure/otp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	sourLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	logger := sourLogger.Sugar()

	conf, err := config.LoadStopsExtractorConfig()
	if err != nil {
		errMsg := errors.Wrap(err, "cannot read configuration").Error()
		logger.Fatal(errMsg)
	}

	otpClient, err := otp.NewClient(logger, conf.OtpURL)
	if err != nil {
		errMsg := errors.Wrap(err, "cannot instantiate OTP client").Error()
		logger.Fatal(errMsg)
	}

	contextBrokerClient, err := ngsild.NewClient(
		logger,
		conf.BrokerURL,
	)
	if err != nil {
		errMsg := errors.Wrap(err, "cannot instantiate context broker client").Error()
		logger.Fatal(errMsg)
	}

	syncStops, err := usecase.NewSyncStops(
		logger,
		otpClient,
		contextBrokerClient,
	)
	if err != nil {
		errMsg := errors.Wrap(err, "cannot instantiate syncStops").Error()
		logger.Fatal(errMsg)
	}

	// Sync!
	err = syncStops.Execute(
		context.Background(),
		conf.BoundingBoxMinLon,
		conf.BoundingBoxMaxLon,
		conf.BoundingBoxMinLat,
		conf.BoundingBoxMaxLat,
	)
	if err != nil {
		errMsg := errors.Wrap(err, "cannot sync stops on context broker").Error()
		logger.Fatal(errMsg)
	}
}
