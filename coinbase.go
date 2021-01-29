package main

import "time"

var COINBASE_URL = "wss://ws-feed.pro.coinbase.com"

var subscribe = Message{
	Type:       "subscribe",
	ProductIds: []string{"BTC-USD"},
	Channels: []MessageChannel{
		MessageChannel{
			Name:       "ticker",
			ProductIds: []string{"BTC-USD"}}}}

type Message struct {
	Type          string           `json:"type"`
	ProductID     string           `json:"product_id"`
	ProductIds    []string         `json:"product_ids"`
	TradeID       int              `json:"trade_id,number"`
	OrderID       string           `json:"order_id"`
	ClientOID     string           `json:"client_oid"`
	Sequence      int64            `json:"sequence,number"`
	MakerOrderID  string           `json:"maker_order_id"`
	TakerOrderID  string           `json:"taker_order_id"`
	Time          time.Time        `json:"time,string"`
	RemainingSize string           `json:"remaining_size"`
	NewSize       string           `json:"new_size"`
	OldSize       string           `json:"old_size"`
	Size          string           `json:"size"`
	Price         string           `json:"price"`
	Side          string           `json:"side"`
	Reason        string           `json:"reason"`
	OrderType     string           `json:"order_type"`
	Funds         string           `json:"funds"`
	NewFunds      string           `json:"new_funds"`
	OldFunds      string           `json:"old_funds"`
	Message       string           `json:"message"`
	Bids          []SnapshotEntry  `json:"bids,omitempty"`
	Asks          []SnapshotEntry  `json:"asks,omitempty"`
	Changes       []SnapshotChange `json:"changes,omitempty"`
	LastSize      string           `json:"last_size"`
	BestBid       string           `json:"best_bid"`
	BestAsk       string           `json:"best_ask"`
	Channels      []MessageChannel `json:"channels"`
	UserID        string           `json:"user_id"`
	ProfileID     string           `json:"profile_id"`
	LastTradeID   int              `json:"last_trade_id"`
}

type MessageChannel struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

type SnapshotChange struct {
	Side  string
	Price string
	Size  string
}

type SnapshotEntry struct {
	Price string
	Size  string
}
