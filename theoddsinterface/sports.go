package theoddsinterface

import (
	"encoding/json"
	"net/http"
)

type SportStruct []struct {
	Key          string `json:"key"`
	Group        string `json:"group"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Active       bool   `json:"active"`
	HasOutrights bool   `json:"has_outrights"`
}

func GetSports() (*SportStruct, error) {
	data := SportStruct{}

	cfg, errConfig := GetConfig()
	if errConfig != nil {
		return &data, errConfig
	}

	urlParts := GetUrlStruct()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", urlParts.BaseUrl+urlParts.SportsEndpoint, nil)

	q := req.URL.Query()
	q.Add("apiKey", cfg.ApiKey)
	req.URL.RawQuery = q.Encode()

	res, errResult := client.Do(req)
	if errResult != nil {
		return &data, errResult
	}
	defer res.Body.Close()
	if errJson := json.NewDecoder(res.Body).Decode(&data); errJson != nil {
		return &data, errJson
	}
	return &data, nil
}

func (ss *SportStruct) CheckMarket(market string) ([]BestOdds, error) {
	betsToMake := []BestOdds{}
	for _, sport := range *ss {
		eventsData, errMarket := GetOdds(sport.Key, market)
		if errMarket != nil {
			return betsToMake, errMarket
		}
		bestOdds := eventsData.GetBestOdds()
		for _, odds := range bestOdds {
			if len(odds.Bookmakers) == 0 {
				continue
			}
			var sum float64 = 0.0
			for _, item := range odds.Bookmakers {
				sum = sum + 1/item.Price
			}
			if sum < 1.0 {
				betsToMake = append(betsToMake, odds)
			}
		}
	}
	return betsToMake, nil
}
