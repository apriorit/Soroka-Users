package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	"github.com/Soroka-EDMS/svc/users/pkgs/models"
	m "github.com/Soroka-EDMS/svc/users/pkgs/models"
)

type UsersEndpoints struct {
	ChangeRoleEndpoint     endpoint.Endpoint
	GetUserListEndpoint    endpoint.Endpoint
	GetUserProfileEndpoint endpoint.Endpoint
	DisableUsersEndpoint   endpoint.Endpoint
	EnableUsersEndpoint    endpoint.Endpoint
}

//Build creates endpoints with authorization by priveledges
func Build(logger log.Logger, service models.UsersService) (endpoints UsersEndpoints) {
	endpoints.ChangeRoleEndpoint = BuildChangeRoleEndpoint(service)
	endpoints.ChangeRoleEndpoint = JwtAuthorization("changeRole")(endpoints.ChangeRoleEndpoint)

	endpoints.GetUserListEndpoint = BuildGetUserListEndpoint(service)
	endpoints.GetUserListEndpoint = JwtAuthorization("userList")(endpoints.GetUserListEndpoint)

	endpoints.GetUserProfileEndpoint = BuildGetUserProfileEndpoint(service)
	endpoints.GetUserProfileEndpoint = JwtAuthorization("userProfile")(endpoints.GetUserProfileEndpoint)

	endpoints.DisableUsersEndpoint = BuildDisableUsersEndpoint(service)
	endpoints.DisableUsersEndpoint = JwtAuthorization("disableUser")(endpoints.DisableUsersEndpoint)

	endpoints.EnableUsersEndpoint = BuildEnableUsersEndpoint(service)
	endpoints.EnableUsersEndpoint = JwtAuthorization("enableUser")(endpoints.EnableUsersEndpoint)

	return endpoints
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
		res, err := svc.DisableUsers(ctx, req.Req)
		return ChangeUserStatusResponse{Res: res, Err: err}, nil
	}
}

type ChangeRoleRequest struct {
	Req m.ChangeRoleRequest
}

type ChangeRoleResponse struct {
	Res m.ChangeUsersResponse
	Err error
}

type UsersListRequest struct {
	Req m.UsersListRequest
}

type UsersListResponse struct {
	Res m.UsersListResponse
	Err error
}

type UserProfileRequest struct {
	Req m.UserProfileRequest
}

type UserProfileResponse struct {
	Res m.UserProfile
	Err error
}

type ChangeUserStatusRequest struct {
	Req m.UsersChangeStatusRequest
}

type ChangeUserStatusResponse struct {
	Res m.ChangeUsersResponse
	Err error
}

func (r ChangeRoleResponse) Error() error       { return r.Err }
func (r UsersListResponse) Error() error        { return r.Err }
func (r UserProfileResponse) Error() error      { return r.Err }
func (r ChangeUserStatusResponse) Error() error { return r.Err }
