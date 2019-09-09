package db

import (
	"sort"

	"github.com/Soroka-EDMS/svc/users/pkgs/errors"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
)

func GetUserProfileFromDb(db *models.UsersDb, userID int, userEmail string) (res *models.UserProfile, err error) {
	if userEmail == "" && userID == 0 {
		return nil, errors.ErrMalformedRequest
	} else if userEmail == "" && (userID > 0 && userID < len(db.Profiles)) {
		return GetById(db, userID)
	} else {
		return GetByEmail(db, userEmail)
	}
}

func GetById(db *models.UsersDb, userID int) (res *models.UserProfile, err error) {
	profile := db.Profiles[userID]
	return &profile, nil
}

func GetByEmail(db *models.UsersDb, userEmail string) (res *models.UserProfile, err error) {
	var profile models.UserProfile
	for _, record := range db.Profiles {
		if record.Email == userEmail {
			profile = record
			return &profile, nil
		}
	}

	return nil, errors.ErrProfileNotFound
}

func GetUserListFromDb(db *models.UsersDb, offset, limit int, sortCriteria, order string) (res *models.UsersListResponse, err error) {
	var resp models.UsersListResponse

	//If offset is bigger than or equal to amount of record in database, then prepare empty response
	if offset >= len(db.Profiles) {
		PrepareEmptyResponse(&resp)
	} else {
		//Check for amount of required users, get user list, sort user list
		nUsers, left := GetAmountOfUsers(offset, limit, len(db.Profiles))
		userList := GetUserList(offset, nUsers, db)

		f := GetSortMethod(userList, sortCriteria, order)
		if f == nil {
			return nil, errors.ErrInvalidSortParameters
		}

		sort.Slice(userList, f)

		resp.Users = userList
		resp.Pagination = models.PaginationInfo{
			Issued: nUsers,
			Left:   left,
		}
	}

	return &resp, nil
}
