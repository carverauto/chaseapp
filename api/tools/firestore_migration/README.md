# Firestore to PostgreSQL Migration (Chases/Users/Push Tokens)

This tool migrates data from Firestore collections to PostgreSQL tables for the ChaseApp API migration.

## Overview

The migration tool:
- Reads directly from Firestore collections (`users`, `chases`, `pushTokens`)
- Transforms and maps Firestore document structure to PostgreSQL schema
- Handles UUID generation, JSONB fields, and proper data types
- Provides comprehensive validation and error reporting
- Supports dry-run mode for testing

## Prerequisites

1. **Access to Firestore project** with read permissions
2. **PostgreSQL database** with schema already created
3. **Go 1.21+** installed
4. **Authentication** configured for GCP

## Setup

1. **Set environment variables:**
```bash
export GCP_PROJECT="your-gcp-project-id"
export DATABASE_URL="postgres://user:password@host:port/database"
```

2. **Install dependencies:**
```bash
make deps
```

3. **Authenticate with Google Cloud:**
```bash
gcloud auth application-default login
```

## Migration Process

### 1. Create Backup (Recommended)
```bash
make backup
```

### 2. Run Migration
```bash
# Dry run to test (no database changes)
make dry-run

# Actual migration
make run
```

### 3. Validate Data
```bash
# Basic validation
make validate

# Comprehensive validation
make validate-full

# Check current status
make status
```

## Data Mapping

### Users Collection
- Firestore field → PostgreSQL column
- `id` → `external_id` (with new UUID primary key)
- `email` → `email`
- `displayName` → `display_name`
- `photoURL` → `photo_url`
- `provider` → `provider`
- `createdAt` → `created_at`
- `updatedAt` → `updated_at`
- `lastLoginAt` → `last_login_at`

### Chases Collection
- Firestore field → PostgreSQL column
- `id` → `id` (with new UUID primary key if needed)
- `title` → `title`
- `description` → `description`
- `chaseType` → `chase_type`
- `location` → `location` (JSONB)
- `city` → `city`
- `state` → `state`
- `country` → `country`
- `live` → `live`
- `startedAt` → `started_at`
- `endedAt` → `ended_at`
- `thumbnailUrl` → `thumbnail_url`
- `streams` → `streams` (JSONB array)
- `viewCount` → `view_count`
- `shareCount` → `share_count`
- `source` → `source`
- `sourceUrl` → `source_url`
- `createdBy` → `created_by` (UUID reference)
- `metadata` → `metadata` (JSONB)
- `createdAt` → `created_at`
- `updatedAt` → `updated_at`

### Push Tokens Collection
- Firestore field → PostgreSQL column
- `id` → `id` (with new UUID primary key if needed)
- `userId` → `user_id` (UUID reference)
- `token` → `token`
- `platform` → `platform`
- `deviceId` → `device_id`
- `deviceName` → `device_name`
- `appVersion` → `app_version`
- `subscribedTopics` → `subscribed_topics` (TEXT array)
- `isActive` → `is_active`
- `lastUsedAt` → `last_used_at`
- `metadata` → `metadata` (JSONB)
- `createdAt` → `created_at`
- `updatedAt` → `updated_at`

## Error Handling

The migration tool handles these scenarios:
- **Missing fields**: Uses NULL/defaults as appropriate
- **Invalid UUIDs**: Generates new UUIDs for primary keys
- **Missing references**: Sets foreign keys to NULL
- **Data type mismatches**: Converts and validates data types
- **Network errors**: Retries with exponential backoff

## Validation

After migration, run these checks:

### Count Validation
```sql
-- Compare counts between original and migrated data
SELECT 'users' as table_name, COUNT(*) as count FROM users WHERE deleted_at IS NULL;
SELECT 'chases' as table_name, COUNT(*) as count FROM chases WHERE deleted_at IS NULL;
SELECT 'push_tokens' as table_name, COUNT(*) as count FROM push_tokens;
```

### Data Quality Checks
```sql
-- Users without email (might need cleanup)
SELECT COUNT(*) FROM users WHERE email IS NULL AND deleted_at IS NULL;

-- Chases with invalid location data
SELECT COUNT(*) FROM chases WHERE location IS NULL AND deleted_at IS NULL;

-- Live chases older than 24 hours (potentially stuck)
SELECT COUNT(*) FROM chases
WHERE live = true AND started_at < NOW() - INTERVAL '24 hours' AND deleted_at IS NULL;

-- Orphaned push tokens
SELECT COUNT(*) FROM push_tokens pt
LEFT JOIN users u ON pt.user_id = u.id WHERE u.id IS NULL;
```

## Rollback

If you need to rollback:
```bash
# Restore from backup
make restore BACKUP_FILE=backup_20231210_120000.sql
```

## Manual Data Correction

For specific data issues after migration:

```sql
-- Fix chase IDs that need updating
UPDATE chases SET id = 'new-uuid' WHERE title = 'Specific Chase';

-- Fix user references
UPDATE chases SET created_by = 'user-uuid' WHERE title = 'Another Chase';

-- Fix token platforms
UPDATE push_tokens SET platform = 'ios' WHERE platform = 'iOS';
```

## Performance Considerations

- **Batch size**: Processes documents in batches of 500
- **Transactions**: Uses database transactions for each collection
- **Indexes**: Migration runs faster with indexes disabled (can re-enable after)
- **Memory**: Low memory usage as streams data rather than loading all at once

## Troubleshooting

### Common Issues

1. **Authentication errors:**
   - Ensure `gcloud auth application-default login` is run
   - Check GCP project permissions

2. **Database connection errors:**
   - Verify DATABASE_URL format
   - Check network connectivity to database
   - Ensure database is accessible from migration environment

3. **Memory issues:**
   - Migration is memory-efficient, but monitor large collections
   - Consider running during off-peak hours for very large datasets

### Logs

The migration provides detailed logging:
- Migration progress for each collection
- Error details for failed documents
- Summary statistics at completion

## Next Steps

After successful migration:

1. **Update application configuration** to point to PostgreSQL
2. **Run application tests** with new data
3. **Deploy to staging** for full integration testing
4. **Monitor performance** and optimize queries if needed
5. **Decommission Firestore** after cutover period
