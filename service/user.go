package service

import (
	"context"
	"social-network/domain"

	"github.com/txix-open/isp-kit/errors"
)

type UserRepo interface {
	RegisterUserRepo(ctx context.Context, email string) error
	UpdateUserRepo(ctx context.Context, oldEmail string, newEmail string) error
	DeleteUserRepo(ctx context.Context, email string) error
}

type User struct {
	repo UserRepo
}

func NewUser(repo UserRepo) *User {
	return &User{
		repo: repo,
	}
}

func (srv User) RegisterUserSrv(ctx context.Context, request domain.Email) error {
	err := srv.repo.RegisterUserRepo(ctx, request.Email)
	if err != nil {
		return errors.WithMessage(err, "RegisterUserSrv")
	}

	return nil
}

func (srv User) UpdateUserSrv(ctx context.Context, request domain.UpdateEmail) error {
	err := srv.repo.UpdateUserRepo(ctx, request.OldEmail, request.NewEmail)
	if err != nil {
		return errors.WithMessage(err, "UpdateUserSrv")
	}

	return nil
}

func (srv User) DeleteUserSrv(ctx context.Context, request domain.DeleteEmail) error {
	err := srv.repo.DeleteUserRepo(ctx, request.Email)
	if err != nil {
		return errors.WithMessage(err, "DeleteUserSrv")
	}

	return nil
}
