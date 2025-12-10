-- Staging tables for Firestore migration validation
-- These tables hold the raw Firestore data before processing

-- Staging table for users
CREATE TABLE IF NOT EXISTS users_staging (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255),
    display_name VARCHAR(255),
    photo_url TEXT,
    provider VARCHAR(50),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    raw_data JSONB
);

-- Staging table for chases
CREATE TABLE IF NOT EXISTS chases_staging (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(500),
    description TEXT,
    chase_type VARCHAR(50),
    location JSONB,
    city VARCHAR(255),
    state VARCHAR(100),
    country VARCHAR(100),
    live BOOLEAN,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    thumbnail_url TEXT,
    streams JSONB,
    view_count INTEGER,
    share_count INTEGER,
    source VARCHAR(255),
    source_url TEXT,
    created_by VARCHAR(255),
    metadata JSONB,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    raw_data JSONB
);

-- Staging table for push tokens
CREATE TABLE IF NOT EXISTS push_tokens_staging (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255),
    token TEXT NOT NULL,
    platform VARCHAR(20) NOT NULL,
    device_id VARCHAR(255),
    device_name VARCHAR(255),
    app_version VARCHAR(50),
    subscribed_topics TEXT[],
    is_active BOOLEAN,
    last_used_at TIMESTAMPTZ,
    metadata JSONB,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    raw_data JSONB
);

-- Validation queries
-- Run these to compare staging data with production tables

-- Compare user counts
SELECT
    'users_staging' as table_name,
    COUNT(*) as count
FROM users_staging
UNION ALL
SELECT
    'users' as table_name,
    COUNT(*) as count
FROM users
WHERE deleted_at IS NULL;

-- Compare chase counts
SELECT
    'chases_staging' as table_name,
    COUNT(*) as count
FROM chases_staging
UNION ALL
SELECT
    'chases' as table_name,
    COUNT(*) as count
FROM chases
WHERE deleted_at IS NULL;

-- Compare push token counts
SELECT
    'push_tokens_staging' as table_name,
    COUNT(*) as count
FROM push_tokens_staging
UNION ALL
SELECT
    'push_tokens' as table_name,
    COUNT(*) as count
FROM push_tokens;

-- Data quality checks
-- Check for users without email (might need cleanup)
SELECT
    COUNT(*) as users_without_email
FROM users
WHERE email IS NULL AND deleted_at IS NULL;

-- Check for chases with invalid location data
SELECT
    COUNT(*) as chases_with_invalid_location,
    COUNT(CASE WHEN location IS NULL THEN 1 END) as null_location,
    COUNT(CASE WHEN location->>'lat' IS NULL THEN 1 END) as missing_lat,
    COUNT(CASE WHEN location->>'lng' IS NULL THEN 1 END) as missing_lng
FROM chases
WHERE deleted_at IS NULL;

-- Check for live chases that should have ended
SELECT
    COUNT(*) as potentially_stuck_live_chases
FROM chases
WHERE live = true
    AND started_at < NOW() - INTERVAL '24 hours'
    AND deleted_at IS NULL;

-- Check for orphaned push tokens (tokens without valid users)
SELECT
    COUNT(*) as orphaned_tokens
FROM push_tokens pt
LEFT JOIN users u ON pt.user_id = u.id
WHERE u.id IS NULL;

-- Sample data verification
-- Get sample of migrated users for manual verification
SELECT
    id,
    email,
    display_name,
    provider,
    created_at,
    updated_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT 5;

-- Get sample of migrated chases for manual verification
SELECT
    id,
    title,
    chase_type,
    live,
    city,
    state,
    started_at,
    created_at
FROM chases
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT 5;

-- Get sample of migrated push tokens for manual verification
SELECT
    id,
    user_id,
    platform,
    device_name,
    is_active,
    created_at
FROM push_tokens
ORDER BY created_at DESC
LIMIT 5;