package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type StopsExtractorConfig struct {
	OtpURL            string  `required:"true" split_words:"true"`
	BrokerURL         string  `required:"true" split_words:"true"`
	BoundingBoxMinLon float32 `required:"true" split_words:"true"`
	BoundingBoxMaxLon float32 `required:"true" split_words:"true"`
	BoundingBoxMinLat float32 `required:"true" split_words:"true"`
	BoundingBoxMaxLat float32 `required:"true" split_words:"true"`
}

func (s StopsExtractorConfig) String() string {
	return fmt.Sprintf(`
OtpURL: %v
BrokerURL: %s
`,
		s.OtpURL,
		s.BrokerURL,
	)
}

func LoadStopsExtractorConfig() (*StopsExtractorConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("could not load configuration from .env file: %v", err)
	}
	var c StopsExtractorConfig
	err = envconfig.Process("", &c)
	if err != nil {
		return nil, err
	}
	log.Printf("Loaded configuration%+s", c)
	return &c, nil
}
