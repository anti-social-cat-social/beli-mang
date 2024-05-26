package user

import "time"

type UserRole string

const (
	ADMIN   UserRole = "admin"
	USER 	UserRole = "user"
)

type User struct {
	ID                  string     	`json:"id" db:"id"`
	Role                UserRole   	`json:"role" db:"role"`
	Username            string     	`json:"username" db:"username"`
	Password            string    	`json:"password" db:"password"`
	Email				string    	`json:"email" db:"email"`
	CreatedAt           time.Time 	`json:"createdAt" db:"created_at"`
}

type UserRegisterDTO struct {
	Username	string     	`json:"username" binding:"required,min=5,max=30"`
	Password	string    	`json:"password" binding:"required,min=5,max=30"`
	Email		string    	`json:"email" binding:"required,email"`
}

type UserRegisterWithRoleDTO struct {
	Username	string     	`json:"username"`
	Password	string    	`json:"password"`
	Email		string    	`json:"email"`
	Role		string		`json:"role"`
}

type UserLoginDTO struct {
	Username	string     	`json:"username" binding:"required,min=5,max=30"`
	Password	string    	`json:"password" binding:"required,min=5,max=30"`
}

type UserLoginWithRoleDTO struct {
	Username	string     	`json:"username"`
	Password	string    	`json:"password"`
	Role		string		`json:"role"`
}

type UserRegisterLoginResponse struct {
	Token string `json:"token"`
}