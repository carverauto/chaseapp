package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/google/uuid"
)

type FirestoreUser struct {
	ID          string    `firestore:"id"`
	Email       string    `firestore:"email"`
	DisplayName string    `firestore:"displayName"`
	PhotoURL    string    `firestore:"photoURL"`
	Provider    string    `firestore:"provider"`
	CreatedAt   time.Time `firestore:"createdAt"`
	UpdatedAt   time.Time `firestore:"updatedAt"`
	LastLoginAt time.Time `firestore:"lastLoginAt"`
}

type FirestoreChase struct {
	ID          string                 `firestore:"id"`
	Title       string                 `firestore:"title"`
	Description string                 `firestore:"description"`
	ChaseType   string                 `firestore:"chaseType"`
	Location    map[string]interface{} `firestore:"location"`
	City        string                 `firestore:"city"`
	State       string                 `firestore:"state"`
	Country     string                 `firestore:"country"`
	Live        bool                   `firestore:"live"`
	StartedAt   time.Time              `firestore:"startedAt"`
	EndedAt     time.Time              `firestore:"endedAt"`
	ThumbnailURL string                `firestore:"thumbnailUrl"`
	Streams     []map[string]interface{} `firestore:"streams"`
	ViewCount   int                    `firestore:"viewCount"`
	ShareCount  int                    `firestore:"shareCount"`
	Source      string                 `firestore:"source"`
	SourceURL   string                 `firestore:"sourceUrl"`
	CreatedBy   string                 `firestore:"createdBy"`
	Metadata    map[string]interface{} `firestore:"metadata"`
	CreatedAt   time.Time              `firestore:"createdAt"`
	UpdatedAt   time.Time              `firestore:"updatedAt"`
}

type FirestorePushToken struct {
	ID               string      `firestore:"id"`
	UserID           string      `firestore:"userId"`
	Token            string      `firestore:"token"`
	Platform         string      `firestore:"platform"`
	DeviceID         string      `firestore:"deviceId"`
	DeviceName       string      `firestore:"deviceName"`
	AppVersion       string      `firestore:"appVersion"`
	SubscribedTopics []string    `firestore:"subscribedTopics"`
	IsActive         bool        `firestore:"isActive"`
	LastUsedAt       time.Time   `firestore:"lastUsedAt"`
	Metadata         interface{} `firestore:"metadata"`
	CreatedAt        time.Time   `firestore:"createdAt"`
	UpdatedAt        time.Time   `firestore:"updatedAt"`
}

type MigrationStats struct {
	UsersMigrated      int
	ChasesMigrated     int
	PushTokensMigrated int
	Errors             []string
}

func main() {
	ctx := context.Background()

	// Get environment variables
	projectID := os.Getenv("GCP_PROJECT")
	databaseURL := os.Getenv("DATABASE_URL")

	if projectID == "" || databaseURL == "" {
		log.Fatal("GCP_PROJECT and DATABASE_URL environment variables are required")
	}

	// Initialize clients
	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer firestoreClient.Close()

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	stats := &MigrationStats{}

	log.Println("Starting Firestore to PostgreSQL migration...")

	// Migrate users
	log.Println("Migrating users...")
	if err := migrateUsers(ctx, firestoreClient, db, stats); err != nil {
		log.Printf("Error migrating users: %v", err)
		stats.Errors = append(stats.Errors, fmt.Sprintf("Users migration error: %v", err))
	}

	// Migrate chases
	log.Println("Migrating chases...")
	if err := migrateChases(ctx, firestoreClient, db, stats); err != nil {
		log.Printf("Error migrating chases: %v", err)
		stats.Errors = append(stats.Errors, fmt.Sprintf("Chases migration error: %v", err))
	}

	// Migrate push tokens
	log.Println("Migrating push tokens...")
	if err := migratePushTokens(ctx, firestoreClient, db, stats); err != nil {
		log.Printf("Error migrating push tokens: %v", err)
		stats.Errors = append(stats.Errors, fmt.Sprintf("Push tokens migration error: %v", err))
	}

	// Print statistics
	printMigrationStats(stats)
}

func migrateUsers(ctx context.Context, firestoreClient *firestore.Client, db *sql.DB, stats *MigrationStats) error {
	usersCollection := firestoreClient.Collection("users")
	docs, err := usersCollection.Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to get users collection: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	insertUserQuery := `
		INSERT INTO users (id, external_id, email, display_name, photo_url, provider, created_at, updated_at, last_login_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (external_id) DO NOTHING
	`

	for _, doc := range docs {
		var user FirestoreUser
		if err := doc.DataTo(&user); err != nil {
			log.Printf("Error parsing user document %s: %v", doc.Ref.ID, err)
			continue
		}

		userID := uuid.New()
		if user.ID != "" {
			if parsedUUID, err := uuid.Parse(user.ID); err == nil {
				userID = parsedUUID
			}
		}

		_, err := tx.Exec(insertUserQuery,
			userID,
			user.ID,
			nullString(user.Email),
			nullString(user.DisplayName),
			nullString(user.PhotoURL),
			nullString(user.Provider),
			user.CreatedAt,
			user.UpdatedAt,
			nullTime(user.LastLoginAt),
		)

		if err != nil {
			log.Printf("Error inserting user %s: %v", user.ID, err)
			continue
		}

		stats.UsersMigrated++
	}

	return tx.Commit()
}

func migrateChases(ctx context.Context, firestoreClient *firestore.Client, db *sql.DB, stats *MigrationStats) error {
	chasesCollection := firestoreClient.Collection("chases")
	docs, err := chasesCollection.Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to get chases collection: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	insertChaseQuery := `
		INSERT INTO chases (id, title, description, chase_type, location, city, state, country,
			live, started_at, ended_at, thumbnail_url, streams, view_count, share_count,
			source, source_url, created_by, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		ON CONFLICT DO NOTHING
	`

	for _, doc := range docs {
		var chase FirestoreChase
		if err := doc.DataTo(&chase); err != nil {
			log.Printf("Error parsing chase document %s: %v", doc.Ref.ID, err)
			continue
		}

		chaseID := uuid.New()
		if chase.ID != "" {
			if parsedUUID, err := uuid.Parse(chase.ID); err == nil {
				chaseID = parsedUUID
			}
		}

		// Convert location to JSONB
		var locationJSON []byte
		if chase.Location != nil {
			locationJSON, _ = json.Marshal(chase.Location)
		}

		// Convert streams to JSONB
		var streamsJSON []byte
		if chase.Streams != nil {
			streamsJSON, _ = json.Marshal(chase.Streams)
		}

		// Convert metadata to JSONB
		var metadataJSON []byte
		if chase.Metadata != nil {
			metadataJSON, _ = json.Marshal(chase.Metadata)
		}

		// Handle created_by - convert to UUID if exists
		var createdByUUID *uuid.UUID
		if chase.CreatedBy != "" {
			if parsedUUID, err := uuid.Parse(chase.CreatedBy); err == nil {
				createdByUUID = &parsedUUID
			}
		}

		_, err := tx.Exec(insertChaseQuery,
			chaseID,
			chase.Title,
			nullString(chase.Description),
			nullString(chase.ChaseType),
			locationJSON,
			nullString(chase.City),
			nullString(chase.State),
			nullString(chase.Country),
			chase.Live,
			nullTime(chase.StartedAt),
			nullTime(chase.EndedAt),
			nullString(chase.ThumbnailURL),
			streamsJSON,
			chase.ViewCount,
			chase.ShareCount,
			nullString(chase.Source),
			nullString(chase.SourceURL),
			createdByUUID,
			metadataJSON,
			chase.CreatedAt,
			chase.UpdatedAt,
		)

		if err != nil {
			log.Printf("Error inserting chase %s: %v", chase.ID, err)
			continue
		}

		stats.ChasesMigrated++
	}

	return tx.Commit()
}

func migratePushTokens(ctx context.Context, firestoreClient *firestore.Client, db *sql.DB, stats *MigrationStats) error {
	tokensCollection := firestoreClient.Collection("pushTokens")
	docs, err := tokensCollection.Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to get pushTokens collection: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	insertTokenQuery := `
		INSERT INTO push_tokens (id, user_id, token, platform, device_id, device_name,
			app_version, subscribed_topics, is_active, last_used_at, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (token, platform) DO NOTHING
	`

	for _, doc := range docs {
		var token FirestorePushToken
		if err := doc.DataTo(&token); err != nil {
			log.Printf("Error parsing push token document %s: %v", doc.Ref.ID, err)
			continue
		}

		tokenID := uuid.New()
		if token.ID != "" {
			if parsedUUID, err := uuid.Parse(token.ID); err == nil {
				tokenID = parsedUUID
			}
		}

		// Handle user_id - convert to UUID if exists
		var userIDUUID *uuid.UUID
		if token.UserID != "" {
			if parsedUUID, err := uuid.Parse(token.UserID); err == nil {
				userIDUUID = &parsedUUID
			}
		}

		// Convert metadata to JSONB
		var metadataJSON []byte
		if token.Metadata != nil {
			metadataJSON, _ = json.Marshal(token.Metadata)
		}

		_, err := tx.Exec(insertTokenQuery,
			tokenID,
			userIDUUID,
			token.Token,
			token.Platform,
			nullString(token.DeviceID),
			nullString(token.DeviceName),
			nullString(token.AppVersion),
			token.SubscribedTopics,
			token.IsActive,
			nullTime(token.LastUsedAt),
			metadataJSON,
			token.CreatedAt,
			token.UpdatedAt,
		)

		if err != nil {
			log.Printf("Error inserting push token %s: %v", token.ID, err)
			continue
		}

		stats.PushTokensMigrated++
	}

	return tx.Commit()
}

func printMigrationStats(stats *MigrationStats) {
	fmt.Println("\n=== Migration Statistics ===")
	fmt.Printf("Users migrated: %d\n", stats.UsersMigrated)
	fmt.Printf("Chases migrated: %d\n", stats.ChasesMigrated)
	fmt.Printf("Push tokens migrated: %d\n", stats.PushTokensMigrated)

	if len(stats.Errors) > 0 {
		fmt.Println("\n=== Errors ===")
		for _, err := range stats.Errors {
			fmt.Printf("- %s\n", err)
		}
	} else {
		fmt.Println("\nMigration completed successfully!")
	}
}

// Helper functions for nullable values
func nullString(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func nullTime(t time.Time) interface{} {
	if t.IsZero() {
		return nil
	}
	return t
}