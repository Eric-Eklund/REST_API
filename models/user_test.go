package models

import (
	"REST_API/db"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) func() {
	originalDB := db.DB

	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	createUsersTable := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)`

	_, err = testDB.Exec(createUsersTable)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	db.DB = testDB

	return func() {
		err := testDB.Close()
		if err != nil {
			return
		}
		db.DB = originalDB
	}
}

func TestUser_Save(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	type fields struct {
		ID       int64
		Email    string
		Password string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid user",
			fields: fields{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Empty password",
			fields: fields{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:       tt.fields.ID,
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}
			err := u.Save()
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && u.ID == 0 {
				t.Errorf("Save() should set user ID, but ID is still 0")
			}
		})
	}

	t.Run("Duplicate email should fail", func(t *testing.T) {
		u1 := &User{
			Email:    "unique@example.com",
			Password: "password123",
		}
		err := u1.Save()
		if err != nil {
			t.Fatalf("First user save should succeed: %v", err)
		}

		u2 := &User{
			Email:    "unique@example.com",
			Password: "different123",
		}
		err = u2.Save()
		if err == nil {
			t.Error("Second user with duplicate email should fail")
		}
	})
}

func TestUser_ValidateCredentials(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	testUser := &User{
		Email:    "test@example.com",
		Password: "correctpassword",
	}
	err := testUser.Save()
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid credentials",
			email:    "test@example.com",
			password: "correctpassword",
			wantErr:  false,
		},
		{
			name:     "Invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "Non-existent email",
			email:    "nonexistent@example.com",
			password: "anypassword",
			wantErr:  true,
		},
		{
			name:     "Empty email",
			email:    "",
			password: "password",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Email:    tt.email,
				Password: tt.password,
			}
			err := u.ValidateCredentials()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && u.ID == 0 {
				t.Errorf("ValidateCredentials() should set user ID, but ID is still 0")
			}
		})
	}
}
