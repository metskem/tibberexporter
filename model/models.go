package model

import (
	"fmt"
	"time"
)

type TibberData struct {
	Total    float64   `json:"total"`
	Energy   float64   `json:"energy"`
	Tax      float64   `json:"tax"`
	StartsAt time.Time `json:"startsAt"`
}

const TimeFormat = "2006-01-02T15:04Z07:00"

func (data TibberData) String() string {
	return fmt.Sprintf("PriceTotal:%.4f, PriceEnergy:%.4f, PriceTax:%.4f, StartsAt:%s", data.Total, data.Energy, data.Tax, data.StartsAt.Format(TimeFormat))
}

type TibberResponse struct {
	Data struct {
		Viewer struct {
			Homes []struct {
				CurrentSubscription struct {
					PriceInfo struct {
						Today    []TibberData `json:"today"`
						Tomorrow []TibberData `json:"tomorrow"`
					} `json:"priceInfo"`
				} `json:"currentSubscription"`
			} `json:"homes"`
		} `json:"viewer"`
	} `json:"data"`
}
