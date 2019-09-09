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

func (db *UserDbConnection) ChangeRole(ctx context.Context, userID int, role string) error {
	return ChangeRoleInDb(db.Db, userID, role)
}
func (db *UserDbConnection) GetUserProfile(ctx context.Context, userID int, userEmail string) (res *models.UserProfile, err error) {
	return GetUserProfileFromDb(db.Db, userID, userEmail)
}
func (db *UserDbConnection) GetUsersList(ctx context.Context, offset, limit int, sort, order string) (res *models.UsersListResponse, err error) {
	return GetUserListFromDb(db.Db, offset, limit, sort, order)
}
func (db *UserDbConnection) ChangeUserStatus(ctx context.Context, userID int, status bool) error {
	return ChangeUserStatusInDb(db.Db, userID, status)
}
