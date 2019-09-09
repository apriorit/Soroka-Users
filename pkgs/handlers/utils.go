package handlers

import (
	"strconv"
	"strings"

	"github.com/joncalhoun/qson"

	er "github.com/Soroka-EDMS/svc/users/pkgs/errors"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
)

func GetIds(path string) ([]int, error) {
	if path == "" {
		return nil, er.ErrMissingPath
	}

	parts := strings.Split(path, "/")
	splitedLastPart := strings.Split(parts[len(parts)-1], ",")

	if len(splitedLastPart) == 0 {
		return nil, er.ErrMalformedPath
	}

	ids := make([]int, len(splitedLastPart))

	for count := 0; count < len(ids); count++ {
		id, err := strconv.Atoi(splitedLastPart[count])

		if err != nil {
			return nil, er.ErrMalformedPath
		}

		ids[count] = id
	}

	return ids, nil
}

//GetPagination converts raw request query to models.UsersListRequest
func QueryToUsersListRequest(query string, req *models.UsersListRequest) error {

	err := qson.Unmarshal(&req, query)
	if err != nil {
		return err
	}

	return nil
}

func QueryToUserProfileRequest(query string, req *models.UserProfileRequest) error {
	err := qson.Unmarshal(&req, query)
	if err != nil {
		return err
	}

	return nil
}
