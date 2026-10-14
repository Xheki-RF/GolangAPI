package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"social-network/cache"
	"social-network/domain"
	"social-network/entity"

	"github.com/txix-open/isp-kit/db"
	"github.com/txix-open/isp-kit/errors"
	"github.com/txix-open/isp-kit/log"
	"github.com/txix-open/isp-kit/metrics/sql_metrics"
)

type User struct {
	l     log.Logger
	db    db.DB
	users *cache.Cache[entity.User]
	tx    db.Transactional
}

func NewUser(l log.Logger, db db.DB, transactionRunner db.Transactional) *User {
	return &User{
		l:     l,
		db:    db,
		users: cache.NewCache[entity.User](),
		tx:    transactionRunner,
	}
}

func (repo User) RegisterUserRepo(ctx context.Context, email string) error {
	ctx = sql_metrics.OperationLabelToContext(ctx, "User.RegisterUserRepo")

	if email == domain.ServiceUserEmail {
		return entity.ErrUserEmailInternalServerError
	}

	err := repo.tx.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		if _, exists := repo.users.Get(email); exists {
			return entity.ErrUserEmailConflict
		}

		err := repo.selectEmailQuery(ctx, tx, email, false)
		if err != nil {
			return errors.WithMessage(err, "selectEmailQuery")
		}

		newUser, err := repo.insertUserEmail(ctx, tx, email)
		if err != nil {
			return errors.WithMessage(err, "updateUsersTable")
		}

		repo.users.Add(email, *newUser)

		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "transaction RegisterUserRepo")
	}

	repo.l.Debug(ctx, fmt.Sprintf("Users cache: %#v", repo.users.GetData()))

	return nil
}
func (repo User) UpdateUserRepo(ctx context.Context, oldEmail string, newEmail string) error {
	ctx = sql_metrics.OperationLabelToContext(ctx, "User.UpdateUserRepo")
	if newEmail == domain.ServiceUserEmail || oldEmail == domain.ServiceUserEmail {
		return entity.ErrUserEmailInternalServerError
	}

	err := repo.tx.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		if _, exists := repo.users.Get(newEmail); exists {
			return entity.ErrUserEmailConflict
		}

		err := repo.selectEmailQuery(ctx, tx, newEmail, false)
		if err != nil {
			return entity.ErrUserEmailConflict
		}

		newUser, err := repo.updateUserEmail(ctx, tx, oldEmail, newEmail)
		if err != nil {
			return errors.WithMessage(err, "updateUserEmail")
		}

		repo.users.Delete(oldEmail)
		repo.users.Add(newEmail, *newUser)

		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "transaction UpdateUserRepo")
	}

	repo.l.Debug(ctx, fmt.Sprintf("Users cache: %#v", repo.users.GetData()))

	return nil
}
func (repo User) DeleteUserRepo(ctx context.Context, email string) error {
	ctx = sql_metrics.OperationLabelToContext(ctx, "User.DeleteUserRepo")

	if email == domain.ServiceUserEmail {
		return entity.ErrUserEmailInternalServerError
	}

	err := repo.tx.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		err := repo.selectEmailQuery(ctx, tx, email, true)
		if err != nil {
			return errors.WithMessage(err, "selectEmailQuery")
		}

		err = repo.deleteUserByEmail(ctx, email)
		if err != nil {
			return errors.WithMessage(err, "deleteUserByEmail")
		}

		repo.users.Delete(email)

		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "transaction DeleteUserRepo")
	}

	repo.l.Debug(ctx, fmt.Sprintf("Users cache: %#v", repo.users.GetData()))

	return nil
}

func (repo User) selectEmailQuery(ctx context.Context, tx *db.Tx, email string, returnNotFound bool) error {
	newUser := entity.User{}
	query := `
		SELECT
			email
		FROM
			users
		WHERE
			email = $1`

	err := tx.QueryRowContext(ctx, query, email).Scan(&newUser.Email)
	switch returnNotFound{
	case true:
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrUserEmailNotFound
		}
	case false:
		if len(newUser.Email) > 0 {
			return entity.ErrUserEmailConflict
		}
	}

	return nil
}

func (repo User) insertUserEmail(ctx context.Context, tx *db.Tx, email string) (*entity.User, error) {
	newUser := entity.User{}
	var statistic json.RawMessage
	query := `
	INSERT INTO
		users(email)
	VALUES
		($1)
	RETURNING id, email, created_at, last_read_id, statistic`

	err := tx.QueryRowContext(ctx, query, email).Scan(
		&newUser.Id,
		&newUser.Email,
		&newUser.CreatedAt,
		&newUser.LastReadId,
		&statistic,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "INSERT user")
	}

	err = json.Unmarshal(statistic, &newUser.Statistic)
	if err != nil {
		return nil, errors.WithMessage(err, "unmarshal json")
	}

	return &newUser, nil
}

func (repo User) deleteUserByEmail(ctx context.Context, email string) error {
	query := `
	DELETE FROM
		users
	WHERE
		email = $1`

	_, err := repo.db.Exec(ctx, query, email)
	if err != nil {
		return errors.WithMessage(err, "DELETE user")
	}

	return nil
}

func (repo User) updateUserEmail(ctx context.Context, tx *db.Tx, oldEmail string, newEmail string) (*entity.User, error) {
	newUser := entity.User{}
	var statistic json.RawMessage
	query := `
		UPDATE
			users
		SET
			email=$1
		WHERE
			email=$2
		RETURNING id, email, created_at, last_read_id, statistic`

	err := tx.QueryRowContext(ctx, query, newEmail, oldEmail).Scan(
		&newUser.Id,
		&newUser.Email,
		&newUser.CreatedAt,
		&newUser.LastReadId,
		&statistic,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrUserEmailNotFound
	}
	if err != nil {
		return nil, errors.WithMessage(err, "UPDATE user email")
	}

	err = json.Unmarshal(statistic, &newUser.Statistic)
	if err != nil {
		return nil, errors.WithMessage(err, "unmarshal json")
	}

	return &newUser, nil
}