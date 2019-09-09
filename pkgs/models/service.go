//  service.go
//  https://github.com/apriorit/Soroka-EDMS/svc/users/pkgs/models
//
//  Created by Ivan Kashuba on 2019.09.03
//  Describe service models

package models

import (
	"context"
	"net/http"
)

type UsersRequest struct {
	Url http.Request
}

type UsersService interface {
	ChangeRole(cntx context.Context, request ChangeRoleRequest) (res ChangeUsersResponse, err error)
	GetUserList(cntx context.Context, request UsersListRequest) (res UsersListResponse, err error)
	GetUserProfile(cntx context.Context, request UserProfileRequest) (res UserProfile, err error)
	DisableUsers(cntx context.Context, request UsersChangeStatusRequest) (res ChangeUsersResponse, err error)
	EnableUsers(cntx context.Context, request UsersChangeStatusRequest) (res ChangeUsersResponse, err error)
}

type Role struct {
	Name string `json:"name"`
	Mask int64  `json:"mask"`
}

type UserInfo struct {
	User_name     string `json:"user_name"`
	User_id       int    `json:"user_id"`
	Role          string `json:"role"`
	Location      string `json:"Dnipro"`
	Email         string `json:"email"`
	Creation_date int64  `json:"creation_date"`
	Status        bool   `json:"status"`
}

type PaginationInfo struct {
	Issued int `json:"issued"`
	Left   int `json:"left"`
}

type ChangeRoleRequest struct {
	Ids  []int  `json:"ids"`
	Role string `json:"role"`
}

type UsersListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Sort   string `json:"sort"`
	Order  string `json:"order"`
}

type UsersListResponse struct {
	Users      []UserInfo     `json:"users"`
	Pagination PaginationInfo `json:"pagination"`
}

type UserProfileRequest struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type UserProfile struct {
	First_name    string   `json:"first_name"`
	Last_name     string   `json:"last_name"`
	Email         string   `json:"email"`
	Phone         string   `json:"phone"`
	Location      string   `json:"location"`
	Position      string   `json:"position"`
	Status        bool     `json:"status"`
	Creation_date int64    `json:"creation_date"`
	Role          UserRole `json:"role"`
}

type UserProfileResponse struct {
	First_name    string   `json:"first_name"`
	Last_name     string   `json:"last_name"`
	Email         string   `json:"email"`
	Phone         string   `json:"phone"`
	Location      string   `json:"location"`
	Position      string   `json:"position"`
	Status        bool     `json:"status"`
	Creation_date int64    `json:"creation_date"`
	Role          UserRole `json:"role"`
}

type UsersChangeStatusRequest struct {
	Ids []int `json:"id"`
}

type ChangeUsersResponse struct {
	Changed   []int  `json:"changed"`
	Not_found []int  `json:"not_found"`
	Message   string `json:"message"`
}
