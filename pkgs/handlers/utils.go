package handlers

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joncalhoun/qson"

	"github.com/Soroka-EDMS/svc/users/pkgs/config"
	cf "github.com/Soroka-EDMS/svc/users/pkgs/config"
	c "github.com/Soroka-EDMS/svc/users/pkgs/constants"
	e "github.com/Soroka-EDMS/svc/users/pkgs/errors"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
	m "github.com/Soroka-EDMS/svc/users/pkgs/models"
)

//GetIds returns slice with ids that enclosed in request path
func GetIds(path string) ([]int, error) {
	if path == "" {
		return nil, e.ErrMissingPath
	}

	parts := strings.Split(path, "/")
	splitedLastPart := strings.Split(parts[len(parts)-1], ",")

	if len(splitedLastPart) == 0 {
		return nil, e.ErrMalformedPath
	}

	ids := make([]int, len(splitedLastPart))

	for count := 0; count < len(ids); count++ {
		id, err := strconv.Atoi(splitedLastPart[count])

		if err != nil {
			return nil, e.ErrMalformedPath
		}

		ids[count] = id
	}

	return ids, nil
}

//QueryToUsersListRequest converts query string from request to models.UsersList that describes pagination parameters
func QueryToUsersListRequest(query string, req *m.UsersList) error {

	err := qson.Unmarshal(&req, query)
	if err != nil {
		return err
	}

	return nil
}

//QueryToUserProfileRequest converts query string from request to models.UserProfileReq
func QueryToUserProfileRequest(query string, req *m.UserProfileReq) error {
	err := qson.Unmarshal(&req, query)
	if err != nil {
		return err
	}

	return nil
}

//GetTokensPair returns a pair of tokens (access and refresh) or error if one of token is missing
func GetTokensPair(req *http.Request) (res m.TokensPair, err error) {
	//Get access token string
	authHeader := req.Header.Get("Authorization")
	config.GetLogger().Logger.Log("method", "GetTokensPair", "authHeader", authHeader)
	if authHeader == "" {
		return res, e.ErrMissingAuthorizationHeader
	}

	if !strings.HasPrefix(authHeader, c.BearerSchema) {
		return res, e.ErrBearerSchemaRequired
	}
	//Get refresh token from cookie
	cookie, err := req.Cookie("refresh_token")
	if err != nil {
		return res, e.ErrRefreshTokenRequired
	}

	res.AccessToken = authHeader[len(c.BearerSchema):]
	res.RefreshToken = cookie.Value

	cf.GetLogger().Logger.Log("a", res.AccessToken, "r", res.RefreshToken)

	return res, nil
}

//CheckJWT calls  check_token endpoint of sessions service in order to update or confirm token validness
func CheckJWT(tp m.TokensPair) (res m.AccessToken, err error) {
	//Prepare request with access token in body and refresh token in cookie
	rawTokenString := `{"access_token": "%s"}`
	req, err := http.NewRequest("POST", c.URIOnCheckTokenProfile, bytes.NewBuffer([]byte(fmt.Sprintf(rawTokenString, tp.AccessToken))))
	if err != nil {
		config.GetLogger().Logger.Log("method", "CheckJWT", "1", "1")
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
	/*var netTransport = &http.Transport{
		MaxIdleConns: 10,
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	// tune net client
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}*/
	netClient, err := MakeHTTPClient()

	if err != nil {
		return res, err
	}

	//Do request
	resp, err := netClient.Do(req)

	if err != nil {
		config.GetLogger().Logger.Log("method", "CheckJWT", "2", "2")
		return res, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusBadRequest:
		cf.GetLogger().Logger.Log("err", "2")
		return res, e.ErrRequestToSessionFailed
	case http.StatusUnauthorized:
		return res, e.ErrNonAuthorized
	default:
		cf.GetLogger().Logger.Log("code", resp.StatusCode)
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

//ParseClaims gets claims and checks client permissions
func ParseClaims(method, token string) error {
	//Get raw token
	rawToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		//return session secret that was used in signing process
		return []byte(cf.GetSecretString()), nil
	})

	if err != nil {
		return err
	}

	//Check token validness
	var claims jwt.MapClaims
	claims, ok := rawToken.Claims.(jwt.MapClaims)
	if !ok {
		return e.ErrInvalidClaims
	}
	mask, ok := claims["mask"].(float64)
	if !ok {
		return e.ErrInvalidMaskType
	}
	//Check priveldge
	if !CheckPriveledge(int64(mask), method) {
		return e.ErrInvalidAuthorization
	}

	return nil
}

func CheckPriveledge(mask int64, method string) bool {
	if method == "" || mask == 0 {
		return false
	}

	privMap := config.GetPriveledges().Priveledges
	bitValue := mask >> (uint(privMap[method] - 1)) & 1
	return bitValue == 1
}

func MakeHTTPClient() (*http.Client, error) {
	//Get public key
	pkey := cf.GetPublicKey()

	if len(pkey) == 0 {
		return nil, e.ErrPublicKeyIsMissing
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(pkey)

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

//CheckAuthorization parse pair of tokens from request, checks their validness by calling sessions service and return an access token.
func CheckAuthorization(rawRequest *http.Request, method string) error {
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
	if at, err = CheckJWT(tPair); err != nil {
		return err
	}

	//ParseClaims
	if err = ParseClaims(method, at.AccessToken); err != nil {
		return err
	}

	return nil
}
