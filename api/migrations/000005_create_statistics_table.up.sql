-- Statistics table
-- Stores aggregated statistics by time period
CREATE TABLE IF NOT EXISTS statistics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Time period
    period_type VARCHAR(20) NOT NULL,  -- daily, weekly, monthly, yearly
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,

    -- Chase statistics
    chase_count INTEGER DEFAULT 0,
    total_duration_seconds BIGINT DEFAULT 0,
    avg_duration_seconds INTEGER DEFAULT 0,
    longest_chase_seconds INTEGER DEFAULT 0,
    shortest_chase_seconds INTEGER DEFAULT 0,

    -- By type breakdown
    stats_by_type JSONB DEFAULT '{}'::jsonb,  -- { "chase": 10, "rocket": 2, "weather": 5 }

    -- By location breakdown
    stats_by_city JSONB DEFAULT '{}'::jsonb,   -- { "Los Angeles": 5, "New York": 2 }
    stats_by_state JSONB DEFAULT '{}'::jsonb,  -- { "CA": 8, "NY": 3 }

    -- User engagement
    total_views INTEGER DEFAULT 0,
    unique_viewers INTEGER DEFAULT 0,
    total_shares INTEGER DEFAULT 0,

    -- Metadata
    metadata JSONB DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Ensure unique period
    UNIQUE(period_type, period_start)
);

-- Indexes
CREATE INDEX idx_statistics_period ON statistics(period_type, period_start DESC);
CREATE INDEX idx_statistics_period_range ON statistics(period_start, period_end);

-- Updated at trigger
CREATE TRIGGER update_statistics_updated_at
    BEFORE UPDATE ON statistics
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Chase durations table for detailed duration tracking
CREATE TABLE IF NOT EXISTS chase_durations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chase_id UUID NOT NULL REFERENCES chases(id) ON DELETE CASCADE,
    duration_seconds INTEGER NOT NULL,
    recorded_at DATE NOT NULL DEFAULT CURRENT_DATE,

    UNIQUE(chase_id)
);

CREATE INDEX idx_chase_durations_recorded_at ON chase_durations(recorded_at);
