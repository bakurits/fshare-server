package db

import (
	"cloud.google.com/go/firestore"
)

// Repository API for accessing database
type Repository struct {
	Users                   UserStore
	PasswordRestoreRequests PasswordRecoveryRequestStore
}

// NewRepository returns new repository object
func NewRepository(client *firestore.Client) (*Repository, error) {
	return &Repository{
		Users:                   NewUserStore(client),
		PasswordRestoreRequests: NewPasswordRecoveryRequestStore(client),
	}, nil
}
