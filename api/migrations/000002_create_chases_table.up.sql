-- Chases table
-- Stores live event data (police chases, etc.)
CREATE TABLE IF NOT EXISTS chases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Basic info
    title VARCHAR(500) NOT NULL,
    description TEXT,
    chase_type VARCHAR(50) NOT NULL DEFAULT 'chase',  -- chase, rocket, weather, aircraft

    -- Location (PostGIS would be better, but JSONB works for now)
    location JSONB,  -- { "lat": 34.0522, "lng": -118.2437, "address": "..." }
    city VARCHAR(255),
    state VARCHAR(100),
    country VARCHAR(100) DEFAULT 'US',

    -- Status
    live BOOLEAN NOT NULL DEFAULT false,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,

    -- Media
    thumbnail_url TEXT,
    streams JSONB DEFAULT '[]'::jsonb,  -- [{ "url": "...", "network": "NBC", "type": "m3u8" }]

    -- Engagement
    view_count INTEGER DEFAULT 0,
    share_count INTEGER DEFAULT 0,

    -- Source info
    source VARCHAR(255),  -- where the chase was reported from
    source_url TEXT,

    -- User who created it (optional - could be system)
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Metadata
    metadata JSONB DEFAULT '{}'::jsonb,  -- flexible field for additional data

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Soft delete
    deleted_at TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_chases_live ON chases(live) WHERE deleted_at IS NULL;
CREATE INDEX idx_chases_chase_type ON chases(chase_type);
CREATE INDEX idx_chases_created_at ON chases(created_at DESC);
CREATE INDEX idx_chases_started_at ON chases(started_at DESC) WHERE started_at IS NOT NULL;
CREATE INDEX idx_chases_city_state ON chases(city, state) WHERE city IS NOT NULL;
CREATE INDEX idx_chases_created_by ON chases(created_by) WHERE created_by IS NOT NULL;

-- GIN index for JSONB location queries
CREATE INDEX idx_chases_location ON chases USING GIN (location);
CREATE INDEX idx_chases_streams ON chases USING GIN (streams);

-- Updated at trigger
CREATE TRIGGER update_chases_updated_at
    BEFORE UPDATE ON chases
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
