package store

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Username       string `binding:"required,min=5,max=30"`
	Password       string `pg:"-" binding:"required,min=7,max=32"`
	HashedPassword []byte `json:"-"`
	Salt           []byte `json:"-"`
	CreatedAt      time.Time
	ModifiedAt     time.Time
	Posts          []*Post `json:"-" pg:"fk:user_id,rel:has-many,on_delete:CASCADE"`
}

var _ pg.AfterSelectHook = (*User)(nil)

func (user *User) AfterSelect(ctx context.Context) error {
	if user.Posts == nil {
		user.Posts = []*Post{}
	}
	return nil
}

func AddUser(user *User) error {
	salt, err := GenerateSalt()
	if err != nil {
		return err
	}
	toHash := append([]byte(user.Password), salt...)
	hashedPassword, err := bcrypt.GenerateFromPassword(toHash, bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing password")
		return err
	}

	user.Salt = salt
	user.HashedPassword = hashedPassword

	_, err = db.Model(user).Returning("*").Insert()
	if err != nil {
		log.Error().Err(err).Msg("Error inserting new user")
		return dbError(err)
	}
	return nil
}

func Authenticate(username, password string) (*User, error) {
	user := new(User)
	if err := db.Model(user).Where(
		"username = ?", username).Select(); err != nil {
		log.Error().Err(err).Str("username", username).Msg("Error fetching user for authentication")
		return nil, dbError(err)
	}
	salted := append([]byte(password), user.Salt...)
	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, salted); err != nil {
		log.Error().Err(err).Msg("Error comparing hash and password")
		return nil, err
	}
	return user, nil
}

func FetchUser(id int) (*User, error) {
	user := new(User)
	user.ID = id
	err := db.Model(user).Returning("*").WherePK().Select()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching user")
		return nil, dbError(err)
	}
	return user, nil
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		log.Error().Err(err).Msg("Unable to create salt")
		return nil, err
	}
	return salt, nil
}
