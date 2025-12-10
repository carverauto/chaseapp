.PHONY: help deps run validate clean backup-restore

# Default target
help:
	@echo "Firestore to PostgreSQL Migration"
	@echo ""
	@echo "Available targets:"
	@echo "  deps         - Install dependencies"
	@echo "  run          - Run the migration"
	@echo "  validate     - Run data validation queries"
	@echo "  backup       - Create database backup before migration"
	@echo "  restore      - Restore from backup (if needed)"
	@echo "  clean        - Clean up build artifacts"

# Install dependencies
deps:
	go mod tidy
	go install cloud.google.com/go/firestore@latest
	go install github.com/jackc/pgx/v5/stdlib@latest

# Set required environment variables
.env-check:
	@if [ -z "$(GCP_PROJECT)" ]; then \
		echo "Error: GCP_PROJECT environment variable is not set"; \
		exit 1; \
	fi
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Error: DATABASE_URL environment variable is not set"; \
		exit 1; \
	fi

# Run the migration
run: .env-check
	@echo "Starting Firestore to PostgreSQL migration..."
	@echo "GCP Project: $(GCP_PROJECT)"
	@echo "Database: $(DATABASE_URL)"
	@echo ""
	@read -p "Continue? (y/N) " confirm && [ "$$confirm" = "y" ] || exit 1
	go run main.go

# Run with dry-run mode (no database writes)
dry-run: .env-check
	@echo "Running migration in dry-run mode..."
	DRY_RUN=true go run main.go

# Create database backup
backup: .env-check
	@echo "Creating database backup..."
	@TIMESTAMP=$$(date +%Y%m%d_%H%M%S); \
	pg_dump "$(DATABASE_URL)" > "backup_$${TIMESTAMP}.sql"; \
	echo "Backup created: backup_$${TIMESTAMP}.sql"

# Validate migrated data
validate: .env-check
	@echo "Running data validation..."
	@echo "Setting up staging tables..."
	psql "$(DATABASE_URL)" -f staging.sql
	@echo "Running validation queries..."
	psql "$(DATABASE_URL)" -c "
		SELECT '=== User Count Comparison ===';
		SELECT 'users_staging' as table_name, COUNT(*) as count FROM users_staging
		UNION ALL
		SELECT 'users' as table_name, COUNT(*) as count FROM users WHERE deleted_at IS NULL;

		SELECT '=== Chase Count Comparison ===';
		SELECT 'chases_staging' as table_name, COUNT(*) as count FROM chases_staging
		UNION ALL
		SELECT 'chases' as table_name, COUNT(*) as count FROM chases WHERE deleted_at IS NULL;

		SELECT '=== Push Token Count Comparison ===';
		SELECT 'push_tokens_staging' as table_name, COUNT(*) as count FROM push_tokens_staging
		UNION ALL
		SELECT 'push_tokens' as table_name, COUNT(*) as count FROM push_tokens;
	"

# Run comprehensive validation
validate-full: .env-check
	@echo "Running comprehensive data validation..."
	psql "$(DATABASE_URL)" -c "
		-- Check for data quality issues
		SELECT '=== Data Quality Issues ===';

		-- Users without email
		SELECT 'Users without email: ' || COUNT(*) FROM users WHERE email IS NULL AND deleted_at IS NULL;

		-- Chases with invalid location
		SELECT 'Chases with NULL location: ' || COUNT(*) FROM chases WHERE location IS NULL AND deleted_at IS NULL;

		-- Live chases older than 24 hours
		SELECT 'Live chases older than 24h: ' || COUNT(*) FROM chases
		WHERE live = true AND started_at < NOW() - INTERVAL '24 hours' AND deleted_at IS NULL;

		-- Orphaned push tokens
		SELECT 'Orphaned push tokens: ' || COUNT(*) FROM push_tokens pt
		LEFT JOIN users u ON pt.user_id = u.id WHERE u.id IS NULL;
	"

# Restore from backup
restore:
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "Error: BACKUP_FILE environment variable is not set"; \
		echo "Usage: make restore BACKUP_FILE=backup_20231210_120000.sql"; \
		exit 1; \
	fi
	@echo "Restoring database from $(BACKUP_FILE)..."
	@read -p "This will overwrite current data. Continue? (y/N) " confirm && [ "$$confirm" = "y" ] || exit 1
	psql "$(DATABASE_URL)" < "$(BACKUP_FILE)"
	@echo "Database restored from $(BACKUP_FILE)"

# Clean up
clean:
	@echo "Cleaning up build artifacts..."
	rm -f main
	rm -f backup_*.sql
	go clean

# Show migration status
status: .env-check
	@echo "Current database status:"
	psql "$(DATABASE_URL)" -c "
		SELECT 'users' as table_name, COUNT(*) as count FROM users WHERE deleted_at IS NULL
		UNION ALL
		SELECT 'chases' as table_name, COUNT(*) as count FROM chases WHERE deleted_at IS NULL
		UNION ALL
		SELECT 'push_tokens' as table_name, COUNT(*) as count FROM push_tokens;
	"