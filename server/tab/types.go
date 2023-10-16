package tab

import "time"

type TimeBandEventResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	Events []Event `json:"events"`
}

type Event struct {
	Name      string    `json:"name"`
	Markets   []Market  `json:"markets"`
	StartTime time.Time `json:"startTime"`
	Teams     []Team    `json:"teams"`
	Meeting   Meeting   `json:"meeting"`
}

type Meeting struct {
	Name string `json:"name"`
}
type Team struct {
	Name string `json:"name"`
}

type Market struct {
	HandicapValue interface{} `json:"handicapValue"` // can be float or null
	GroupCode     string      `json:"groupCode"`
	Outcomes      []Outcome   `json:"outcomes"`
}

type Outcome struct {
	DisplayOrder int     `json:"displayOrder"`
	Prices       []Price `json:"prices"`
}

type Price struct {
	Decimal float64 `json:"decimal"`
}
