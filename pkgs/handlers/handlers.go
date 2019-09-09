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
	ep "github.com/Soroka-EDMS/svc/users/pkgs/endpoints"
	er "github.com/Soroka-EDMS/svc/users/pkgs/errors"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joncalhoun/qson"
)

func NewHTTPHandler(endpoints ep.UsersEndpoints) http.Handler {
	methods := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(config.GetLogger().Logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST    /users/change_role                  changes roles for users with certain ids
	// GET     /users                      		   returns an array with users
	// GET     /user                   			   returns a profile for user with certain id or email
	// POST    /users/disable/{id}                 changes user status to disabled
	// POST    /users/enable/{id}                  changes user status to enabled

	methods.Methods("POST").Path(constants.RoleEndpoint).Handler(httptransport.NewServer(
		endpoints.ChangeRoleEndpoint,
		DecodeChangeRoleRequest,
		encodeChangeRoleResponse,
		options...,
	))

	methods.Methods("GET").Path(constants.ListEndpoint).Handler(httptransport.NewServer(
		endpoints.GetUserListEndpoint,
		DecodeGetUserListRequest,
		encodeGetUserListResponse,
		options...,
	))

	methods.Methods("GET").Path(constants.GetProfileEndpoint).Handler(httptransport.NewServer(
		endpoints.GetUserProfileEndpoint,
		DecodeGetUserProfileRequest,
		encodeGetUserProfileResponse,
		options...,
	))

	methods.Methods("POST").Path(constants.DisableUserEndpoint).Handler(httptransport.NewServer(
		endpoints.DisableUsersEndpoint,
		DecodeChangeUsersRequest,
		encodeChangeUsersResponse,
		options...,
	))

	methods.Methods("POST").Path(constants.EnableUserEndpoint).Handler(httptransport.NewServer(
		endpoints.EnableUsersEndpoint,
		DecodeChangeUsersRequest,
		encodeChangeUsersResponse,
		options...,
	))

	return methods
}

func DecodeChangeRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.ChangeRoleRequest

	//Parse request body
	if r.Body == nil {
		return nil, er.ErrMissingBodyContent
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Role == "" {
		return nil, er.ErrMalformedBodyContent
	}

	//Parse request path
	ids, err := GetIds(r.URL.Path)
	if err != nil {
		return nil, err
	}

	req.Ids = ids

	return ep.ChangeRoleRequest{Req: req}, nil
}

func DecodeGetUserListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.UsersListRequest
	queryPart := r.URL.RawQuery
	err := qson.Unmarshal(&req, queryPart)

	if err != nil {
		return nil, err
	}

	return ep.UsersListRequest{Req: req}, nil
}

func DecodeGetUserProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.UserProfileRequest
	queryPart := r.URL.RawQuery
	err := qson.Unmarshal(&req, queryPart)

	if err != nil {
		return nil, err
	}

	return ep.UserProfileRequest{Req: req}, nil
}

func DecodeChangeUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.UsersChangeStatusRequest

	//Parse request path
	ids, err := GetIds(r.URL.Path)
	if err != nil {
		return nil, err
	}

	req.Ids = ids

	return ep.ChangeUserStatusRequest{Req: req}, nil
}

func encodeChangeRoleResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(ep.ChangeRoleResponse)
	if !ok {
		return er.ErrEncodingResponse
	}

	err := e.Error()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.Res)
}

func encodeGetUserListResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(ep.UsersListResponse)
	if !ok {
		return er.ErrEncodingResponse
	}

	err := e.Error()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.Res)
}

func encodeGetUserProfileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(ep.UserProfileResponse)
	if !ok {
		return er.ErrEncodingResponse
	}

	err := e.Error()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.Res)
}

func encodeChangeUsersResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(ep.ChangeUserStatusResponse)
	if !ok {
		return er.ErrEncodingResponse
	}

	err := e.Error()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.Res)
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
	case er.ErrMissingBody, er.ErrMissingPath:
		return http.StatusBadRequest, constants.Required
	case er.ErrMalformedBodyContent, er.ErrMalformedPath:
		return http.StatusBadRequest, constants.Invalid
	case qson.ErrInvalidParam:
		return http.StatusBadRequest, err.Error()
	case er.ErrEncodingResponse:
		return http.StatusInternalServerError, err.Error()
	case er.ErrMissingPayload, er.ErrInvalidAuthorization:
		return http.StatusUnauthorized, constants.InvalidAuthorization
	/*Add other error cases*/
	default:
		return http.StatusInternalServerError, err.Error()
	}
}

// errorWrapper represents the JSON struct used in responses that
// specify an error.

type ErrorResponse struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}
