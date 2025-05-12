package storage

import "errors"

var (
	ErrDeviceExists = errors.New("device already in use in another tradepoint")
	//ErrTradePointNotFound = errors.New("tradepoint not found")
)
