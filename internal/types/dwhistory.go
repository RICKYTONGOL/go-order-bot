// Package types implements different core types of the API.
package types

import "encoding/json"

// DWHistory is the dwhistory base row
type DWHistory struct {
	From        string  `bson:"from"`
	To          string  `bson:"to"`
	UserID      uint64  `bson:"userid"`
	TxHash      string  `bson:"txhash"`
	BlockNumber uint64  `bson:"blocknum"`
	Timestamp   uint64  `bson:"timestamp"`
	Value       string  `bson:"value"`
	FloatValue  float64 `bson:"fvalue"`
	Symbol      string  `bson:"symbol"`
	Type        uint64  `bson:"type"`
}

// DWHistoryList is the list of dwhistory
type DWHistoryList struct {
	Success  bool           `bson:"success"`
	TotalCnt uint64         `bson:"totalcnt"`
	Data     []DWHistory    `bson:"data"`
	Error    DWHistoryError `bson:"error"`
}

// Error is the signup field error
type DWHistoryError struct {
	UnknownError uint64 `bson:"unknown"`
	Msg          string `bson:"msg"`
}

// UnmarshalDWHistory parses the JSON-encoded dwhistory data.
func UnmarshalDWHistory(data []byte) (*DWHistory, error) {
	var history DWHistory
	err := json.Unmarshal(data, &history)
	return &history, err
}

// Marshal returns the JSON encoding of dwhistory.
func (acc *DWHistory) Marshal() ([]byte, error) {
	return json.Marshal(acc)
}
