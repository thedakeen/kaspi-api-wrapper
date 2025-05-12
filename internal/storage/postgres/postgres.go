package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"kaspi-api-wrapper/internal/storage"
	"time"
)

// Storage represents a PostgreSQL storage implementation
type Storage struct {
	db *sql.DB
}

// New initializes a new Storage instance by connecting to the PostgreSQL database
func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &Storage{db: db}, nil
}

// Stop closes the database connection and releases any open resources.
func (s *Storage) Stop() error {
	return s.db.Close()
}

// SaveDevice saves device in the database
func (s *Storage) SaveDevice(ctx context.Context, deviceID string, deviceToken string, tradePointID int64) error {
	const op = "storage.postgres.SaveDevice"

	// check whether device id already exists
	var existingTradePointID int64
	var existingDeviceToken string

	checkQuery := `
		SELECT tradepoint_id, device_token
		FROM devices
		WHERE device_id = $1
		LIMIT 1
	`

	err := s.db.QueryRowContext(ctx, checkQuery, deviceID).Scan(&existingTradePointID, &existingDeviceToken)

	if err == nil {
		if existingTradePointID != tradePointID {
			// device already in use in another tradepoint
			return storage.ErrDeviceExists
		}

		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s:%w", op, err)
	}

	insertQuery := `
		INSERT INTO devices (device_id, device_token, tradepoint_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (device_token) DO UPDATE 
		SET device_id = $1, tradepoint_id = $3
	`

	_, err = s.db.ExecContext(ctx, insertQuery,
		deviceID,
		deviceToken,
		tradePointID,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

// SaveDeviceEnhanced saves device in the database
func (s *Storage) SaveDeviceEnhanced(ctx context.Context, deviceID string, deviceToken string, tradePointID int64, organizationBin string) error {
	const op = "storage.postgres.SaveDeviceEnhanced"

	// check whether device id already exists
	var existingTradePointID int64
	var existingDeviceToken string

	checkQuery := `
		SELECT tradepoint_id, device_token
		FROM devices_enhanced
		WHERE device_id = $1
		LIMIT 1
	`

	err := s.db.QueryRowContext(ctx, checkQuery, deviceID).Scan(&existingTradePointID, &existingDeviceToken)

	if err == nil {
		if existingTradePointID != tradePointID {
			// device already in use in another tradepoint
			return storage.ErrDeviceExists
		}

		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s:%w", op, err)
	}

	insertQuery := `
		INSERT INTO devices_enhanced (device_id, device_token, tradepoint_id, organization_bin, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (device_token) DO UPDATE 
		SET device_id = $1, tradepoint_id = $3
	`

	_, err = s.db.ExecContext(ctx, insertQuery,
		deviceID,
		deviceToken,
		tradePointID,
		organizationBin,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}
