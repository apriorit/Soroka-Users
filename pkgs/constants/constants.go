//  constants.go
//  https://github.com/apriorit/Soroka-EDMS/svc/users/pkgs/constants
//
//  Created by Ivan Kashuba on 2019.09.03
//  Defines all known constant across all the packages
package constants

const (
	AuthEndpoint               string = "/users/check_auth"
	RoleEndpoint               string = "/users/change_role/{id}"
	ListEndpoint               string = "/users"
	GetProfileEndpoint         string = "/user"
	DisableUserEndpoint        string = "/users/disable/{id}"
	EnableUserEndpoint         string = "/users/enable/{id}"
	URIOnCheckTokenProfile     string = "https://sessions_sessions.service_1:443/session/check_token"
	IncorrectParameter         string = "Incorrect parameter %s"
	MissingBody                string = "Missing content"
	MissingPath                string = "Missing path"
	MalformedBody              string = "Malformed body"
	MalformedPath              string = "Malformed path"
	InvalidClaims              string = "Invalid claims"
	InvalidMaskType            string = "Mask claim has invalid type"
	EncodingErorr              string = "Enconding response error"
	MissingAuthorizationHeader string = "Authorization header required"
	RefreshTokenRequired       string = "Refresh token is missing"
	MissingPayload             string = "Missing payload"
	InvalidAuthorization       string = "No priveledge"
	RequestToSessionFailed     string = "Request to session service failed"
	RequiredAuthorized         string = "Authorization required"
	InvalidSortParameters      string = "Sort criteria is undefined"
	ProfileNotFound            string = "No profile for such user id%d"
	PublicKeyIsMissing         string = "Public key is missing"
	MalformedRequest           string = "Malformed request"
	DisabledPartialy           string = "Disabled partially"
	EnabledPartialy            string = "Enabled partially"
	ChangedPartialy            string = "Changed partially"
	ChangedSuccessfully        string = "Successfully changed"
	NoRecordChanged            string = "Nor record has been changed"
	UserDoesNotExist           string = "User does not exist"
	InvalidRole                string = "InvalidRole"
	BearerSchema               string = "Bearer "
	BasicSchema                string = "Basic "
	BearerSchemaRequired       string = "Authorization requires Bearer scheme"
	BasicSchemaRequired        string = "Authorization requires Basic scheme"
	Required                   string = "Required"
	Invalid                    string = "Invalid"
)
