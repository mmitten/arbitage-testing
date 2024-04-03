package theoddsinterface

import (
	"encoding/json"
	"net/http"
	"time"
)

type Event []struct {
	ID           string      `json:"id"`
	SportKey     string      `json:"sport_key"`
	SportTitle   string      `json:"sport_title"`
	CommenceTime time.Time   `json:"commence_time"`
	HomeTeam     string      `json:"home_team"`
	AwayTeam     string      `json:"away_team"`
	Bookmakers   []Bookmaker `json:"bookmakers"`
}

type Bookmaker struct {
	Key        string    `json:"key"`
	Title      string    `json:"title"`
	LastUpdate time.Time `json:"last_update"`
	Markets    []Market  `json:"markets"`
}

type Market struct {
	Key        string    `json:"key"`
	LastUpdate time.Time `json:"last_update"`
	Outcomes   []Outcome `json:"outcomes"`
}

type Outcome struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func GetOdds(sport_key string, market string) (*Event, error) {
	data := Event{}

	cfg, errConfig := GetConfig()
	if errConfig != nil {
		return &data, errConfig
	}

	urlParts := GetUrlStruct()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", urlParts.BaseUrl+urlParts.SportsEndpoint+"/"+sport_key+urlParts.OddsEndpoint, nil)

	q := req.URL.Query()
	q.Add("apiKey", cfg.ApiKey)
	q.Add("regions", "us")
	q.Add("markets", market)
	req.URL.RawQuery = q.Encode()

	res, errResult := client.Do(req)
	if errResult != nil {
		return &data, errResult
	}
	defer res.Body.Close()
	if res.StatusCode > 305 {
		return &data, nil
	}
	if errJson := json.NewDecoder(res.Body).Decode(&data); errJson != nil {
		return &data, errJson
	}
	return &data, nil
}

type BestOdds struct {
	Event      string
	Bookmakers map[string]BestBookmaker
}

type BestBookmaker struct {
	Name    string
	Market  string
	Outcome string
	Price   float64
}

func (e Event) GetBestOdds() []BestOdds {
	eventBestOdds := []BestOdds{}
	for _, event := range e {
		bestOdds := BestOdds{
			Event:      event.AwayTeam + " at " + event.HomeTeam + " in " + event.SportTitle,
			Bookmakers: make(map[string]BestBookmaker),
		}
		for bId, bookmaker := range event.Bookmakers {
			if len(bookmaker.Markets) < 1 {
				continue
			}
			market := bookmaker.Markets[0]
			for _, outcome := range market.Outcomes {
				bb := BestBookmaker{
					Name:    bookmaker.Title,
					Market:  market.Key,
					Outcome: outcome.Name,
					Price:   outcome.Price,
				}
				if bId == 0 {
					bestOdds.Bookmakers[outcome.Name] = bb
					continue
				}
				if outcome.Price > bestOdds.Bookmakers[outcome.Name].Price {
					bestOdds.Bookmakers[outcome.Name] = bb
				}
			}
		}
		eventBestOdds = append(eventBestOdds, bestOdds)
	}
	return eventBestOdds
}
