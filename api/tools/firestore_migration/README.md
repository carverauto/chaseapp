# Firestore to PostgreSQL Migration (Chases/Users/Push Tokens)

This is a lightweight outline for extracting Firestore data and loading it into the new PostgreSQL schema.

## Export Firestore (JSON)
```bash
# Export collections to JSON using gcloud (requires Firestore export permissions)
COLLECTIONS="chases,users,push_tokens"
gcloud firestore export gs://<bucket>/firestore-export --collection-ids=$COLLECTIONS

# Download and convert export to JSON (using firestore-export-parser or jq)
# Example with python tool (install separately):
firestore-export-parser gs://<bucket>/firestore-export > firestore.json
```

## Transform to CSV
Use a small script (example pseudo-Python) to emit CSV files matching the PostgreSQL tables:
- `chases.csv`: id,title,description,chase_type,live,started_at,ended_at,city,state,country,streams,metadata,created_at,updated_at
- `users.csv`: id,email,display_name,created_at,updated_at
- `push_tokens.csv`: id,user_id,token,platform,device_id,device_name,topics,is_active,created_at,updated_at

## Load into PostgreSQL
```bash
# Create staging tables if desired
psql "$DATABASE_URL" -f tools/firestore_migration/staging.sql

# Load CSVs
psql "$DATABASE_URL" -c "\copy chases FROM 'chases.csv' WITH (FORMAT csv, HEADER true)"
psql "$DATABASE_URL" -c "\copy users FROM 'users.csv' WITH (FORMAT csv, HEADER true)"
psql "$DATABASE_URL" -c "\copy push_tokens FROM 'push_tokens.csv' WITH (FORMAT csv, HEADER true)"
```

## Validation Checklist
- Counts match between Firestore and PostgreSQL for each collection/table
- Spot-check a handful of recent records (chases/users/tokens)
- Verify live chases have `ended_at` null and streams populated
- Ensure indexes and constraints are present after import

This doc is a starting point; add concrete scripts as migration details are finalized.
