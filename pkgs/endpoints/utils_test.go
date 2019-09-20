package endpoints

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	er "github.com/Soroka-EDMS/svc/users/pkgs/errors"
)

//GetTokensPair_BearerAndCookiePresent test
func TestGetTokensPair_BearerAndCookiePresent(t *testing.T) {
	cookie := &http.Cookie{
		Name:  "refresh_token",
		Value: "p672b44A-0e5b-4c2d-u12b-c7be8c7hV3y9",
	}

	rawRequest, err := http.NewRequest("POST", "https://edms.com/session/check_token", nil)
	assert.NoError(t, err)

	rawRequest.Header.Add("Authorization", "Bearer d352b45d-0e5b-4c2d-a10b-c7be8c7cd3ff")
	rawRequest.AddCookie(cookie)

	resp, err := GetTokensPair(rawRequest)
	assert.NoError(t, err)
	assert.Equal(t, resp.AccessToken, "d352b45d-0e5b-4c2d-a10b-c7be8c7cd3ff")
	assert.Equal(t, resp.RefreshToken, "p672b44A-0e5b-4c2d-u12b-c7be8c7hV3y9")
}

func TestGetTokensPair_AuthorizationDoesNotExist(t *testing.T) {
	rawRequest, err := http.NewRequest("POST", "https://edms.com/session/check_token", nil)
	assert.NoError(t, err)

	_, err = GetTokensPair(rawRequest)
	assert.Error(t, err)
	assert.Equal(t, err, er.ErrMissingAuthorizationHeader)
}

func TestGetTokensPair_BearerSchemaDoesNotExist(t *testing.T) {
	rawRequest, err := http.NewRequest("POST", "https://edms.com/session/check_token", nil)
	assert.NoError(t, err)

	rawRequest.Header.Add("Authorization", "Basic d352b45d-0e5b-4c2d-a10b-c7be8c7cd3ff")

	_, err = GetTokensPair(rawRequest)
	assert.Error(t, err)
	assert.Error(t, er.ErrBearerSchemaRequired)
}

func TestCheckPriveledge_ExpectedTrueForAll(t *testing.T) {
	var tests = []struct {
		role   int64
		method string
	}{
		{2048, "changeRole"},
		{16384, "userList"},
		{8192, "userProfile"},
		{4096, "changeUserStatus"},
	}

	for _, test := range tests {
		if res := CheckPriveledge(test.role, test.method); !res {
			t.Errorf("Mask %d has no priveledge for %s", test.role, test.method)
		}
	}
}

func TestCheckPriveledge_EmptyMask(t *testing.T) {
	var tests = []struct {
		role   int64
		method string
	}{
		{2048, ""},
		{8, ""},
		{64, ""},
		{4096, ""},
		{512, ""},
	}

	for _, test := range tests {
		if res := CheckPriveledge(test.role, test.method); res {
			t.Errorf("Expected false, find %v", res)
		}
	}
}

func TestCheckPriveledge_HasNoPriveledge(t *testing.T) {
	var tests = []struct {
		role   int64
		method string
	}{
		{8, "userList"},
		{64, "userProfile"},
		{512, "enableUser"},
	}

	for _, test := range tests {
		if res := CheckPriveledge(test.role, test.method); res {
			t.Errorf("Expected false, find %v", res)
		}
	}
}
