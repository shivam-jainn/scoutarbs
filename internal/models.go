package internal

import "time"

type Venue string

const (
	VenueBinance Venue = "BINANCE"
	VenueBybit   Venue = "BYBIT"
	VenueMock    Venue = "MOCK"
)

type TickerData struct {
	Venue   Venue
	Bid     float64
	Ask     float64
	BidSize float64
	AskSize float64
}

type Snapshot struct {
	A    TickerData
	B    TickerData
	Time time.Time
}
