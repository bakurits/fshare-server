package db

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// User db model for users
type User struct {
	Email    string
	Password string
	Token    *oauth2.Token
}

type userFireStoreModel struct {
	Email    string
	Password string
	Token    []byte
}

func userToFirestore(user User) (userFireStoreModel, error) {
	t, err := json.Marshal(user.Token)
	if err != nil {
		return userFireStoreModel{}, err
	}
	return userFireStoreModel{
		Email:    user.Email,
		Password: user.Password,
		Token:    t,
	}, nil
}

func firestoreToUser(user userFireStoreModel) (User, error) {
	var t oauth2.Token
	if err := json.Unmarshal(user.Token, &t); err != nil {
		return User{}, err
	}
	return User{
		Email:    user.Email,
		Password: user.Password,
		Token:    &t,
	}, nil
}

// UserStore API for user's store
type UserStore interface {
	Get(context.Context, string) (User, error)
	Add(context.Context, User) error
	Update(context.Context, User) error
}

// NewUserStore returns new UserStore instance
func NewUserStore(client *firestore.Client) UserStore {
	return &userStore{client: client}
}

type userStore struct {
	client *firestore.Client
}

func (s *userStore) Get(ctx context.Context, email string) (User, error) {
	doc, err := s.client.Collection("users").Doc(email).Get(ctx)
	if err != nil {
		return User{}, errors.Wrap(err, "error while retrieving user")
	}

	var u userFireStoreModel
	if err := doc.DataTo(&u); err != nil {
		return User{}, errors.Wrap(err, "error while retrieving user")
	}
	return firestoreToUser(u)
}

func (s *userStore) Add(ctx context.Context, user User) error {
	u, err := userToFirestore(user)
	if err != nil {
		return errors.Wrap(err, "error while adding new user")
	}

	_, err = s.client.Collection("users").Doc(user.Email).Create(ctx, u)
	if err != nil {
		return errors.Wrap(err, "error while adding new user")
	}
	return nil
}

func (s *userStore) Update(ctx context.Context, user User) error {
	u, err := userToFirestore(user)
	if err != nil {
		return errors.Wrap(err, "error while updating user")
	}

	_, err = s.client.Collection("users").Doc(user.Email).Update(ctx,
		[]firestore.Update{
			{
				Path:  "password",
				Value: u.Password,
			},
			{
				Path:  "token",
				Value: u.Token,
			},
		},
	)
	if err != nil {
		return errors.Wrap(err, "error while updating user")
	}
	return nil
}
