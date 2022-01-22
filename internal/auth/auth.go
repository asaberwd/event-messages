package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// HTTPClient ...
type HTTPClient interface {
	Request(url string, method string, header map[string]string, body []byte, params map[string]string) (statusCode int, resBody []byte, err error)
}

// Provider service API
type Provider struct {
	HTTPClient   HTTPClient
	AuthAudience string
}

func NewProvider(HTTPClient HTTPClient, autAudience string) *Provider {
	return &Provider{
		HTTPClient:   HTTPClient,
		AuthAudience: autAudience,
	}
}

// JWTAuth implements the authorization logic for the JWT security scheme.
func (a *Provider) JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		header := make(map[string]string)
		header["Authorization"] = fmt.Sprintf("Bearer %s", token)
		URL := fmt.Sprintf("%s%s", a.AuthAudience, "authping")

		statusCode, body, err := a.HTTPClient.Request(URL, "GET", header, nil, nil)
		if err != nil {
			log.Errorf(fmt.Sprintf("auth error: %+v for GET request to authenticate ", err))
			return errors.New("ErrInvalidToken")
		}

		var result Result
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Errorf(fmt.Sprintf("Unmarshal error: %+v", err))
			return err
		}

		if statusCode != 200 {
			log.Warn(fmt.Sprintf("auth request return with status code %v , url %s", statusCode, URL))
			return echo.ErrUnauthorized
		}
		c.Set("creator", result.Sub )
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

// Result ...
type Result struct {
	Sub   string `json:"sub"`
	Iss   string `json:"iss"`
	Aud   string `json:"aud"`
	Iat   int    `json:"iat"`
	Exp   int    `json:"exp"`
	Scope string `json:"scope"`
	Gty   string `json:"gty"`
}
