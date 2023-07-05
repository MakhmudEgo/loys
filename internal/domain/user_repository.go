package domain

import "context"

type UserRepository interface {
	GetByID(ctx context.Context, ID string) (*User, error)
	Create(ctx context.Context, user *User) (string, error)
	// Update(ctx context.Context, user *User) error
	// Delete(ctx context.Context, user *User) error
}
