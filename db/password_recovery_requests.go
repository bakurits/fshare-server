package db

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

const (
	deleteAfter       = time.Minute * 30
	prrCollectionName = "prr"
)

// PasswordRecoveryRequest stores data about password recovery requests
type PasswordRecoveryRequest struct {
	Token       string
	Email       string
	RequestDate time.Time
}

// PasswordRecoveryRequestStore API for password recovery request's store
type PasswordRecoveryRequestStore interface {
	Get(ctx context.Context, token string) (PasswordRecoveryRequest, error)
	Add(ctx context.Context, req PasswordRecoveryRequest) error
	Delete(ctx context.Context, token string) error
}

// NewPasswordRecoveryRequestStore returns new PasswordRecoveryRequestStore instance
func NewPasswordRecoveryRequestStore(client *firestore.Client) PasswordRecoveryRequestStore {
	return &passwordRecoveryRequestStore{client: client}
}

type passwordRecoveryRequestStore struct {
	client *firestore.Client
}

func (s *passwordRecoveryRequestStore) Get(ctx context.Context, token string) (prr PasswordRecoveryRequest, err error) {
	doc, err := s.client.Collection(prrCollectionName).Doc(token).Get(ctx)
	if err != nil {
		return PasswordRecoveryRequest{}, errors.Wrap(err, "error while retrieving recovery request")
	}
	if err = doc.DataTo(&prr); err != nil {
		return PasswordRecoveryRequest{}, errors.Wrap(err, "error while retrieving recovery request")
	}
	return
}

func (s *passwordRecoveryRequestStore) Add(ctx context.Context, req PasswordRecoveryRequest) error {
	req.RequestDate = time.Now()
	_, err := s.client.Collection(prrCollectionName).Doc(req.Token).Create(ctx, req)
	if err != nil {
		return errors.Wrap(err, "error while adding recovery request")
	}
	return nil
}

func (s *passwordRecoveryRequestStore) Delete(ctx context.Context, token string) error {
	_, err := s.client.Collection(prrCollectionName).Doc(token).Delete(ctx)
	if err != nil {
		return errors.Wrap(err, "error while deleting recovery request")
	}
	return nil
}
