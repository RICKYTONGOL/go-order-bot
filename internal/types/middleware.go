// Package types implements different core types of the API.
package types

import "encoding/json"

// Middleware reponse structure
type Middleware struct {
	Success bool            `bson:"success"`
	Error   MiddlewareError `bson:"error"`
}

// MiddlewareError
type MiddlewareError struct {
	ExpiredAt int64  `bson:"expiredat"`
	UserId    string `bson:"user"`
	Email     string `bson:"email"`
	Msg       string `bson:"msg"`
}

// UnmarshalMiddleware parses the JSON-encoded middleware data.
func UnmarshalMiddleware(data []byte) (*Middleware, error) {
	var acc Middleware
	err := json.Unmarshal(data, &acc)
	return &acc, err
}

// Marshal returns the JSON encoding of middleware.
func (acc *Middleware) Marshal() ([]byte, error) {
	return json.Marshal(acc)
}
