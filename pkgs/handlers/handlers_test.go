package handlers

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Soroka-EDMS/svc/users/pkgs/endpoints"
	"github.com/Soroka-EDMS/svc/users/pkgs/errors"
)

//Test DecodeChangeRoleRequest
func TestDecodeChangeRoleRequest_GoodRequest(t *testing.T) {
	//Build a request
	var jsonStr = []byte(`{"role": "admin"}`)
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users/change_role/11,78,53,45", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)

	ctx := context.Background()
	request, err := DecodeChangeRoleRequest(ctx, rawRequest)
	assert.NoError(t, err)

	req, ok := request.(endpoints.ChangeRoleRequest)
	if !ok {
		t.Errorf("Unexpected request")
	}

	assert.Equal(t, "admin", req.Req.Role)
	assert.Equal(t, len(req.Req.Ids), 4)
	assert.Equal(t, req.Req.Ids[0], 11)
	assert.Equal(t, req.Req.Ids[1], 78)
	assert.Equal(t, req.Req.Ids[2], 53)
	assert.Equal(t, req.Req.Ids[3], 45)
}

func TestDecodeChangeRoleRequest_NoIdsInPath(t *testing.T) {
	var jsonStr = []byte(`{"role": "admin"}`)
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users/change_role/", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = DecodeChangeRoleRequest(ctx, rawRequest)
	assert.Error(t, errors.ErrMalformedRequest)
}

func TestDecodeChangeRoleRequest_EmptyBody(t *testing.T) {
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users/change_role/11,78,53,45", nil)
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = DecodeChangeRoleRequest(ctx, rawRequest)
	assert.Error(t, errors.ErrMissingBodyContent)
}

func TestDecodeChangeRoleRequest_MalformedBody(t *testing.T) {
	var jsonStr = []byte(`{"unknown": "admin"}`)
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users/change_role/11,78,53,45", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = DecodeChangeRoleRequest(ctx, rawRequest)
	assert.Error(t, errors.ErrMalformedBodyContent)
}

func TestDecodeChangeRoleRequest_OneIdInPath(t *testing.T) {
	var jsonStr = []byte(`{"role": "admin"}`)
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users/change_role/45", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)

	ctx := context.Background()
	request, err := DecodeChangeRoleRequest(ctx, rawRequest)
	assert.NoError(t, err)

	req, ok := request.(endpoints.ChangeRoleRequest)
	if !ok {
		t.Errorf("Unexpected request")
	}

	assert.Equal(t, "admin", req.Req.Role)
	assert.Equal(t, len(req.Req.Ids), 1)
	assert.Equal(t, req.Req.Ids[0], 45)
}

func TestDecodeGetUserListRequest_GoodQuery(t *testing.T) {
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users?offset=10&limit=100&sort=id&order=up", nil)
	assert.NoError(t, err)

	ctx := context.Background()
	request, err := DecodeGetUserListRequest(ctx, rawRequest)
	assert.NoError(t, err)

	req, ok := request.(endpoints.UsersListRequest)

	if !ok {
		t.Errorf("Unexpected request")
	}

	assert.Equal(t, 10, req.Req.Offset)
	assert.Equal(t, 100, req.Req.Limit)
	assert.Equal(t, "id", req.Req.Sort)
	assert.Equal(t, "up", req.Req.Order)
}

func TestDecodeGetUserListRequest_UnexpectedQuery(t *testing.T) {
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users?a=xyz&b[c]=456", nil)
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = DecodeGetUserListRequest(ctx, rawRequest)
	assert.NoError(t, err)
}

func TestDecodeGetUserProfileRequest_GoodQuery(t *testing.T) {
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users?id=57&email=user@example.com", nil)
	assert.NoError(t, err)

	ctx := context.Background()
	request, err := DecodeGetUserProfileRequest(ctx, rawRequest)
	assert.NoError(t, err)

	req, ok := request.(endpoints.UserProfileRequest)

	if !ok {
		t.Errorf("Unexpected request")
	}

	assert.Equal(t, 57, req.Req.Id)
	assert.Equal(t, "user@example.com", req.Req.Email)
}

func TestDecodeGetUserProfileRequest_NoEmail(t *testing.T) {
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users?id=57", nil)
	assert.NoError(t, err)

	ctx := context.Background()
	request, err := DecodeGetUserProfileRequest(ctx, rawRequest)
	assert.NoError(t, err)

	_, ok := request.(endpoints.UserProfileRequest)
	assert.True(t, ok)
}

func TestDecodeGetUserProfileRequest_NoId(t *testing.T) {
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users?email=user@example.com", nil)
	assert.NoError(t, err)

	ctx := context.Background()
	request, err := DecodeGetUserProfileRequest(ctx, rawRequest)
	assert.NoError(t, err)

	_, ok := request.(endpoints.UserProfileRequest)
	assert.True(t, ok)
}

func TestDecodeGetUserProfileRequest_NoParamsInQuery(t *testing.T) {
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users", nil)
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = DecodeGetUserProfileRequest(ctx, rawRequest)
	assert.Error(t, err)
}

func TestDecodeDisableUsersRequest_GoodQuery(t *testing.T) {
	testData := []int{11, 57, 44, 89, 23, 111}

	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users/disable/11,57,44,89,23,111", nil)
	assert.NoError(t, err)

	ctx := context.Background()
	request, err := DecodeChangeUsersRequest(ctx, rawRequest)
	assert.NoError(t, err)

	req, ok := request.(endpoints.ChangeUserStatusRequest)
	if !ok {
		t.Errorf("Unexpected request")
	}

	assert.NotEmpty(t, req.Req.Ids)
	assert.Equal(t, testData, req.Req.Ids)
}

func TestDecodeDisableUsersRequest_IvalidTypeInQuery(t *testing.T) {
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users/disable/11,57,invalid,89,23,111", nil)
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = DecodeChangeUsersRequest(ctx, rawRequest)
	assert.Error(t, err)
}
