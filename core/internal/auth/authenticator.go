package auth

import (
	"sentinelos/core/pkg/utils"
)

type AuthService struct {
	users *UsersFile
}

func NewAuthService(users *UsersFile) *AuthService {
	return &AuthService{users: users}
}

func (a *AuthService) Authenticate(username, password string) (*User, error) {
	for _, u := range a.users.Users {
		if u.Username != username {
			continue
		}

		if !u.Enabled {
			return nil, &utils.APIError{Code: "ERR_SEC_5004", Message: "User account disabled"}
		}

		if err := CheckPassword(u.PasswordHash, password); err != nil {
			return nil, &utils.APIError{Code: "ERR_SEC_5005", Message: "Invalid credentials"}
		}

		return &u, nil
	}

	return nil, &utils.APIError{Code: "ERR_SEC_5005", Message: "Invalid credentials"}
}
