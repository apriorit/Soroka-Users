package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	m "github.com/Soroka-EDMS/svc/users/pkgs/models"
)

type UsersEndpoints struct {
	CheckAuthEndpoint      endpoint.Endpoint
	ChangeRoleEndpoint     endpoint.Endpoint
	GetUserListEndpoint    endpoint.Endpoint
	GetUserProfileEndpoint endpoint.Endpoint
	DisableUsersEndpoint   endpoint.Endpoint
	EnableUsersEndpoint    endpoint.Endpoint
}

//Build creates endpoints with authorization by priveledges
func Build(logger log.Logger, service m.UsersService) (endpoints UsersEndpoints) {
	return UsersEndpoints{
		CheckAuthEndpoint:      BuildCheckAuthEndpoint(service),
		ChangeRoleEndpoint:     BuildChangeRoleEndpoint(service),
		GetUserListEndpoint:    BuildGetUserListEndpoint(service),
		GetUserProfileEndpoint: BuildGetUserProfileEndpoint(service),
		DisableUsersEndpoint:   BuildDisableUsersEndpoint(service),
		EnableUsersEndpoint:    BuildEnableUsersEndpoint(service),
	}
}

//BuildChangeRoleEndpoint returns an endpoint that call CheckAuth handler.
func BuildCheckAuthEndpoint(svc m.UsersService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CheckAuthRequest)
		err := svc.CheckAuth(ctx, req.Req)
		return CheckAuthResponse{Err: err}, nil
	}
}

// BuildChangeRoleEndpoint returns an endpoint that call ChangeRole handler.
func BuildChangeRoleEndpoint(svc m.UsersService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeRoleRequest)
		res, err := svc.ChangeRole(ctx, req.Req)
		return ChangeRoleResponse{Res: res, Err: err}, nil
	}
}

// BuildGetUserListEndpoint returns an endpoint that call GetUserList handler.
func BuildGetUserListEndpoint(svc m.UsersService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UsersListRequest)
		res, err := svc.GetUserList(ctx, req.Req)
		return UsersListResponse{Res: res, Err: err}, nil
	}
}

// BuildGetUserProfileEndpoint returns an endpoint that call GetUserProfile handler.
func BuildGetUserProfileEndpoint(svc m.UsersService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UserProfileRequest)
		res, err := svc.GetUserProfile(ctx, req.Req)
		return UserProfileResponse{Res: res, Err: err}, nil
	}
}

// BuildDisableUsersEndpoint returns an endpoint that call DisableUsers handler.
func BuildDisableUsersEndpoint(svc m.UsersService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeUserStatusRequest)
		res, err := svc.DisableUsers(ctx, req.Req)
		return ChangeUserStatusResponse{Res: res, Err: err}, nil
	}
}

// BuildEnableUsersEndpoint returns an endpoint that call EnableUsers handler.
func BuildEnableUsersEndpoint(svc m.UsersService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeUserStatusRequest)
		res, err := svc.EnableUsers(ctx, req.Req)
		return ChangeUserStatusResponse{Res: res, Err: err}, nil
	}
}

type CheckAuthRequest struct {
	Req m.UserCredentials
}

type CheckAuthResponse struct {
	Err error
}

type ChangeRoleRequest struct {
	Req m.ChangeRole
}

type ChangeRoleResponse struct {
	Res m.ChangeUsers
	Err error
}

type UsersListRequest struct {
	Req m.UsersList
}

type UsersListResponse struct {
	Res m.UsersListResp
	Err error
}

type UserProfileRequest struct {
	Req m.UserProfileReq
}

type UserProfileResponse struct {
	Res m.UserProfileResp
	Err error
}

type ChangeUserStatusRequest struct {
	Req m.UsersChangeStatus
}

type ChangeUserStatusResponse struct {
	Res m.ChangeUsers
	Err error
}

func (r CheckAuthResponse) Error() error        { return r.Err }
func (r ChangeRoleResponse) Error() error       { return r.Err }
func (r UsersListResponse) Error() error        { return r.Err }
func (r UserProfileResponse) Error() error      { return r.Err }
func (r ChangeUserStatusResponse) Error() error { return r.Err }
