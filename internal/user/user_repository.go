package user

import (
	"database/sql"
	"errors"
	// "fmt"
	localError "belimang/pkg/error"
	"log"
	// "strings"

	// "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	FindById(id string) (*User, *localError.GlobalError)
	FindByUsernameWithRole(username string, role string) (*User, *localError.GlobalError)
	FindByEmailWithRole(email string, role string) (*User, *localError.GlobalError)
	Create(entity User) *localError.GlobalError
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

// This can be use for authentication process
func (u *userRepository) FindById(id string) (*User, *localError.GlobalError) {
	user := User{}

	if err := u.db.Get(&user, "SELECT * FROM users where id=$1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, localError.ErrNotFound("User data not found", err)
		}

		log.Println(err)

		return nil, &localError.GlobalError{
			Code:    400,
			Message: "Not found",
			Error:   err,
		}

	}

	return &user, nil
}

// This can be use for authentication process
func (u *userRepository) FindByUsernameWithRole(username string, role string) (*User, *localError.GlobalError) {
	user := User{}
	var err error

	if role != "" {
		err = u.db.Get(&user, "SELECT * FROM users where username=$1 AND role=$2;", username, role);
	} else {
		err = u.db.Get(&user, "SELECT * FROM users where username=$1;", username);
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, localError.ErrNotFound("User data not found", err)
		}

		log.Println(err)

		return nil, &localError.GlobalError{
			Code:    400,
			Message: "Not found",
			Error:   err,
		}

	}

	return &user, nil
}

// This can be use for authentication process
func (u *userRepository) FindByEmailWithRole(email string, role string) (*User, *localError.GlobalError) {
	user := User{}

	if err := u.db.Get(&user, "SELECT * FROM users where email=$1 AND role=$2", email, role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, localError.ErrNotFound("User data not found", err)
		}

		log.Println(err)

		return nil, &localError.GlobalError{
			Code:    400,
			Message: "Not found",
			Error:   err,
		}

	}

	return &user, nil
}

// Store new user to database
func (u *userRepository) Create(entity User) *localError.GlobalError {
	q := "INSERT INTO users (id, role, username, password, email) values (:id, :role, :username, :password, :email);"

	// Insert into database
	_, err := u.db.NamedExec(q, &entity)
	if err != nil {
		return localError.ErrInternalServer(err.Error(), err)
	}

	return nil
}