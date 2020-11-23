package service

type JSONParams struct {
	Action string   `json:"action"`
	Subs   []string `json:"subs"`
}

type Message struct {
	// Type of the message, this is 8 for an order book update and 9 for an order book snapshot.
	Type string `json:"TYPE"`

	// The market you have requested (name of the market e.g. Coinbase, Kraken, etc)
	M string `json:"M"`

	// The mapped from asset (base symbol / coin) you have requested (e.g. BTC, ETH, etc.)
	FSYM string `json:"FSYM"`

	// The mapped to asset (quote/counter symbol/coin) you have requested (e.g. BTC, USD, etc.)
	TSYM string `json:"TSYM"`

	// The side is 0 for BID and 1 for ASK
	SIDE int `json:"SIDE"`

	// The action you need to apply on the snapshot,
	// 1 for ADD (add this position to your ordebook),
	// 2 for REMOVE (take this position out of your orderbook, REMOVE orders also have a quantity of 0),
	// 3 for NOACTION (you should not see these messages they represent updates that we receive from the exchange but
	// have no impact on the orderbook) and 4 for CHANGE/UPDATE (update the available quantity for this position)
	ACTION int `json:"ACTION"`

	// Our internal order book sequence, the snapshot you get when you subscribe will have the starting sequence and
	// all other updates will be increments of 1 on the sequence. This does not reset.
	CCSEQ int `json:"CCSEQ"`

	// The price in the to asset (quote/counter symbol/coin) of the order book position
	// (e.g. for a BTC-USD order book update you would get price in USD how much do i need to pay per BTC to hit this possition)
	P float32 `json:"P"`

	// The from asset (base symbol/coin) volume of position (e.g. for a BTC-USD order book update you would get volume in BTC -
	// I can buy or sell x BTC at the y price).
	Q float32 `json:"Q"`

	// The external exchange sequence if they have one.
	SEQ int64 `json:"SEQ"`

	// The external exchange reported timestamp in nanoseconds.
	// If they do not provide one, we store the time we receive the message in this field.
	REPORTEDNS int64 `json:"REPORTEDNS"`

	// The difference in nanoseconds between the REPORTEDNS and the time we publish the update on our internal network
	// (network delay from exchange to us + internal aggregation delay + internal delay in propagating the message)
	DELAYNS int `json:"DELAYNS"`

	BID []Bid `json:"BID"`
}

type Bid struct {
	P float32 `json:"P"`
	Q float32 `json:"Q"`
}
