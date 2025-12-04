-- Aircraft table
-- Stores ADSB aircraft tracking data
CREATE TABLE IF NOT EXISTS aircraft (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Aircraft identification
    icao VARCHAR(10) NOT NULL,  -- ICAO 24-bit address (hex)
    callsign VARCHAR(20),       -- Flight callsign
    registration VARCHAR(20),   -- Aircraft registration (tail number)

    -- Position
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    altitude INTEGER,           -- Altitude in feet
    ground_speed INTEGER,       -- Speed in knots
    track INTEGER,              -- Heading in degrees (0-359)
    vertical_rate INTEGER,      -- Climb/descent rate in ft/min

    -- Aircraft info
    aircraft_type VARCHAR(10),  -- ICAO aircraft type code
    category VARCHAR(50),       -- Aircraft category (e.g., "media", "law_enforcement")
    operator VARCHAR(255),      -- Operator name

    -- Status
    on_ground BOOLEAN DEFAULT false,
    squawk VARCHAR(4),          -- Transponder code
    emergency VARCHAR(50),      -- Emergency status if any

    -- Tracking
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    first_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- For clustering
    cluster_id VARCHAR(50),     -- Current cluster assignment

    -- Metadata
    metadata JSONB DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE UNIQUE INDEX idx_aircraft_icao ON aircraft(icao);
CREATE INDEX idx_aircraft_callsign ON aircraft(callsign) WHERE callsign IS NOT NULL;
CREATE INDEX idx_aircraft_category ON aircraft(category) WHERE category IS NOT NULL;
CREATE INDEX idx_aircraft_last_seen ON aircraft(last_seen_at DESC);
CREATE INDEX idx_aircraft_position ON aircraft(latitude, longitude)
    WHERE latitude IS NOT NULL AND longitude IS NOT NULL;
CREATE INDEX idx_aircraft_cluster ON aircraft(cluster_id) WHERE cluster_id IS NOT NULL;

-- Updated at trigger
CREATE TRIGGER update_aircraft_updated_at
    BEFORE UPDATE ON aircraft
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Aircraft history table for tracking positions over time
CREATE TABLE IF NOT EXISTS aircraft_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aircraft_id UUID NOT NULL REFERENCES aircraft(id) ON DELETE CASCADE,

    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    altitude INTEGER,
    ground_speed INTEGER,
    track INTEGER,

    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Partition by time for efficient cleanup (keep last 24 hours typically)
CREATE INDEX idx_aircraft_history_aircraft_id ON aircraft_history(aircraft_id);
CREATE INDEX idx_aircraft_history_recorded_at ON aircraft_history(recorded_at DESC);
