package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"web_exporter/internal/models"
	"web_exporter/pkg/logger"
)

var LOGINHandler = "/login"
var BASEHandler = "/panel/api/inbounds"
var ONLINESHandler = "/onlines"
var INBOUNDSHandler = "/list"

type SessionExpiredError struct{}

func (e SessionExpiredError) Error() string {
	return "Your session has expired, please log in again"
}

func Auth(callUrl, username, password string) ([]*http.Cookie, error) {
	logger.Logf(logger.WebLogPrefix, "authenticating")

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	credentials := strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", callUrl+LOGINHandler, credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result models.API3XUIOnlinesResp
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to auth, error message: %s", result.Msg)
	}

	logger.Logf(logger.WebLogPrefix, "successifuly authenticated")
	return resp.Cookies(), nil
}

func sendRequest(method, uri string, cookies []*http.Cookie, result any) error {
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for i := range cookies {
		req.AddCookie(cookies[i])
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	logger.Logf(logger.WebLogPrefix, "REQUEST: (%s)%s RESPONSE (truncated): %s", method, uri, string(body[:100]))

	err = json.Unmarshal(body, result)

	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

func getInbounds(callUrl string, cookies []*http.Cookie) ([]models.API3XUIInbound, error) {
	logger.Logf(logger.WebLogPrefix, "getting inbounds")

	var response models.API3XUIInboundsResp

	err := sendRequest("GET", callUrl+BASEHandler+INBOUNDSHandler, cookies, &response)

	if err != nil {
		return nil, err
	}

	if !response.Success {
		if response.Msg == "Your session has expired, please log in again" {
			return nil, SessionExpiredError{}
		}

		return nil, fmt.Errorf("%s", response.Msg)
	}

	return response.Inbounds, nil
}

func GetClients(callUrl string, cookies []*http.Cookie) ([]models.API3XUIClientStats, error) {
	logger.Logf(logger.WebLogPrefix, "getting clients")

	var clients []models.API3XUIClientStats

	inbounds, err := getInbounds(callUrl, cookies)
	if err != nil {
		return nil, fmt.Errorf("failed to get inbounds: %w", err)
	}

	for _, inbound := range inbounds {
		clients = append(clients, inbound.ClientStats...)
	}

	return clients, nil
}

func GetOnlines(callUrl string, cookies []*http.Cookie) ([]string, error) {
	logger.Logf(logger.WebLogPrefix, "getting onlines")
	var response models.API3XUIOnlinesResp

	err := sendRequest("POST", callUrl+BASEHandler+ONLINESHandler, cookies, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		if response.Msg == "Your session has expired, please log in again" {
			return nil, SessionExpiredError{}
		}

		return nil, fmt.Errorf("%s", response.Msg)
	}

	return response.Emails, nil
}
