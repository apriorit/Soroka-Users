package db

import (
	"github.com/Soroka-EDMS/svc/users/pkgs/config"
	"github.com/Soroka-EDMS/svc/users/pkgs/errors"
)

//ChangeRoleInDb changes role for current user
func (db *UserDbStub) ChangeRoleInDb(userID int, role string) error {
	if 0 >= userID || userID >= len(db.Profiles) {
		return errors.ErrUserDoesNotExist
	}

	//Find role
	err := db.FindRole(role)

	if err != nil {
		return err
	}

	//Define mask for this role
	mask := db.FindMask(role)
	//Change profile
	newProfile := db.Profiles[userID]
	newProfile.Role.Name = role
	newProfile.Role.Mask = mask
	db.Profiles[userID] = newProfile

	return nil
}

//ChangeUserStatusInDb blocks/unblocks user according to status value
func (db *UserDbStub) ChangeUserStatusInDb(userID int, status bool) error {
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
