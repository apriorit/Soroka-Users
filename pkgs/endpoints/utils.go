package endpoints

import (
	"net/http"

	"strconv"
	"strings"

	"github.com/joncalhoun/qson"

	"github.com/Soroka-EDMS/svc/users/pkgs/config"
	"github.com/Soroka-EDMS/svc/users/pkgs/constants"
	"github.com/Soroka-EDMS/svc/users/pkgs/errors"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
)

//GetIds returns slice with ids that enclosed in request path
func GetIds(path string) ([]int, error) {
	if path == "" {
		return nil, errors.ErrMissingPath
	}

	parts := strings.Split(path, "/")
	splitedLastPart := strings.Split(parts[len(parts)-1], ",")

	if len(splitedLastPart) == 0 {
		return nil, errors.ErrMalformedPath
	}

	ids := make([]int, len(splitedLastPart))

	for count := 0; count < len(ids); count++ {
		id, err := strconv.Atoi(splitedLastPart[count])

		if err != nil {
			return nil, errors.ErrMalformedPath
		}

		ids[count] = id
	}

	return ids, nil
}

//QueryToUsersListRequest converts query string from request to models.UsersList that describes pagination parameters
func QueryToUsersListRequest(query string, req *models.UsersList) error {

	err := qson.Unmarshal(&req, query)
	if err != nil {
		return err
	}

	return nil
}

//QueryToUserProfileRequest converts query string from request to models.UserProfileReq
func QueryToUserProfileRequest(query string, req *models.UserProfileReq) error {
	err := qson.Unmarshal(&req, query)
	if err != nil {
		return err
	}

	return nil
}

//GetTokensPair returns a pair of tokens (access and refresh) or error if one of token is missing
func GetTokensPair(req *http.Request) (res models.TokensPair, err error) {
	//Get access token string
	authHeader := req.Header.Get("Authorization")
	config.GetLogger().Logger.Log("method", "GetTokensPair", "authHeader", authHeader)
	if authHeader == "" {
		return res, errors.ErrMissingAuthorizationHeader
	}

	if !strings.HasPrefix(authHeader, constants.BearerSchema) {
		return res, errors.ErrBearerSchemaRequired
	}
	//Get refresh token from cookie
	cookie, err := req.Cookie("refresh_token")
	if err != nil {
		return res, errors.ErrRefreshTokenRequired
	}

	res.AccessToken = authHeader[len(constants.BearerSchema):]
	res.RefreshToken = cookie.Value

	return res, nil
}

func CheckPriveledge(mask int64, method string) bool {
	if method == "" || mask == 0 {
		return false
	}

	privMap := config.GetPrivileges().Privileges
	bitValue := mask >> (uint(privMap[method] - 1)) & 1
	return bitValue == 1
}
