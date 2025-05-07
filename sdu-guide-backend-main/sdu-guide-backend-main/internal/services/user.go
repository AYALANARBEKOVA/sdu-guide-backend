package services

import (
	"errors"
	"net/http"
	"sdu-guide/internal/logger"
	"sdu-guide/internal/repositories"
	"sdu-guide/internal/structures"
	"sdu-guide/internal/utils"
	"time"
)

type User struct {
	repo *repositories.Repository
}

func newUserService(repo *repositories.Repository) *User {
	return &User{repo: repo}
}

func (m *User) Get(id int64) (structures.User, error) {

	return m.repo.UserRepo.Get(id)
}

func (m *User) GetBy(field string, value interface{}) (structures.User, error) {

	return m.repo.UserRepo.GetBy(field, value)
}

func (m *User) Create(userRegister structures.UserRegister) error {

	if !utils.IsValidEmail(userRegister.Email) {
		return errors.New("wrong email format")
	}

	if len(userRegister.Login) < 5 || len(userRegister.Password) < 6 {
		return errors.New("wrong input value")
	}

	hashPassword, err := utils.HashPassword(userRegister.Password)
	if err != nil {
		return err
	}

	user := structures.User{
		PasswordHash:     hashPassword,
		Username:         userRegister.Login,
		Email:            userRegister.Email,
		RegistrationDate: time.Now(),
		LastLogin:        time.Now(),
	}
	return m.repo.UserRepo.Create(user)
}

func (m *User) Upadte(user structures.User) error {

	prevUser, err := m.repo.UserRepo.GetBy("_id", user.ID)
	if err != nil {
		return err
	}

	user.PasswordHash = prevUser.PasswordHash

	return m.repo.UserRepo.Upadte(user)
}

func (s *User) LoginToSystem(login structures.Login) (*http.Cookie, structures.Session, error) {

	user, err := s.repo.UserRepo.GetBy("email", login.Email)
	if err != nil {
		return nil, structures.Session{}, err
	}
	if !utils.CheckPasswordHash(login.Password, user.PasswordHash) {
		return nil, structures.Session{}, errors.New("wrong password")
	}
	var session structures.Session
	cookie := &http.Cookie{
		Path:     "/",
		HttpOnly: true,
		Name:     "Token",
	}
	newToken, err := utils.CreateToken(login.Email)
	if err != nil {
		return nil, structures.Session{}, err
	}

	user.LastLogin = time.Now()

	newTime := time.Now().Add(time.Minute * 30)
	session.Token = newToken
	session.ExpiredDate = newTime
	session.UserID = int(user.ID)
	cookie.Value = newToken
	cookie.Expires = newTime

	if err := s.repo.UserRepo.Upadte(user); err != nil {
		logger.Error.Println("Can't update time", err)
		return cookie, session, nil
	}

	return cookie, session, nil
}

func (s *User) Logout(cookie *http.Cookie) error {

	cookie.Path = "/"
	cookie.MaxAge = -1
	cookie.HttpOnly = false
	cookie.Value = ""
	return nil

}

func (m *User) GetUserFromSession(session structures.Session) (structures.User, error) {

	user, err := m.repo.UserRepo.GetBy("_id", session.UserID)
	if err != nil {
		return structures.User{}, err
	}
	return user, nil
}
