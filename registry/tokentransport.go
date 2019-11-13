package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)



type TokenTransport struct {
	Transport http.RoundTripper
}

type authToken struct {
	ExpiresIn   int   `json:"expires_in"`
	Token string `json:"token"`
}

type ResponseError struct {
	Errors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Detail  []struct {
			Type   string `json:"Type"`
			Class  string `json:"Class"`
			Name   string `json:"Name"`
			Action string `json:"Action"`
		} `json:"detail"`
	} `json:"errors"`
}

func NewTokenTransport(transport http.RoundTripper) *TokenTransport {
	return &TokenTransport{Transport:transport}
}

func (t *TokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	if authService := isTokenDemand(resp); authService != nil {
		defer resp.Body.Close()
		resp, err = t.authAndRetry(authService, req)
	}
	return resp, err
}


func (t *TokenTransport) authAndRetry(authService *authService, req *http.Request) (*http.Response, error) {
	token, authResp, err := t.auth(authService)
	if err != nil {
		return authResp, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := t.Transport.RoundTrip(req)
	return resp, err
}

func (t *TokenTransport) auth(authService *authService) (string, *http.Response, error) {
	authReq, err := authService.Request()
	if err != nil {
		return "", nil, err
	}

	client := http.Client{
		Transport: t.Transport,
	}

	response, err := client.Do(authReq)
	if err != nil {
		return "", nil, err
	}

	if response.StatusCode != http.StatusOK {
		return "", response, err
	}
	defer response.Body.Close()

	var authToken authToken
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&authToken)
	if err != nil {
		return "", nil, err
	}

	return authToken.Token, nil, nil
}


type authService struct {
	Realm   string
	Service string
	Scope   string
}

func (authService *authService) Request() (*http.Request, error) {
	url, err := url.Parse(authService.Realm)
	if err != nil {
		return nil, err
	}

	q := url.Query()
	q.Set("service", authService.Service)
	if authService.Scope != "" {
		q.Set("scope", authService.Scope)
	}
	url.RawQuery = q.Encode()

	request, err := http.NewRequest("GET", url.String(), nil)

	return request, err
}

func isTokenDemand(resp *http.Response) *authService {
	if resp == nil {
		return nil
	}
	if resp.StatusCode != http.StatusUnauthorized {
		return nil
	}
	return parseOauthHeader(resp)
}

func parseOauthHeader(resp *http.Response) *authService {
	if len(resp.Header["Www-Authenticate"]) > 0 {
		result := make(map[string]string)
		wantedHeaders := []string{"realm", "service", "scope"}
		authHeaderValueSlice := strings.Split(resp.Header["Www-Authenticate"][0], ",")
		for _, r := range authHeaderValueSlice {
			for _, w := range wantedHeaders {
				if strings.Contains(r, w) {
					result[w] = strings.Split(r, `"`)[1]
				}
			}
		}
		return &authService{
			Realm:   result["realm"],
			Service: result["service"],
			Scope:   result["scope"],
		}
	}

	return nil
}