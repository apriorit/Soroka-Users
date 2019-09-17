package db

import (
	"github.com/Soroka-EDMS/svc/users/pkgs/config"
	"github.com/Soroka-EDMS/svc/users/pkgs/errors"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
)

//ChangeRoleInDb changes role for current user
func ChangeRoleInDb(db *models.UsersDb, userID int, role string) error {
	if 0 >= userID || userID >= len(db.Profiles) {
		return errors.ErrUserDoesNotExist
	}

	//Find role
	err := FindRole(db, role)

	if err != nil {
		return err
	}

	//Define mask for this role
	mask := FindMask(db, role)
	//Change profile
	newProfile := db.Profiles[userID]
	newProfile.Role.Name = role
	newProfile.Role.Mask = mask
	db.Profiles[userID] = newProfile

	return nil
}

//ChangeUserStatusInDb blocks/unblocks user according to status value
func ChangeUserStatusInDb(db *models.UsersDb, userID int, status bool) error {
	db.Mtx.Lock()
	defer db.Mtx.Unlock()
	if 0 >= userID || userID >= len(db.Profiles) {
		return errors.ErrUserDoesNotExist
	}

	config.GetLogger().Logger.Log("method", "ChangeUserStatusInDb", "neaded status", status)
	newProfile := db.Profiles[userID]
	newProfile.Status = status
	db.Profiles[userID] = newProfile

	return nil
}
