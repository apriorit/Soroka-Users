//  errors.go
//  https://github.com/apriorit/Soroka-EDMS/svc/users/pkgs/errors
//
//  Created by Ivan Kashuba on 2019.09.03

package errors

import (
	"errors"

	constnt "github.com/Soroka-EDMS/svc/users/pkgs/constants"
)

//All known error definitions for current service
var (
	ErrIncorrectParameter    = errors.New(constnt.IncorrectParameter)
	ErrMissingBody           = errors.New(constnt.MissingBody)
	ErrMissingPayload        = errors.New(constnt.MissingPayload)
	ErrInvalidAuthorization  = errors.New(constnt.InvalidAuthorization)
	ErrProfileNotFound       = errors.New(constnt.ProfileNotFound)
	ErrMalformedRequest      = errors.New(constnt.MalformedRequest)
	ErrMalformedPath         = errors.New(constnt.MalformedPath)
	ErrMissingBodyContent    = errors.New(constnt.MissingBody)
	ErrMalformedBodyContent  = errors.New(constnt.MalformedBody)
	ErrDisabledPartialy      = errors.New(constnt.DisabledPartialy)
	ErrEnabledPartialy       = errors.New(constnt.EnabledPartialy)
	ErrEncodingResponse      = errors.New(constnt.EncodingErorr)
	ErrMissingPath           = errors.New(constnt.MissingPath)
	ErrInvalidRole           = errors.New(constnt.InvalidRole)
	ErrUserDoesNotExist      = errors.New(constnt.UserDoesNotExist)
	ErrInvalidSortParameters = errors.New(constnt.InvalidSortParameters)
)
