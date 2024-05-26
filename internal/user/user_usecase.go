package user

import (
	"errors"
	localError "belimang/pkg/error"
	"belimang/pkg/hasher"
	tokenizer "belimang/pkg/jwt"
	// "strconv"
	// "time"
	"github.com/google/uuid"
)

type IUserUsecase interface {
	Login(req UserLoginWithRoleDTO) (*UserRegisterLoginResponse, *localError.GlobalError)
	Register(req UserRegisterWithRoleDTO) (*UserRegisterLoginResponse, *localError.GlobalError)
}

type userUsecase struct {
	repo IUserRepository
}

func NewUserUsecase(repo IUserRepository) IUserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (a *userUsecase) Login(req UserLoginWithRoleDTO) (*UserRegisterLoginResponse, *localError.GlobalError) {
	// Searcd user by username
	user, err := a.repo.FindByUsernameWithRole(req.Username, req.Role)
	if err != nil {
		return nil, localError.ErrNotFound("Account not found", err.Error)
	}

	// Check password
	passErr := hasher.CheckPassword(user.Password, req.Password)
	if passErr != nil {
		return nil, localError.ErrBadRequest(passErr.Error(), passErr)
	}

	// Generate Token
	tokenData := tokenizer.TokenData{
		ID:   user.ID,
		Name: user.Username,
		Role: string(user.Role),
	}

	token, tokenErr := tokenizer.GenerateToken(tokenData)
	if tokenErr != nil {
		return nil, localError.ErrInternalServer(tokenErr.Error(), tokenErr)
	}

	response := UserRegisterLoginResponse{
		Token: token,
	}

	return &response, nil
}

func (uc *userUsecase) Register(req UserRegisterWithRoleDTO) (*UserRegisterLoginResponse, *localError.GlobalError) {
	existingUser, _ := uc.repo.FindByUsernameWithRole(req.Username, "")
	if existingUser != nil {
		return nil, localError.ErrConflict("User already exists", errors.New("user already exists"))
	}

	existingUser, _ = uc.repo.FindByEmailWithRole(req.Email, req.Role)
	if existingUser != nil {
		return nil, localError.ErrConflict("User already exists", errors.New("user already exists"))
	}

	// Generate Password
	password, errPass := hasher.HashPassword(req.Password)
	if errPass != nil {
		return nil, localError.ErrInternalServer(errPass.Error(), errPass)
	}

	user := User{
		ID: uuid.NewString(),
		Role: UserRole(req.Role),
		Username: req.Username,
		Password: password,
		Email: req.Email,
	}

	// Generate token
	tokenData := tokenizer.TokenData{
		ID:   user.ID,
		Name: user.Username,
		Role: req.Role,
	}

	token, errToken := tokenizer.GenerateToken(tokenData)
	if errToken != nil {
		return nil, localError.ErrInternalServer(errToken.Error(), errToken)
	}

	// Create User
	err := uc.repo.Create(user)
	if err != nil {
		return nil, err
	}

	response := UserRegisterLoginResponse{
		Token: token,
	}

	return &response, nil
}