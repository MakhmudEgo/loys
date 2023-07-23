package domain

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUser(ctx context.Context, userID string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*UserCreateResp, error)
	UpdateUser(userID string, updatedUser User) error
	DeleteUser(userID string) error
	AuthenticateUser(ctx context.Context, userID, password string) (*User, error)
	AuthorizeUser(userID string, permission string) bool
	SearchUsersByNames(ctx context.Context, firstName string, lastName string) ([]User, error)
	// ...
}

type userService struct {
	userRepository UserRepository
	// ...
}

func (u userService) GetUser(ctx context.Context, userID string) (*User, error) {
	return u.userRepository.GetByID(ctx, userID)
}

func NewUserService(userRepository UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (u userService) CreateUser(ctx context.Context, user *User) (*UserCreateResp, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(password)

	userID, err := u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &UserCreateResp{
		ID: userID,
	}, nil
}

func (u userService) SearchUsersByNames(ctx context.Context, firstName string, lastName string) ([]User, error) {
	return u.userRepository.FindByNames(ctx, firstName, lastName)
}

func (u userService) UpdateUser(userID string, updatedUser User) error {
	// TODO implement me
	panic("implement me")
}

func (u userService) DeleteUser(userID string) error {
	// TODO implement me
	panic("implement me")
}

func (u userService) AuthenticateUser(ctx context.Context, userID, password string) (*User, error) {
	user, err := u.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (u userService) AuthorizeUser(userID string, permission string) bool {
	// TODO implement me
	panic("implement me")
}
