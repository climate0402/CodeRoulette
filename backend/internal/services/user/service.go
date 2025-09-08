package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(req UserCreate) (*User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := &User{
		Username:   req.Username,
		Email: 		req.Email,
		Password:   string(hash),
		Rating: 	1200,
		Wins:		0,
		Losses: 	0,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(req UserLogin) (*User, error) {
	// user, err := s.repo.GetByEmail(req.Email)
	// if err != nil {
	// 	return nil, errors.New("user email not found")
	// }
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("username not found")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}