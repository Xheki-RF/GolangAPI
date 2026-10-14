package controller

import (
	"context"
	"net/http"

	"social-network/domain"
	"social-network/entity"

	"github.com/txix-open/isp-kit/errors"
	"github.com/txix-open/isp-kit/http/apierrors"
)

type UserService interface {
	RegisterUserSrv(ctx context.Context, request domain.Email) error
	UpdateUserSrv(ctx context.Context, request domain.UpdateEmail) error
	DeleteUserSrv(ctx context.Context, request domain.DeleteEmail) error
}

type User struct {
	srv UserService
}

func NewUser(srv UserService) *User {
	return &User{
		srv: srv,
	}
}

func (ctrl User) RegisterUserCtrl(ctx context.Context, request domain.Email) error {
	err := ctrl.srv.RegisterUserSrv(ctx, request)
	switch {
	case errors.Is(err, entity.ErrUserEmailConflict):
		return apierrors.New(http.StatusConflict, domain.ErrCodeUserEmailConflict, "Email already exists", err)
	case errors.Is(err, entity.ErrUserEmailInternalServerError):
		return apierrors.New(http.StatusInternalServerError, domain.ErrCodeUserInternalServerError, "Something went wrong", err)
	case err != nil:
		return errors.WithMessage(err, "RegisterUserCtrl")
	default:
		return nil
	}
}

func (ctrl User) UpdateUserCtrl(ctx context.Context, request domain.UpdateEmail) error {
	err := ctrl.srv.UpdateUserSrv(ctx, request)
	switch {
	case errors.Is(err, entity.ErrUserEmailConflict):
		return apierrors.New(http.StatusConflict, domain.ErrCodeUserEmailConflict, "Email already exists", err)
	case errors.Is(err, entity.ErrUserEmailNotFound):
		return apierrors.New(http.StatusNotFound, domain.ErrCodeUserEmailNotFound, "Email not found", err)
	case errors.Is(err, entity.ErrUserEmailInternalServerError):
		return apierrors.New(http.StatusInternalServerError, domain.ErrCodeUserInternalServerError, "Something went wrong", err)
	case err != nil:
		return errors.WithMessage(err, "UpdateUserCtrl")
	default:
		return nil
	}
}

func (ctrl User) DeleteUserCtrl(ctx context.Context, request domain.DeleteEmail) error {
	err := ctrl.srv.DeleteUserSrv(ctx, request)
	switch {
	case errors.Is(err, entity.ErrUserEmailNotFound):
		return apierrors.New(http.StatusNotFound, domain.ErrCodeUserEmailNotFound, "Email not found", err)
	case errors.Is(err, entity.ErrUserEmailInternalServerError):
		return apierrors.New(http.StatusInternalServerError, domain.ErrCodeUserInternalServerError, "Something went wrong", err)
	case err != nil:
		return errors.WithMessage(err, "DeleteUserCtrl")
	default:
		return nil
	}
}
