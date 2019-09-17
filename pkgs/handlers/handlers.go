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
	c "github.com/Soroka-EDMS/svc/users/pkgs/constants"
	ep "github.com/Soroka-EDMS/svc/users/pkgs/endpoints"
	e "github.com/Soroka-EDMS/svc/users/pkgs/errors"
	er "github.com/Soroka-EDMS/svc/users/pkgs/errors"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joncalhoun/qson"
)

//NewHTTPHandler creates http handler that enclosing all endpoint handlers
func NewHTTPHandler(endpoints ep.UsersEndpoints) http.Handler {
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
		endpoints.CheckAuthEndpoint,
		DecodeCheckAuthRequest,
		encodeCheckAuthResponse,
		options...,
	))

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

//DecodeCheckAuthRequest decodes raw request to corresponding service model
func DecodeCheckAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req models.UserCredentials
		ok  bool
	)

	if req.Username, req.Password, ok = r.BasicAuth(); !ok {
		return nil, er.ErrBearerSchemaRequired
	}

	return ep.CheckAuthRequest{Req: req}, nil
}

//DecodeChangeRoleRequest decodes raw request to corresponding service model
func DecodeChangeRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req models.ChangeRole
		err error
	)

	if err = CheckAuthorization(r, "changeRole"); err != nil {
		return nil, err
	}

	//Parse request body
	if r.Body == nil {
		config.GetLogger().Logger.Log("1", "1")
		return nil, er.ErrMissingBodyContent
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.GetLogger().Logger.Log("2", "2")
		return nil, err
	}
	if req.Role == "" {
		config.GetLogger().Logger.Log("3", "3")
		return nil, er.ErrMalformedBodyContent
	}

	//Parse request path
	ids, err := GetIds(r.URL.Path)
	if err != nil {
		config.GetLogger().Logger.Log("4", "4")
		return nil, err
	}

	req.Ids = ids

	return ep.ChangeRoleRequest{Req: req}, nil
}

//DecodeGetUserListRequest decodes raw request to corresponding service model
func DecodeGetUserListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req models.UsersList
		err error
	)

	if err = CheckAuthorization(r, "userList"); err != nil {
		return nil, err
	}
	//Parse request query
	queryPart := r.URL.RawQuery
	if err = qson.Unmarshal(&req, queryPart); err != nil {
		return nil, err
	}

	return ep.UsersListRequest{Req: req}, nil
}

//DecodeGetUserProfileRequest decodes raw request to corresponding service model
func DecodeGetUserProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req models.UserProfileReq
		err error
	)

	config.GetLogger().Logger.Log("method", "DecodeGetUserProfileRequest", "URI", r.RequestURI)

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, e.ErrMissingAuthorizationHeader
	}

	config.GetLogger().Logger.Log("method", "DecodeGetUserProfileRequest", "secret", authHeader[len(c.BearerSchema):])
	if authHeader[len(c.BearerSchema):] != config.GetSecretString() {
		if err = CheckAuthorization(r, "userProfile"); err != nil {
			return nil, err
		}
	}

	//Parse request query
	queryPart := r.URL.RawQuery
	if err := qson.Unmarshal(&req, queryPart); err != nil {
		return nil, err
	}

	return ep.UserProfileRequest{Req: req}, nil
}

//DecodeChangeUsersRequest decodes raw request to corresponding service model
func DecodeChangeUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req models.UsersChangeStatus
		err error
	)

	if err = CheckAuthorization(r, "changeUserStatus"); err != nil {
		return nil, err
	}
	//Parse request path
	ids, err := GetIds(r.URL.Path)
	if err != nil {
		return nil, err
	}
	req.Ids = ids

	return ep.ChangeUserStatusRequest{Req: req}, nil
}

func encodeCheckAuthResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(ep.CheckAuthResponse)
	if !ok {
		return er.ErrEncodingResponse
	}

	if err := e.Error(); err != nil {
		return err
	}

	return nil
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
		return http.StatusBadRequest, c.Required
	case er.ErrMalformedBodyContent, er.ErrMalformedPath:
		return http.StatusBadRequest, c.Invalid
	case er.ErrMissingAuthorizationHeader:
		return http.StatusUnauthorized, c.MissingAuthorizationHeader
	case er.ErrBearerSchemaRequired:
		return http.StatusUnauthorized, c.BearerSchemaRequired
	case er.ErrBearerSchemaRequired:
		return http.StatusUnauthorized, c.BasicSchemaRequired
	case er.ErrMissingPayload, er.ErrInvalidAuthorization:
		return http.StatusUnauthorized, c.InvalidAuthorization
	case qson.ErrInvalidParam:
		return http.StatusBadRequest, err.Error()
	case er.ErrEncodingResponse:
		return http.StatusInternalServerError, err.Error()
	case er.ErrInvalidClaims:
		return http.StatusUnauthorized, err.Error()
	default:
		return http.StatusInternalServerError, err.Error()
	}
}

type ErrorResponse struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}
