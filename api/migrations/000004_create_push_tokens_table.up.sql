-- Push tokens table
-- Stores device tokens for push notifications
CREATE TABLE IF NOT EXISTS push_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    -- Token info
    token TEXT NOT NULL,
    platform VARCHAR(20) NOT NULL,  -- ios, android, web, safari

    -- Device info
    device_id VARCHAR(255),
    device_name VARCHAR(255),
    app_version VARCHAR(50),

    -- Subscription preferences
    subscribed_topics TEXT[] DEFAULT ARRAY[]::TEXT[],  -- ['chases', 'rockets', 'weather']

    -- Status
    is_active BOOLEAN DEFAULT true,
    last_used_at TIMESTAMPTZ,

    -- Metadata
    metadata JSONB DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Ensure unique token per platform
    UNIQUE(token, platform)
);

-- Indexes
CREATE INDEX idx_push_tokens_user_id ON push_tokens(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_push_tokens_platform ON push_tokens(platform);
CREATE INDEX idx_push_tokens_active ON push_tokens(is_active) WHERE is_active = true;
CREATE INDEX idx_push_tokens_topics ON push_tokens USING GIN (subscribed_topics);

-- Updated at trigger
CREATE TRIGGER update_push_tokens_updated_at
    BEFORE UPDATE ON push_tokens
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
