package models

import "time"

type Device struct {
	DeviceID     string    `db:"device_id"`
	DeviceToken  string    `db:"device_token"`
	TradePointID int64     `db:"tradepoint_id"`
	CreatedAt    time.Time `db:"created_at"`
}

type DeviceEnhanced struct {
	DeviceID        string    `db:"device_id"`
	DeviceToken     string    `db:"device_token"`
	TradePointID    int64     `db:"tradepoint_id"`
	OrganizationBin string    `db:"organization_bin"`
	CreatedAt       time.Time `db:"created_at"`
}
