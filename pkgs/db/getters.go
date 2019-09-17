package db

import (
	"sort"

	"github.com/Soroka-EDMS/svc/users/pkgs/errors"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
)

func CheckAuthInDb(db *models.UsersDb, username, password string) error {
	db.Mtx.Lock()
	defer db.Mtx.Unlock()
	if username == "" || password == "" {
		return errors.ErrNonAuthorized
	}

	for _, record := range db.Creds {
		if record.Username == username && record.Password == password {
			return nil
		}
	}

	return errors.ErrNonAuthorized
}

func GetUserProfileFromDb(db *models.UsersDb, userID int, userEmail string) (res *models.UserProfileResp, err error) {
	db.Mtx.Lock()
	defer db.Mtx.Unlock()
	if userEmail == "" && userID == 0 {
		return nil, errors.ErrMalformedRequest
	} else if userEmail == "" && (userID > 0 && userID < len(db.Profiles)) {
		return GetById(db, userID)
	} else {
		return GetByEmail(db, userEmail)
	}
}

func GetById(db *models.UsersDb, userID int) (res *models.UserProfileResp, err error) {
	profile := db.Profiles[userID]
	return &profile, nil
}

func GetByEmail(db *models.UsersDb, userEmail string) (res *models.UserProfileResp, err error) {
	var profile models.UserProfileResp
	for _, record := range db.Profiles {
		if record.Email == userEmail {
			profile = record
			return &profile, nil
		}
	}

	return nil, errors.ErrProfileNotFound
}

//GetUserListFromDb returns list with users and pagination parameters according to pagination parameters in request
func GetUserListFromDb(db *models.UsersDb, offset, limit int, sortCriteria, order string) (res *models.UsersListResp, err error) {
	db.Mtx.Lock()
	defer db.Mtx.Unlock()
	var resp models.UsersListResp

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
