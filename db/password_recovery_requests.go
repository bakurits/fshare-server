package db

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const deleteAfter = time.Minute * 30

// PasswordRecoveryRequest stores data about password recovery requests
type PasswordRecoveryRequest struct {
	Token       string
	Email       string
	RequestDate time.Time
}

// PasswordRecoveryRequestStore API for password recovery request's store
type PasswordRecoveryRequestStore interface {
	Get(token string) (PasswordRecoveryRequest, error)
	Add(req PasswordRecoveryRequest) error
	Delete(token string) error
	DeleteOld() error
}

// NewPasswordRecoveryRequestStore returns new PasswordRecoveryRequestStore instance
func NewPasswordRecoveryRequestStore(db *gorm.DB) PasswordRecoveryRequestStore {
	return &passwordRecoveryRequestStore{db: db}
}

type passwordRecoveryRequestStore struct {
	db *gorm.DB
}

func (s *passwordRecoveryRequestStore) Get(token string) (prr PasswordRecoveryRequest, err error) {
	if err = s.db.Where("token = ?", token).First(&prr).Error; err != nil {
		err = errors.Wrap(err, "error while getting token")
		return
	}
	return
}

func (s *passwordRecoveryRequestStore) Add(req PasswordRecoveryRequest) error {
	req.RequestDate = time.Now()
	if err := s.db.Create(req).Error; err != nil {
		return errors.Wrap(err, "error while adding new password restore request")
	}
	time.AfterFunc(deleteAfter, func() { s.db.Delete(PasswordRecoveryRequest{Email: req.Email}) })
	return nil
}

func (s *passwordRecoveryRequestStore) Delete(token string) error {
	var r PasswordRecoveryRequest
	if err := s.db.Where("token = ?", token).Delete(&r).Error; err != nil {
		return errors.Wrap(err, "error while deleting password restore request")
	}
	return nil
}

func (s *passwordRecoveryRequestStore) DeleteOld() error {
	var r PasswordRecoveryRequest
	if err := s.db.Where("request_date < ?", time.Now().Add(-deleteAfter)).Delete(&r).Error; err != nil {
		return errors.Wrap(err, "error while deleting old password restore requests")
	}
	return nil
}
