package models

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO snippetbox.users (name, email, hashed_password, created)
VALUES ($1, $2, $3, $4)
`
	createTime := time.Now()
	_, err = m.DB.Exec(context.Background(), stmt, name, email, string(hashedPassword), createTime)
	if err != nil {
		if strings.Contains(err.Error(), "users_uc_email") {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (m UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m UserModel) Exists(id int) (bool, error) {
	return false, nil
}
