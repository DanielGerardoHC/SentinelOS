package auth

import "errors"

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
			return nil, errors.New("usuario deshabilitado")
		}

		if err := CheckPassword(u.PasswordHash, password); err != nil {
			return nil, errors.New("credenciales inválidas")
		}

		return &u, nil
	}

	return nil, errors.New("usuario no existe")
}
