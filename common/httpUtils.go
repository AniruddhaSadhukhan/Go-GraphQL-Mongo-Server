package common

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/logger"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")

	// Set HSTS header is HTTPS is enabled
	if config.Store.HTTPSCert.HTTPSEnabled {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}

	if config.Store.CORSAllowOrigins != "" {
		w.Header().Set("Access-Control-Allow-Origin", config.Store.CORSAllowOrigins)
	}
	w.WriteHeader(code)

	if code != http.StatusNoContent {
		_, err := w.Write(response)

		if err != nil {
			logger.Log.Errorf("Error in writing response %+v", err)
		}
	}

}

func RespondWithUnauthorized(w http.ResponseWriter) {
	RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
}

type HTTPClient struct {
	url        string
	httpClient *http.Client
	headers    http.Header
}

func GetHTTPClient(useProxy bool) *http.Client {

	tr := &http.Transport{}
	if config.Store.InsecureSkipVerify {
		// #nosec G402
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if useProxy {
		tr.Proxy = http.ProxyFromEnvironment
	}

	return &http.Client{
		Transport: tr,
	}
}

func NewHTTPClient(url string, token string, tlsSkipVerify bool) *HTTPClient {

	client := &HTTPClient{
		url:        url,
		httpClient: &http.Client{},
		headers: http.Header{
			"Content-Type": []string{"application/json"},
			"Accept":       []string{"application/json"},
		},
	}

	//nolint:gosec
	if tlsSkipVerify {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client.httpClient.Transport = customTransport
	}

	if token != "" {
		client.headers.Add("Authorization", token)
	}

	return client
}

func (client *HTTPClient) ExecuteGraphQl(query string, variables map[string]any, response any) error {

	var reqBody bytes.Buffer
	reqBodyObj := struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}{
		Query:     query,
		Variables: variables,
	}
	if err := json.NewEncoder(&reqBody).Encode(reqBodyObj); err != nil {
		return fmt.Errorf("error encoding request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, client.url, &reqBody)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header = client.headers

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error executing request: %v", resp.Status)
	}

	defer resp.Body.Close()

	gqlResponse := struct {
		Data   any
		Errors []struct {
			Message string
		}
	}{
		Data: response,
	}

	if err := json.NewDecoder(resp.Body).Decode(&gqlResponse); err != nil {
		x := new(bytes.Buffer)
		_, resError := x.ReadFrom(resp.Body)
		logger.Log.Errorf("error decoding response body: %v", resError)
		return fmt.Errorf("error decoding response body : %v : %v", err, x.String())
	}

	if len(gqlResponse.Errors) > 0 {
		return fmt.Errorf("graphql error : %v", gqlResponse.Errors[0].Message)
	}
	return nil

}
