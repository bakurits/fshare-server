package db

import (
	"log"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Repository API for accessing database
type Repository struct {
	Users                   UserStore
	PasswordRestoreRequests PasswordRecoveryRequestStore
}

// NewRepository returns new repository object
func NewRepository(dialect, connectionString string) (*Repository, error) {
	db := initGorm(dialect, connectionString)
	if db == nil {
		return nil, errors.New("error while initializing gorm object")
	}

	return &Repository{
		Users:                   NewUserStore(db),
		PasswordRestoreRequests: NewPasswordRecoveryRequestStore(db),
	}, nil
}

func openDB(dialect string, connectionString string) (*gorm.DB, error) {
	switch dialect {
	case "sqlite3":
		return gorm.Open(sqlite.Open(connectionString), nil)
	case "mysql":
		return gorm.Open(mysql.Open(connectionString), nil)
	default:
		return nil, errors.New("unknown dialect")
	}
}

func initGorm(dialect string, connectionString string) *gorm.DB {
	db, err := openDB(dialect, connectionString)
	if err != nil {
		log.Println(err)
		return nil
	}

	db.Set("gorm:table_options", "charset=utf8")
	_ = db.AutoMigrate(User{})
	_ = db.AutoMigrate(PasswordRecoveryRequest{})
	return db
}
