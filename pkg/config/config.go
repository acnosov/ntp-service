package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"time"
)

type Config struct {
	Service struct {
		Name       string        `default:"ntp-service" json:"name"`
		SyncPeriod time.Duration `default:"1m" json:"sync_period"`
	} `json:"service"`
	Zap struct {
		//debug, info, warn, error, fatal, panic
		Level string `default:"info" json:"level"`
		//console, json
		Encoding string `default:"console" json:"encoding"`
		//disable, short, full
		Caller string `default:"disable" json:"caller"`
	} `json:"zap"`
}

func NewConfig() *Config {
	var cfg Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		//SkipFlags:          true,
		AllFieldRequired:   true,
		AllowUnknownFlags:  true,
		AllowUnknownEnvs:   true,
		AllowUnknownFields: true,
		AllowDuplicates:    true,
		SkipEnv:            false,
		FileFlag:           "config",
		FailOnFileNotFound: false,
		MergeFiles:         true,
		Files:              []string{"config.yaml", "ntp-service.yaml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})
	err := loader.Load()
	if err != nil {
		panic(err)
	}

	return &cfg
}
