package db

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/Soroka-EDMS/svc/users/pkgs/models"
)

type UserDbConnection struct {
	Db     *models.UsersDb
	Logger log.Logger
}

//CheckAuth calls CheckAuthInDb described in getters.go
func (db *UserDbConnection) CheckAuth(ctx context.Context, username, password string) error {
	db.Logger.Log("method", "main", "dbCredsLen", len(db.Db.Creds), "dbRolesLen", len(db.Db.Roles), "dbProfilesLen", len(db.Db.Profiles))
	return CheckAuthInDb(db.Db, username, password)
}

//ChangeRole calls ChangeRoleInDb described in posters.go
func (db *UserDbConnection) ChangeRole(ctx context.Context, userID int, role string) error {
	return ChangeRoleInDb(db.Db, userID, role)
}

//GetUserProfile calls GetUserProfileFromDb described in getters.go
func (db *UserDbConnection) GetUserProfile(ctx context.Context, userID int, userEmail string) (res *models.UserProfileResp, err error) {
	return GetUserProfileFromDb(db.Db, userID, userEmail)
}

//GetUsersList calls GetUserListFromDb described in getters.go
func (db *UserDbConnection) GetUsersList(ctx context.Context, offset, limit int, sort, order string) (res *models.UsersListResp, err error) {
	return GetUserListFromDb(db.Db, offset, limit, sort, order)
}

//ChangeUserStatus calls ChangeUserStatusInDb described in posters.go
func (db *UserDbConnection) ChangeUserStatus(ctx context.Context, userID int, status bool) error {
	return ChangeUserStatusInDb(db.Db, userID, status)
}
