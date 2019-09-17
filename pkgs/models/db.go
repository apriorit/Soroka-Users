package models

import (
	"context"
	"sync"
)

type UsersDatabase interface {
	CheckAuth(ctx context.Context, username, password string) (err error)
	ChangeRole(ctx context.Context, userID int, role string) (err error)
	GetUserProfile(ctx context.Context, userID int, userEmail string) (res *UserProfileResp, err error)
	GetUsersList(ctx context.Context, offset, limit int, sort, order string) (res *UsersListResp, err error)
	ChangeUserStatus(ctx context.Context, userID int, status bool) (err error)
}

type UserCredentials struct {
	Username string
	Password string
}

type UserRole struct {
	Name string
	Mask int64
}

type UsersDb struct {
	Creds    map[int]UserCredentials
	Roles    map[int]UserRole
	Profiles map[int]UserProfileResp
	Mtx      sync.RWMutex
}
