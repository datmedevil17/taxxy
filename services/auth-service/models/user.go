package models


import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ─── User ────────────────────────────────────────────────────
// Single table in auth-service's private database.
// role is enforced at the application layer (not a PG enum)
// so we can extend it later without a migration.

type User struct {
	ID        string    `gorm:"primaryKey;type:uuid"         json:"id"`
	Email     string    `gorm:"uniqueIndex;not null"         json:"email"`
	Password  string    `gorm:"not null"                     json:"-"`       // bcrypt hash
	Name      string    `gorm:"not null"                     json:"name"`
	Phone     string    `gorm:"index"                        json:"phone"`
	Role      string    `gorm:"not null;check:role IN ('rider','driver')"  json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime"               json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"               json:"updated_at"`
}

// ─── Hooks ───────────────────────────────────────────────────

// BeforeCreate sets a UUID and hashes the plain-text password
// that was placed in the Password field by the handler.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	// Only hash if not already hashed (len 60 is bcrypt output length).
	if len(u.Password) < 60 {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return nil
}

// ─── Helpers ─────────────────────────────────────────────────

// ComparePassword checks plain text against the stored bcrypt hash.
func (u *User) ComparePassword(plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
}

// TableName overrides GORM's default pluralisation.
func (User) TableName() string { return "auth_users" }
