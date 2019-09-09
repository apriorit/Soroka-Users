//  constants.go
//  https://github.com/apriorit/Soroka-EDMS/svc/users/pkgs/constants
//
//  Created by Ivan Kashuba on 2019.09.03
//  Defines all known constant across all the packages
package constants

const (
	RoleEndpoint          string = "/users/change_role"
	ListEndpoint          string = "/users"
	GetProfileEndpoint    string = "/user"
	DisableUserEndpoint   string = "/users/disable/{id}"
	EnableUserEndpoint    string = "/users/enable/{id}"
	RegisterEndpoint      string = "/users/register"
	IncorrectParameter    string = "Incorrect parameter %s"
	MissingBody           string = "Missing content"
	MissingPath           string = "Missing path"
	MalformedBody         string = "Malformed body"
	MalformedPath         string = "Malformed path"
	EncodingErorr         string = "Enconding response error"
	MissingPayload        string = "Missing payload"
	InvalidAuthorization  string = "No priveledge"
	InvalidSortParameters string = "Sort criteria is undefined"
	ProfileNotFound       string = "No profile for such user id%d"
	MalformedRequest      string = "Malformed request"
	DisabledPartialy      string = "Disabled partialy"
	EnabledPartialy       string = "Enabled partialy"
	ChangedPartialy       string = "Changed partialy"
	ChangedSuccessfully   string = "Successfully changed"
	NoRecordChanged       string = "Nor record has been changed"
	UserDoesNotExist      string = "User does not exist"
	InvalidRole           string = "InvalidRole"
	Invalid               string = "Invalid"
	Required              string = "Required"
	UserName_admin        string = "admin"
	UserName_user         string = "user"
	UserName_ordinaryUser string = "ordinaryUser"
	UserName_reducedUser  string = "reducedUser"
	Password_admin        string = "a@m1n"
	Password_user         string = "@Sr!1"
	Password_ordinaryUser string = "@Sua1pwd"
	Password_reducedUser  string = "Canguro!1"
	RoleName_admin        string = "admin"
	RoleName_user         string = "user"
	RoleName_ordinaryUser string = "ordinaryUser"
	RoleName_reducedUser  string = "reducedUser"
	RoleMask_admin        int64  = 32767
	RoleMask_user         int64  = 25588
	RoleMask_ordinaryUser int64  = 25524
	RoleMask_reducedUser  int64  = 8628
	DbSize                int    = 5
	ChangeRole            int    = 12
	ChangeUserStatus      int    = 13
	GetProfile            int    = 14
	QueryUsers            int    = 15
	NumOfPriveledges      int    = 5
)
