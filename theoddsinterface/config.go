package theoddsinterface

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type InterfaceConfig struct {
	ApiKey string
}

type URLConstants struct {
	BaseUrl        string
	SportsEndpoint string
	OddsEndpoint   string
}

func GetConfig() (*InterfaceConfig, error) {
	cfg := &InterfaceConfig{}
	if errLoadEnv := godotenv.Load(); errLoadEnv != nil {
		return cfg, errLoadEnv
	}

	opts := env.Options{UseFieldNameByDefault: true}

	if errParseEnv := env.ParseWithOptions(cfg, opts); errParseEnv != nil {
		return cfg, errParseEnv
	}

	return cfg, nil
}

func GetUrlStruct() *URLConstants {
	url := &URLConstants{
		BaseUrl:        "https://api.the-odds-api.com",
		SportsEndpoint: "/v4/sports",
		OddsEndpoint:   "/odds",
	}
	return url
}
