package database

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *Queries {
	t.Helper()

	// try loading .env from project root; ignore error (CI sets vars directly)
	godotenv.Load("../../.env")

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		t.Fatal("DB_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	t.Cleanup(func() { db.Close() })
	return New(db)
}

// --- User tests ---

func TestCreateUser(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()

	user, err := q.CreateUser(ctx, CreateUserParams{
		Email:          "testcreate@example.com",
		HashedPassword: "fakehash",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	t.Cleanup(func() { q.DeleteAllUsers(ctx) })

	if user.Email != "testcreate@example.com" {
		t.Errorf("expected email testcreate@example.com, got %s", user.Email)
	}
	if user.ID.String() == "" {
		t.Error("expected non-empty UUID")
	}
	if user.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestGetUserByEmail(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()

	_, err := q.CreateUser(ctx, CreateUserParams{
		Email:          "testget@example.com",
		HashedPassword: "fakehash",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	t.Cleanup(func() { q.DeleteAllUsers(ctx) })

	user, err := q.GetUserByEmail(ctx, "testget@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if user.Email != "testget@example.com" {
		t.Errorf("expected email testget@example.com, got %s", user.Email)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()

	_, err := q.GetUserByEmail(ctx, "nonexistent@example.com")
	if err == nil {
		t.Fatal("expected error for missing user, got nil")
	}
}

func TestDeleteAllUsers(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()

	_, err := q.CreateUser(ctx, CreateUserParams{
		Email:          "todelete@example.com",
		HashedPassword: "fakehash",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if err := q.DeleteAllUsers(ctx); err != nil {
		t.Fatalf("DeleteAllUsers failed: %v", err)
	}

	_, err = q.GetUserByEmail(ctx, "todelete@example.com")
	if err == nil {
		t.Fatal("expected user to be deleted, but found it")
	}
}

// --- Chirp tests ---

func TestCreateChirp(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()

	user, err := q.CreateUser(ctx, CreateUserParams{
		Email:          "chirpauthor@example.com",
		HashedPassword: "fakehash",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	t.Cleanup(func() { q.DeleteAllUsers(ctx) })

	chirp, err := q.CreateChirp(ctx, CreateChirpParams{
		Body:   "hello chirpy world",
		UserID: user.ID,
	})
	if err != nil {
		t.Fatalf("CreateChirp failed: %v", err)
	}

	if chirp.Body != "hello chirpy world" {
		t.Errorf("expected body 'hello chirpy world', got %s", chirp.Body)
	}
	if chirp.UserID != user.ID {
		t.Errorf("expected userID %v, got %v", user.ID, chirp.UserID)
	}
	if chirp.ID.String() == "" {
		t.Error("expected non-empty chirp UUID")
	}
}

func TestGetChirp(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()

	user, err := q.CreateUser(ctx, CreateUserParams{
		Email:          "chirpget@example.com",
		HashedPassword: "fakehash",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	t.Cleanup(func() { q.DeleteAllUsers(ctx) })

	created, err := q.CreateChirp(ctx, CreateChirpParams{
		Body:   "get me by id",
		UserID: user.ID,
	})
	if err != nil {
		t.Fatalf("CreateChirp failed: %v", err)
	}

	fetched, err := q.GetChirp(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetChirp failed: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("expected chirp ID %v, got %v", created.ID, fetched.ID)
	}
	if fetched.Body != "get me by id" {
		t.Errorf("expected body 'get me by id', got %s", fetched.Body)
	}
}

func TestListChirps(t *testing.T) {
	q := setupTestDB(t)
	ctx := context.Background()

	user, err := q.CreateUser(ctx, CreateUserParams{
		Email:          "chirplist@example.com",
		HashedPassword: "fakehash",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	t.Cleanup(func() { q.DeleteAllUsers(ctx) })

	bodies := []string{"first chirp", "second chirp", "third chirp"}
	for _, body := range bodies {
		_, err := q.CreateChirp(ctx, CreateChirpParams{Body: body, UserID: user.ID})
		if err != nil {
			t.Fatalf("CreateChirp failed: %v", err)
		}
	}

	chirps, err := q.ListChirps(ctx)
	if err != nil {
		t.Fatalf("ListChirps failed: %v", err)
	}
	if len(chirps) < 3 {
		t.Errorf("expected at least 3 chirps, got %d", len(chirps))
	}
}
