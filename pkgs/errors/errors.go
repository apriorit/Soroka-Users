//  errors.go
//  https://github.com/apriorit/Soroka-EDMS/svc/users/pkgs/errors
//
//  Created by Ivan Kashuba on 2019.09.03

package errors

import (
	"errors"

	c "github.com/Soroka-EDMS/svc/users/pkgs/constants"
)

//All known error definitions for current service
var (
	ErrIncorrectParameter         = errors.New(c.IncorrectParameter)
	ErrMissingAuthorizationHeader = errors.New(c.MissingAuthorizationHeader)
	ErrBearerSchemaRequired       = errors.New(c.BearerSchemaRequired)
	ErrBasicSchemaRequired        = errors.New(c.BasicSchemaRequired)
	ErrRefreshTokenRequired       = errors.New(c.RefreshTokenRequired)
	ErrMissingBody                = errors.New(c.MissingBody)
	ErrMissingPayload             = errors.New(c.MissingPayload)
	ErrInvalidAuthorization       = errors.New(c.InvalidAuthorization)
	ErrNonAuthorized              = errors.New(c.RequiredAuthorized)
	ErrRequestToSessionFailed     = errors.New(c.RequestToSessionFailed)
	ErrProfileNotFound            = errors.New(c.ProfileNotFound)
	ErrMalformedRequest           = errors.New(c.MalformedRequest)
	ErrMalformedPath              = errors.New(c.MalformedPath)
	ErrMissingBodyContent         = errors.New(c.MissingBody)
	ErrInvalidMaskType            = errors.New(c.InvalidMaskType)
	ErrInvalidClaims              = errors.New(c.InvalidClaims)
	ErrMalformedBodyContent       = errors.New(c.MalformedBody)
	ErrDisabledPartialy           = errors.New(c.DisabledPartialy)
	ErrEnabledPartialy            = errors.New(c.EnabledPartialy)
	ErrEncodingResponse           = errors.New(c.EncodingErorr)
	ErrMissingPath                = errors.New(c.MissingPath)
	ErrInvalidRole                = errors.New(c.InvalidRole)
	ErrUserDoesNotExist           = errors.New(c.UserDoesNotExist)
	ErrInvalidSortParameters      = errors.New(c.InvalidSortParameters)
	ErrPublicKeyIsMissing         = errors.New(c.PublicKeyIsMissing)
)
