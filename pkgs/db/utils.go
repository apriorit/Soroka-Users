package db

import (
	"encoding/json"
	_ "encoding/json"

	"github.com/go-kit/kit/log"

	"github.com/Soroka-EDMS/svc/users/pkgs/constants"
	"github.com/Soroka-EDMS/svc/users/pkgs/errors"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
	m "github.com/Soroka-EDMS/svc/users/pkgs/models"
	"github.com/Soroka-EDMS/svc/users/pkgs/stub"
)

//1. Prepare stub database
//2. Connect to database

//PrepareStubDatabase fills database stub with records
func PrepareStubDatabase(db *m.UsersDb) error {
	var (
		err       error
		userCreds m.UserCredentials
		role      m.UserRole
		profile   m.UserProfileResp
	)
	roleMarker := [constants.DbSize]string{"default", "admin", "user", "ordinaryUser", "reducedUser"}

	//Create maps
	db.Creds = make(map[int]m.UserCredentials)
	db.Profiles = make(map[int]m.UserProfileResp)
	db.Roles = make(map[int]m.UserRole)

	//Fill fake database
	for count, value := range roleMarker {
		err = PrepareDatabaseRecord(value, &userCreds, &role, &profile)
		if err != nil {
			return err
		}

		db.Creds[count] = userCreds
		db.Roles[count] = role
		db.Profiles[count] = profile
	}

	return nil
}

//PrepareDatabaseRecord fills creds, role and profile according to a given marker that can be one of:  default, admin, user, ordinaryUser, reducedUser
func PrepareDatabaseRecord(marker string, creds *m.UserCredentials, role *m.UserRole, profile *m.UserProfileResp) (err error) {

	switch marker {
	case "default":
		rawCreds := []byte(stub.DefaultUserCreds)
		err = json.Unmarshal(rawCreds, &creds)
		if err != nil {
			return err
		}

		rawRole := []byte(stub.DefaultRole)
		err = json.Unmarshal(rawRole, &role)
		if err != nil {
			return err
		}

		rawProfile := []byte(stub.DefaultProfile)
		err = json.Unmarshal(rawProfile, &profile)
		if err != nil {
			return err
		}
	case "admin":
		rawCreds := []byte(stub.AdminCreds)
		err = json.Unmarshal(rawCreds, &creds)
		if err != nil {
			return err
		}

		rawRole := []byte(stub.AdminRole)
		err = json.Unmarshal(rawRole, &role)
		if err != nil {
			return err
		}

		rawProfile := []byte(stub.AdminProfile)
		err = json.Unmarshal(rawProfile, &profile)
		if err != nil {
			return err
		}
	case "user":
		rawCreds := []byte(stub.UserCreds)
		err = json.Unmarshal(rawCreds, &creds)
		if err != nil {
			return err
		}

		rawRole := []byte(stub.UserRole)
		err = json.Unmarshal(rawRole, &role)
		if err != nil {
			return err
		}

		rawProfile := []byte(stub.UserProfile)
		err = json.Unmarshal(rawProfile, &profile)
		if err != nil {
			return err
		}
	case "ordinaryUser":
		rawCreds := []byte(stub.OrdinaryUserCreds)
		err = json.Unmarshal(rawCreds, &creds)
		if err != nil {
			return err
		}

		rawRole := []byte(stub.OrdinaryUserRole)
		err = json.Unmarshal(rawRole, &role)
		if err != nil {
			return err
		}

		rawProfile := []byte(stub.OrdinaryUserProfile)
		err = json.Unmarshal(rawProfile, &profile)
		if err != nil {
			return err
		}
	case "reducedUser":
		rawCreds := []byte(stub.ReducedUserCreds)
		err = json.Unmarshal(rawCreds, &creds)
		if err != nil {
			return err
		}

		rawRole := []byte(stub.ReducedUserRole)
		err = json.Unmarshal(rawRole, &role)
		if err != nil {
			return err
		}

		rawProfile := []byte(stub.ReducedUserProfile)
		err = json.Unmarshal(rawProfile, &profile)
		if err != nil {
			return err
		}
	}

	return nil
}

//Connection returns object that enclosing database after connection to it
func Connection(logger log.Logger, conn string) (*UserDbConnection, error) {
	db := m.UsersDb{}
	if conn == "stub" {
		PrepareStubDatabase(&db)
	} else {
		//Real database
	}

	logger.Log("method", "Connection", "dbCredsLen", len(db.Creds), "dbRolesLen", len(db.Roles), "dbProfilesLen", len(db.Profiles))

	return &UserDbConnection{Db: &db, Logger: log.With(logger, "pkg", "connection")}, nil
}

//FindRole checks whether a given role is contained in the database
func FindRole(db *models.UsersDb, role string) error {
	db.Mtx.Lock()
	defer db.Mtx.Unlock()
	roles := db.Roles
	for _, r := range roles {
		if role == r.Name {
			return nil
		}
	}
	return errors.ErrInvalidRole
}

//FindMask returns mask value according to a given role value. Pair of role and mask is unique
func FindMask(db *models.UsersDb, role string) int64 {
	db.Mtx.Lock()
	defer db.Mtx.Unlock()
	for _, r := range db.Roles {
		if r.Name == role {
			return r.Mask
		}
	}

	return 0
}

//GetSortMethod returns sorting function according to a given sorting criteria
func GetSortMethod(userList []models.UserInfo, criteria, order string) func(left, right int) bool {
	switch criteria {
	case "user_id":
		if order == "up" {
			return func(left, right int) bool { return userList[left].User_id < userList[right].User_id }
		} else if order == "down" {
			return func(left, right int) bool { return userList[left].User_id > userList[right].User_id }
		} else {
			return nil
		}
	case "user_name":
		if order == "up" {
			return func(left, right int) bool { return userList[left].User_name < userList[right].User_name }
		} else if order == "down" {
			return func(left, right int) bool { return userList[left].User_name > userList[right].User_name }
		} else {
			return nil
		}
	case "role":
		if order == "up" {
			return func(left, right int) bool { return userList[left].Role < userList[right].Role }
		} else if order == "down" {
			return func(left, right int) bool { return userList[left].Role > userList[right].Role }
		} else {
			return nil
		}
	case "email":
		if order == "up" {
			return func(left, right int) bool { return userList[left].Email < userList[right].Email }
		} else if order == "down" {
			return func(left, right int) bool { return userList[left].Email > userList[right].Email }
		} else {
			return nil
		}
	case "creation_date":
		if order == "up" {
			return func(left, right int) bool { return userList[left].Creation_date < userList[right].Creation_date }
		} else if order == "down" {
			return func(left, right int) bool { return userList[left].Creation_date > userList[right].Creation_date }
		} else {
			return nil
		}
	}

	return nil
}

//GetAmountOfUsers returns a value that represents possible number of users  that can be returned according to given limit and offset
func GetAmountOfUsers(offset, limit, dbSize int) (nUsers, left int) {
	//Check for amount of required users
	if (limit - offset) <= (dbSize - offset) {
		nUsers = limit - offset
		left = dbSize - limit
	} else {
		nUsers = dbSize - offset
		left = 0
	}

	return
}

//GetUserList returns a list with users according to offset and amount of required users
func GetUserList(offset, nUsers int, db *models.UsersDb) []m.UserInfo {
	//Prepare user list
	userList := make([]models.UserInfo, nUsers)

	//Get users info
	for pos, count := 0, offset; count < nUsers+offset; count++ {
		profile := db.Profiles[count]

		userList[pos].User_id = count
		userList[pos].User_name = db.Creds[count].Username
		userList[pos].Email = profile.Email
		userList[pos].Location = profile.Location
		userList[pos].Role = profile.Role.Name
		userList[pos].Creation_date = profile.Creation_date
		userList[pos].Status = profile.Status
		pos++
	}

	return userList
}

//PrepareEmptyResponse returns empty users list
func PrepareEmptyResponse(resp *models.UsersListResp) {
	resp.Users = make([]models.UserInfo, 0)
	resp.Pagination = models.PaginationInfo{
		Issued: 0,
		Left:   0,
	}
}
