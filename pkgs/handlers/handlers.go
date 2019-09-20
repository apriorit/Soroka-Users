//  handlers.go
//  https://github.com/apriorit/Soroka-EDMS/svc/users/pkgs/handlers
//
//  Created by Ivan Kashuba on 2019.09.03
//  Describe HTTP handlers

package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Soroka-EDMS/svc/users/pkgs/config"
	"github.com/Soroka-EDMS/svc/users/pkgs/constants"
	"github.com/Soroka-EDMS/svc/users/pkgs/endpoints"
	"github.com/Soroka-EDMS/svc/users/pkgs/errors"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joncalhoun/qson"
)

//NewHTTPHandler creates http handler that enclosing all endpoint handlers
func NewHTTPHandler(eps endpoints.IUsersEndpoints) http.Handler {
	methods := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(config.GetLogger().Logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// GET     /users/check_auth 				   checks validness of a given user credentials
	// POST    /users/change_role                  changes roles for users with certain ids
	// GET     /users                      		   returns an array with users
	// GET     /user                   			   returns a profile for user with certain id or email
	// POST    /users/disable/{id}                 changes user status to disabled
	// POST    /users/enable/{id}                  changes user status to enabled
	methods.Methods("GET").Path(constants.AuthEndpoint).Handler(httptransport.NewServer(
		eps.GetEndpoint(endpoints.CheckAuth),
		eps.DecodeCheckAuthRequest,
		eps.EncodeCheckAuthResponse,
		options...,
	))

	methods.Methods("POST").Path(constants.RoleEndpoint).Handler(httptransport.NewServer(
		eps.GetEndpoint(endpoints.ChangeRole),
		eps.DecodeChangeRoleRequest,
		eps.EncodeChangeRoleResponse,
		options...,
	))

	methods.Methods("GET").Path(constants.ListEndpoint).Handler(httptransport.NewServer(
		eps.GetEndpoint(endpoints.GetUserList),
		eps.DecodeGetUserListRequest,
		eps.EncodeGetUserListResponse,
		options...,
	))

	methods.Methods("GET").Path(constants.GetProfileEndpoint).Handler(httptransport.NewServer(
		eps.GetEndpoint(endpoints.GetUserProfile),
		eps.DecodeGetUserProfileRequest,
		eps.EncodeGetUserProfileResponse,
		options...,
	))

	methods.Methods("GET").Path(constants.DisableUserEndpoint).Handler(httptransport.NewServer(
		eps.GetEndpoint(endpoints.DisableUsers),
		eps.DecodeChangeUsersRequest,
		eps.EncodeChangeUsersResponse,
		options...,
	))

	methods.Methods("GET").Path(constants.EnableUserEndpoint).Handler(httptransport.NewServer(
		eps.GetEndpoint(endpoints.EnableUsers),
		eps.DecodeChangeUsersRequest,
		eps.EncodeChangeUsersResponse,
		options...,
	))

	return methods
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}

	errCode, errReason := codeFrom(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Reason:  errReason,
		Message: err.Error(),
	})
}

func codeFrom(err error) (int, string) {
	switch err {
	case errors.ErrMissingBody, errors.ErrMissingPath:
		return http.StatusBadRequest, constants.Required
	case errors.ErrMalformedBodyContent, errors.ErrMalformedPath:
		return http.StatusBadRequest, constants.Invalid
	case errors.ErrMissingAuthorizationHeader:
		return http.StatusUnauthorized, constants.MissingAuthorizationHeader
	case errors.ErrBearerSchemaRequired:
		return http.StatusUnauthorized, constants.BearerSchemaRequired
	case errors.ErrBearerSchemaRequired:
		return http.StatusUnauthorized, constants.BasicSchemaRequired
	case errors.ErrMissingPayload, errors.ErrInvalidAuthorization:
		return http.StatusUnauthorized, constants.InvalidAuthorization
	case qson.ErrInvalidParam:
		return http.StatusBadRequest, err.Error()
	case errors.ErrEncodingResponse:
		return http.StatusInternalServerError, err.Error()
	case errors.ErrInvalidClaims:
		return http.StatusUnauthorized, err.Error()
	default:
		return http.StatusInternalServerError, err.Error()
	}
}

type ErrorResponse struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}
