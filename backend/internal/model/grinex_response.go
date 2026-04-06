package model

type GrinexResponse struct {
	Timestamp int64      `json:"timestamp"`
	Asks      []RateItem `json:"asks"`
	Bids      []RateItem `json:"bids"`
}