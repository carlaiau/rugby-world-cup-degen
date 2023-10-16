package types

import "time"

type DegenOutcome struct {
	DisplayOrder int     `json:"displayOrder"`
	Price        float64 `json:"price"`
}
type DegenMarket struct {
	Name       string         `json:"name"`
	Outcomes   []DegenOutcome `json:"outcomes"`
	Line       float64        `json:"line"`
	Bookmaker  string         `json:"bookmaker"`
	LastUpdate time.Time      `json:"last_update"`
}
type DegenEvent struct {
	StartTime time.Time              `json:"startTime"`
	Name      string                 `json:"name"`
	Markets   map[string]DegenMarket `json:"markets"`
	Home      string                 `json:"home"`
	Away      string                 `json:"away"`
}
type DegenEventForClient struct {
	StartTime time.Time     `json:"startTime"`
	Name      string        `json:"name"`
	Markets   []DegenMarket `json:"markets"`
	Home      string        `json:"home"`
	Away      string        `json:"away"`
}
