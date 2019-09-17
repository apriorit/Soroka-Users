package handlers

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Soroka-EDMS/svc/users/pkgs/errors"
)

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
