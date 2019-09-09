package endpoints

import (
	"context"

	"github.com/Soroka-EDMS/svc/users/pkgs/config"
	jwtg "github.com/dgrijalva/jwt-go"
	jwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"

	er "github.com/Soroka-EDMS/svc/users/pkgs/errors"
)

type PayloadAuth struct {
	Mask int64 `json:"mask"`
	jwtg.StandardClaims
}

func CheckPriveledge(mask int64, method string) bool {
	if method == "" || mask == 0 {
		return false
	}

	privMap := config.GetPriveledges().Priveledges
	bitValue := mask >> (uint(privMap[method] - 1)) & 1
	return bitValue == 1
}

func JwtAuthorization(method string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			// retrieve auth claims
			claims, err := getPayload(ctx)
			if err != nil {
				return nil, err
			}

			priveledge := CheckPriveledge(claims.Mask, method)
			if !priveledge {
				return nil, er.ErrInvalidAuthorization
			}

			return next(ctx, request)
		}
	}
}

// getPayload takes payload from context and returns Payload.
func getPayload(ctx context.Context) (*PayloadAuth, error) {
	payloadClaims, ok := ctx.Value(jwt.JWTClaimsContextKey).(*PayloadAuth)
	if !ok {
		return nil, er.ErrMissingPayload
	}

	return payloadClaims, nil
}
