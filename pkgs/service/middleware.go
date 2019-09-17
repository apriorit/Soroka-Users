package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	m "github.com/Soroka-EDMS/svc/users/pkgs/models"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(m.UsersService) m.UsersService

type loggingMiddleware struct {
	next   m.UsersService
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next m.UsersService) m.UsersService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (lmw loggingMiddleware) CheckAuth(ctx context.Context, r m.UserCredentials) (err error) {
	defer func(begin time.Time) {
		lmw.logger.Log("method", "CheckAuth", "usename", r.Username, "password", r.Password, "took", time.Since(begin), "err", err)
	}(time.Now())
	return lmw.next.CheckAuth(ctx, r)
}

func (lmw loggingMiddleware) ChangeRole(ctx context.Context, r m.ChangeRole) (res m.ChangeUsers, err error) {
	defer func(begin time.Time) {
		lmw.logger.Log("method", "ChangeRole", "took", time.Since(begin), "err", err)
	}(time.Now())
	return lmw.next.ChangeRole(ctx, r)
}

func (lmw loggingMiddleware) GetUserList(ctx context.Context, r m.UsersList) (res m.UsersListResp, err error) {
	defer func(begin time.Time) {
		lmw.logger.Log("method", "GetUserList", "offset", r.Offset, "limit", r.Limit, "Sort", r.Sort, "Order", r.Order, "took", time.Since(begin), "err", err)
	}(time.Now())
	return lmw.next.GetUserList(ctx, r)
}

func (lmw loggingMiddleware) GetUserProfile(ctx context.Context, r m.UserProfileReq) (res m.UserProfileResp, err error) {
	defer func(begin time.Time) {
		lmw.logger.Log("method", "GetUserProfile", "id", r.Email, "email", r.Email, "took", time.Since(begin), "err", err)
	}(time.Now())
	return lmw.next.GetUserProfile(ctx, r)
}

func (lmw loggingMiddleware) DisableUsers(ctx context.Context, r m.UsersChangeStatus) (res m.ChangeUsers, err error) {
	defer func(begin time.Time) {
		lmw.logger.Log("method", "DisableUsers", "took", time.Since(begin), "err", err)
	}(time.Now())
	return lmw.next.DisableUsers(ctx, r)
}

func (lmw loggingMiddleware) EnableUsers(ctx context.Context, r m.UsersChangeStatus) (res m.ChangeUsers, err error) {
	defer func(begin time.Time) {
		lmw.logger.Log("method", "EnableUsers", "took", time.Since(begin), "err", err)
	}(time.Now())
	return lmw.next.EnableUsers(ctx, r)
}
