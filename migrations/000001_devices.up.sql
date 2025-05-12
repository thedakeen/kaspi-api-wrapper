CREATE TABLE IF NOT EXISTS devices (
                                       device_id TEXT PRIMARY KEY,
                                       device_token TEXT UNIQUE NOT NULL,
                                       tradepoint_id BIGINT NOT NULL,
                                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),

                                       UNIQUE(device_id, tradepoint_id)
);

CREATE TABLE IF NOT EXISTS devices_enhanced (
                                                device_id TEXT PRIMARY KEY,
                                                device_token TEXT UNIQUE NOT NULL,
                                                tradepoint_id BIGINT NOT NULL,
                                                organization_bin TEXT NOT NULL,
                                                created_at TIMESTAMP NOT NULL DEFAULT NOW(),

                                                UNIQUE(device_id, tradepoint_id)
);