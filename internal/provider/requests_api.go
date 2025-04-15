// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"io"
	//"bytes"
	//	"encoding/json"
	"net/http" //Package http provides HTTP client and server implementations.
	//"net/http/httputil" //TO DEBUG
	"net/url"
)

//func get_request(params map[string]string, api_url string) ([]byte, error) {
//	// create http object
//	client := &http.Client{}
//
//	//new request
//	//http.NewRequest allow to customize the request vs client.Get()
//	api_request, err := http.NewRequest("GET", api_url, nil)
//	if err != nil {
//		 return nil, err
//	}
//
//	//Set headers
//	api_request.Header.Set("authtoken", SecurdenAuthToken)
//
//	// extract params of request
//	query := api_request.URL.Query()
//
//	// Get every params
//	for key, value := range params {
//		if value != "" {
//			query.Add(key, value)
//		}
//	}
//
//	// encode params in request URL
//	api_request.URL.RawQuery = query.Encode()
//
//	//Send request
//	resp, err := client.Do(api_request)
//	if err != nil {
//		return nil, fmt.Errorf("%v", err)
//	}
//	// Close request
//	defer resp.Body.Close()
//	// Use to read body of response
//	body, err := io.ReadAll(resp.Body)
//
//	return body, err
//}

func (c *securdenClient) getRequest(ctx context.Context, apiURL string, params map[string]string) ([]byte, error) {
	// full URL
	fullURL := fmt.Sprintf("%s/%s", c.ServerURL, apiURL) //  https://bel.securden-vault.com/api/get_password

	// Logs to debug
	//tflog.Debug(ctx, fmt.Sprintf("Making GET request to: %s", fullURL))

	// Set query
	//u, err := url.Parse(fullURL)
	u, _ := url.Parse(fullURL)
	query := u.Query()

	// Get every params
	for key, value := range params {
		if value != "" {
			query.Add(key, value)
			//query.Set(key, value)
		}
	}

	// encode params in request URL
	//req.URL.RawQuery = query.Encode()
	u.RawQuery = query.Encode() //account_id=2000900957088
	updateURL := u.String()

	// HTTP GET request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, updateURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set Header
	req.Header.Set("authtoken", c.AuthToken) // c1**********63

	//DEBUG
	// Dump de la requête complète avec le corps
	//reqDump, err := httputil.DumpRequestOut(req, true)
	//if err == nil {
	//	fmt.Printf("Request send:\n%s\n", string(reqDump))
	//}
	//END DEBUG

	//Send request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	// Close request
	defer resp.Body.Close()

	// Use to read body of response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// DEBUG
	fmt.Printf("Réponse brute: %s\n", string(body))
	// END DEBUG

	// HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	//return body, nil
	return body, err
}

//func post_request(params map[string]interface{}, apiURL string) ([]byte, error) {
//	client := &http.Client{}
//
//	requestBody, err := json.Marshal(params)
//	if err != nil {
//		return nil, fmt.Errorf("failed to serialize request body: %v", err)
//	}
//
//	apiRequest, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
//	if err != nil {
//		return nil, fmt.Errorf("failed to create request: %v", err)
//	}
//
//	apiRequest.Header.Set("Content-Type", "application/json")
//	apiRequest.Header.Set("authtoken", SecurdenAuthToken)
//
//	resp, err := client.Do(apiRequest)
//	if err != nil {
//		return nil, fmt.Errorf("request failed: %v", err)
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return nil, fmt.Errorf("failed to read response body: %v", err)
//	}
//
//	return body, nil
//}
