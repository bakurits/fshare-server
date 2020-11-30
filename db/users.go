package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// TokenStore wraper on oauth2.Token for storing in db
type TokenStore oauth2.Token

// User db model for users
type User struct {
	Email    string `gorm:"primaryKey"`
	Password string
	Token    TokenStore
}

// Scan scan value into TokenStore, implements sql.Scanner interface
func (t *TokenStore) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("failed to unmarshal TokenStore value:", value))
	}

	result := TokenStore{}
	err := json.Unmarshal(b, &result)
	*t = result
	return err
}

// Value return json value, implement driver.Valuer interface
func (t TokenStore) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// UserStore API for user's store
type UserStore interface {
	Get(string) (User, error)
	Add(User) error
	Update(user User) error
}

// NewUserStore returns new UserStore instance
func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{db: db}
}

type userStore struct {
	db *gorm.DB
}

func (s *userStore) Get(email string) (User, error) {
	var u User
	if err := s.db.Where("email = ?", email).First(&u).Error; err != nil {
		return u, errors.Wrap(err, "error while getting token")
	}
	return u, nil
}

func (s *userStore) Add(user User) error {
	if err := s.db.Create(user).Error; err != nil {
		return errors.Wrap(err, "error while adding new user")
	}
	return nil
}

func (s *userStore) Update(user User) error {
	if err := s.db.Updates(user).Error; err != nil {
		return errors.Wrap(err, "error while updating user")
	}
	return nil
}
