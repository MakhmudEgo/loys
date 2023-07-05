package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"dz1/internal/domain"
)

var (
	ErrUserNotFound = errors.New("user not exist")
)

type UserRepositoryImpl struct {
	db *pgxpool.Pool
}

var _ domain.UserRepository = (*UserRepositoryImpl)(nil)

func NewUser(db *pgxpool.Pool) domain.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, ID string) (*domain.User, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT password, first_name, second_name, birthdate, gender, biography, city
				FROM users WHERE id = $1`

	rows, err := conn.Query(ctx, query, ID)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, ErrUserNotFound
	}

	id, _ := uuid.Parse(ID)
	user := domain.User{
		ID: id,
	}

	if err = rows.Scan(
		&user.Password,
		&user.FirstName,
		&user.SecondName,
		&user.Birthdate,
		&user.Gender,
		&user.Biography,
		&user.City,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) (string, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	query := `INSERT INTO users(password, first_name, second_name, birthdate, gender, biography, city)
					VALUES($1, $2,$3,$4,$5,$6,$7) RETURNING id`
	rows, err := conn.Query(ctx, query,
		user.Password,
		user.FirstName,
		user.SecondName,
		user.Birthdate,
		user.Gender,
		user.Biography,
		user.City,
	)

	var ID string
	if rows.Next() {
		if err = rows.Scan(&ID); err != nil {
			return "", err
		}
	}

	return ID, nil
}

func (r *UserRepositoryImpl) Update(user *domain.User) error {
	panic("impl")
}

func (r *UserRepositoryImpl) Delete(user *domain.User) error {
	panic("impl")
}
