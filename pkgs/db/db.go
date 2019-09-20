package db

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"

	"github.com/Soroka-EDMS/svc/users/pkgs/models"
)

type UserDbStub struct {
	Creds    map[int]models.UserCredentials
	Roles    map[int]models.UserRole
	Profiles map[int]models.UserProfile
	Mtx      sync.RWMutex
	Logger   log.Logger
}

//CheckAuth calls CheckAuthInDb described in getters.go
func (db *UserDbStub) CheckAuth(ctx context.Context, username, password string) error {
	db.Logger.Log("method", "main", "dbCredsLen", len(db.Creds), "dbRolesLen", len(db.Roles), "dbProfilesLen", len(db.Profiles))
	return db.CheckAuthInDb(username, password)
}

//ChangeRole calls ChangeRoleInDb described in posters.go
func (db *UserDbStub) ChangeRole(ctx context.Context, userID int, role string) error {
	return db.ChangeRoleInDb(userID, role)
}

//GetUserProfile calls GetUserProfileFromDb described in getters.go
func (db *UserDbStub) GetUserProfile(ctx context.Context, userID int, userEmail string) (res *models.UserProfile, err error) {
	return db.GetUserProfileFromDb(userID, userEmail)
}

//GetUsersList calls GetUserListFromDb described in getters.go
func (db *UserDbStub) GetUsersList(ctx context.Context, offset, limit int, sort, order string) (res *models.UsersListResp, err error) {
	return db.GetUserListFromDb(offset, limit, sort, order)
}

//ChangeUserStatus calls ChangeUserStatusInDb described in posters.go
func (db *UserDbStub) ChangeUserStatus(ctx context.Context, userID int, status bool) error {
	return db.ChangeUserStatusInDb(userID, status)
}
