package endpoints

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/joncalhoun/qson"

	"github.com/Soroka-EDMS/svc/users/pkgs/config"
	"github.com/Soroka-EDMS/svc/users/pkgs/constants"
	"github.com/Soroka-EDMS/svc/users/pkgs/errors"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
)

type EndpointType int

const (
	CheckAuth EndpointType = iota
	ChangeRole
	GetUserList
	GetUserProfile
	DisableUsers
	EnableUsers
)

type IUsersEndpoints interface {
	DecodeCheckAuthRequest(context.Context, *http.Request) (interface{}, error)
	DecodeChangeRoleRequest(context.Context, *http.Request) (interface{}, error)
	DecodeGetUserListRequest(context.Context, *http.Request) (interface{}, error)
	DecodeGetUserProfileRequest(context.Context, *http.Request) (interface{}, error)
	DecodeChangeUsersRequest(context.Context, *http.Request) (interface{}, error)
	EncodeCheckAuthResponse(context.Context, http.ResponseWriter, interface{}) error
	EncodeChangeRoleResponse(context.Context, http.ResponseWriter, interface{}) error
	EncodeGetUserListResponse(context.Context, http.ResponseWriter, interface{}) error
	EncodeGetUserProfileResponse(context.Context, http.ResponseWriter, interface{}) error
	EncodeChangeUsersResponse(context.Context, http.ResponseWriter, interface{}) error
	getPKey() []byte
	getSignSecret() []byte
	GetEndpoint(EndpointType) endpoint.Endpoint
}

type UsersEndpoints struct {
	PKey                   []byte
	signSecret             []byte
	CheckAuthEndpoint      endpoint.Endpoint
	ChangeRoleEndpoint     endpoint.Endpoint
	GetUserListEndpoint    endpoint.Endpoint
	GetUserProfileEndpoint endpoint.Endpoint
	DisableUsersEndpoint   endpoint.Endpoint
	EnableUsersEndpoint    endpoint.Endpoint
}

//Build creates endpoints with authorization by priveledges
func Build(logger log.Logger, service models.UsersService, certKeyData, secret []byte) IUsersEndpoints {
	return UsersEndpoints{
		PKey:                   certKeyData,
		signSecret:             secret,
		CheckAuthEndpoint:      BuildCheckAuthEndpoint(service),
		ChangeRoleEndpoint:     BuildChangeRoleEndpoint(service),
		GetUserListEndpoint:    BuildGetUserListEndpoint(service),
		GetUserProfileEndpoint: BuildGetUserProfileEndpoint(service),
		DisableUsersEndpoint:   BuildDisableUsersEndpoint(service),
		EnableUsersEndpoint:    BuildEnableUsersEndpoint(service),
	}
}

func (eps UsersEndpoints) getPKey() []byte {
	return eps.PKey
}
func (eps UsersEndpoints) getSignSecret() []byte {
	return eps.signSecret
}

func (eps UsersEndpoints) GetEndpoint(et EndpointType) (ep endpoint.Endpoint) {
	switch et {
	case CheckAuth:
		ep = eps.CheckAuthEndpoint
	case ChangeRole:
		ep = eps.ChangeRoleEndpoint
	case GetUserList:
		ep = eps.GetUserListEndpoint
	case GetUserProfile:
		ep = eps.GetUserProfileEndpoint
	case DisableUsers:
		ep = eps.DisableUsersEndpoint
	case EnableUsers:
		ep = eps.EnableUsersEndpoint
	}

	return ep
}

func (eps UsersEndpoints) DecodeCheckAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req models.UserCredentials
		ok  bool
	)

	if req.Username, req.Password, ok = r.BasicAuth(); !ok {
		return nil, errors.ErrBearerSchemaRequired
	}

	return CheckAuthRequest{Req: req}, nil
}

//DecodeChangeRoleRequest decodes raw request to corresponding service model
func (eps UsersEndpoints) DecodeChangeRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req models.ChangeRole
		err error
	)

	//Pass public key
	if err = eps.CheckAuthorization(r, "changeRole"); err != nil {
		return nil, err
	}

	//Parse request body
	if r.Body == nil {
		return nil, errors.ErrMissingBodyContent
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Role == "" {
		return nil, errors.ErrMalformedBodyContent
	}

	//Parse request path
	ids, err := GetIds(r.URL.Path)
	if err != nil {
		return nil, err
	}

	req.Ids = ids

	return ChangeRoleRequest{Req: req}, nil
}

//DecodeGetUserListRequest decodes raw request to corresponding service model
func (eps UsersEndpoints) DecodeGetUserListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req models.UsersList
		err error
	)

	if err = eps.CheckAuthorization(r, "userList"); err != nil {
		return nil, err
	}
	//Parse request query
	queryPart := r.URL.RawQuery
	if err = qson.Unmarshal(&req, queryPart); err != nil {
		return nil, err
	}

	return UsersListRequest{Req: req}, nil
}

//DecodeGetUserProfileRequest decodes raw request to corresponding service model
func (eps UsersEndpoints) DecodeGetUserProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req models.UserProfileReq
		err error
	)

	config.GetLogger().Logger.Log("method", "DecodeGetUserProfileRequest", "URI", r.RequestURI)

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.ErrMissingAuthorizationHeader
	}

	config.GetLogger().Logger.Log("method", "DecodeGetUserProfileRequest", "secret", authHeader[len(constants.BearerSchema):])
	if authHeader[len(constants.BearerSchema):] != string(eps.getSignSecret()) {
		if err = eps.CheckAuthorization(r, "userProfile"); err != nil {
			return nil, err
		}
	}

	//Parse request query
	queryPart := r.URL.RawQuery
	if err := qson.Unmarshal(&req, queryPart); err != nil {
		return nil, err
	}

	return UserProfileRequest{Req: req}, nil
}

//DecodeChangeUsersRequest decodes raw request to corresponding service model
func (eps UsersEndpoints) DecodeChangeUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req models.UsersChangeStatus
		err error
	)

	if err = eps.CheckAuthorization(r, "changeUserStatus"); err != nil {
		return nil, err
	}
	//Parse request path
	ids, err := GetIds(r.URL.Path)
	if err != nil {
		return nil, err
	}
	req.Ids = ids

	return ChangeUserStatusRequest{Req: req}, nil
}

func (eps UsersEndpoints) EncodeCheckAuthResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(CheckAuthResponse)
	if !ok {
		return errors.ErrEncodingResponse
	}

	if err := e.Error(); err != nil {
		return err
	}

	return nil
}

func (eps UsersEndpoints) EncodeChangeRoleResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(ChangeRoleResponse)
	if !ok {
		return errors.ErrEncodingResponse
	}

	err := e.Error()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.Res)
}

func (eps UsersEndpoints) EncodeGetUserListResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(UsersListResponse)
	if !ok {
		return errors.ErrEncodingResponse
	}

	err := e.Error()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.Res)
}

func (eps UsersEndpoints) EncodeGetUserProfileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(UserProfileResponse)
	if !ok {
		return errors.ErrEncodingResponse
	}

	err := e.Error()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.Res)
}

func (eps UsersEndpoints) EncodeChangeUsersResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(ChangeUserStatusResponse)
	if !ok {
		return errors.ErrEncodingResponse
	}

	err := e.Error()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.Res)
}

//CheckAuthorization parse pair of tokens from request, checks their validness by calling sessions service and return an access token.
func (eps UsersEndpoints) CheckAuthorization(rawRequest *http.Request, method string) error {
	var (
		err   error
		tPair models.TokensPair
		at    models.AccessToken
	)

	//Get Tokens pair from raw request
	if tPair, err = GetTokensPair(rawRequest); err != nil {
		return err
	}
	//Check_token_validness
	if at, err = eps.CheckJWT(tPair); err != nil {
		return err
	}

	//ParseClaims
	if err = eps.ParseClaims(method, at.AccessToken); err != nil {
		return err
	}

	return nil
}

//ParseClaims gets claims and checks client permissions
func (eps UsersEndpoints) ParseClaims(method, token string) error {
	//Get raw token
	rawToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		//return session secret that was used in signing process
		return eps.getSignSecret(), nil
	})

	if err != nil {
		return err
	}

	//Check token validness
	var claims jwt.MapClaims
	claims, ok := rawToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.ErrInvalidClaims
	}
	mask, ok := claims["mask"].(float64)
	if !ok {
		return errors.ErrInvalidMaskType
	}
	//Check priveldge
	if !CheckPriveledge(int64(mask), method) {
		return errors.ErrInvalidAuthorization
	}

	return nil
}

//CheckJWT calls  check_token endpoint of sessions service in order to update or confirm token validness
func (eps UsersEndpoints) CheckJWT(tp models.TokensPair) (res models.AccessToken, err error) {
	//Prepare request with access token in body and refresh token in cookie
	rawTokenString := `{"access_token": "%s"}`
	req, err := http.NewRequest("POST", constants.URIOnCheckTokenProfile, bytes.NewBuffer([]byte(fmt.Sprintf(rawTokenString, tp.AccessToken))))
	if err != nil {
		return res, err
	}

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    tp.RefreshToken,
		HttpOnly: true,
		Secure:   true,
	}
	req.AddCookie(cookie)

	//Create http client with public key
	netClient, err := eps.MakeHTTPClient()

	if err != nil {
		return res, err
	}

	//Do request
	resp, err := netClient.Do(req)

	if err != nil {
		return res, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusBadRequest:
		config.GetLogger().Logger.Log("err", "2")
		return res, errors.ErrRequestToSessionFailed
	case http.StatusUnauthorized:
		return res, errors.ErrNonAuthorized
	default:
		config.GetLogger().Logger.Log("code", resp.StatusCode)
		return res, fmt.Errorf("Check token request failed with code: %v", resp.StatusCode)
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		return res, fmt.Errorf("Check token received content type: %v", resp.Header.Get("Content-Type"))
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return res, fmt.Errorf("Desponse decoding response failed on check token ")
	}

	return res, nil
}

func (eps UsersEndpoints) MakeHTTPClient() (*http.Client, error) {
	//Get public key
	pKey := eps.getPKey()

	if len(pKey) == 0 {
		return nil, errors.ErrPublicKeyIsMissing
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(pKey)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				InsecureSkipVerify: true,
			},
		},
	}

	return client, nil
}

//BuildChangeRoleEndpoint returns an endpoint that call CheckAuth handler.
func BuildCheckAuthEndpoint(svc models.UsersService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CheckAuthRequest)
		err := svc.CheckAuth(ctx, req.Req)
		return CheckAuthResponse{Err: err}, nil
	}
}

// BuildChangeRoleEndpoint returns an endpoint that call ChangeRole handler.
func BuildChangeRoleEndpoint(svc models.UsersService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeRoleRequest)
		res, err := svc.ChangeRole(ctx, req.Req)
		return ChangeRoleResponse{Res: res, Err: err}, nil
	}
}

// BuildGetUserListEndpoint returns an endpoint that call GetUserList handler.
func BuildGetUserListEndpoint(svc models.UsersService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UsersListRequest)
		res, err := svc.GetUserList(ctx, req.Req)
		return UsersListResponse{Res: res, Err: err}, nil
	}
}

// BuildGetUserProfileEndpoint returns an endpoint that call GetUserProfile handler.
func BuildGetUserProfileEndpoint(svc models.UsersService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UserProfileRequest)
		res, err := svc.GetUserProfile(ctx, req.Req)
		return UserProfileResponse{Res: res, Err: err}, nil
	}
}

// BuildDisableUsersEndpoint returns an endpoint that call DisableUsers handler.
func BuildDisableUsersEndpoint(svc models.UsersService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeUserStatusRequest)
		res, err := svc.DisableUsers(ctx, req.Req)
		return ChangeUserStatusResponse{Res: res, Err: err}, nil
	}
}

// BuildEnableUsersEndpoint returns an endpoint that call EnableUsers handler.
func BuildEnableUsersEndpoint(svc models.UsersService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeUserStatusRequest)
		res, err := svc.EnableUsers(ctx, req.Req)
		return ChangeUserStatusResponse{Res: res, Err: err}, nil
	}
}

type CheckAuthRequest struct {
	Req models.UserCredentials
}

type CheckAuthResponse struct {
	Err error
}

type ChangeRoleRequest struct {
	Req models.ChangeRole
}

type ChangeRoleResponse struct {
	Res models.ChangeUsers
	Err error
}

type UsersListRequest struct {
	Req models.UsersList
}

type UsersListResponse struct {
	Res models.UsersListResp
	Err error
}

type UserProfileRequest struct {
	Req models.UserProfileReq
}

type UserProfileResponse struct {
	Res models.UserProfile
	Err error
}

type ChangeUserStatusRequest struct {
	Req models.UsersChangeStatus
}

type ChangeUserStatusResponse struct {
	Res models.ChangeUsers
	Err error
}

func (r CheckAuthResponse) Error() error        { return r.Err }
func (r ChangeRoleResponse) Error() error       { return r.Err }
func (r UsersListResponse) Error() error        { return r.Err }
func (r UserProfileResponse) Error() error      { return r.Err }
func (r ChangeUserStatusResponse) Error() error { return r.Err }
