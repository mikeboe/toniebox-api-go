package toniebox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// requestHandler handles all HTTP requests to the Toniebox API
type requestHandler struct {
	client   *http.Client
	jwtToken *JWTToken
}

// newRequestHandler creates a new request handler with default settings
func newRequestHandler() *requestHandler {
	return &requestHandler{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// newRequestHandlerWithProxy creates a new request handler with proxy settings
func newRequestHandlerWithProxy(proxyURL string) (*requestHandler, error) {
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %w", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	return &requestHandler{
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}, nil
}

// login authenticates the user and stores the JWT token
func (rh *requestHandler) login(loginData *Login) error {
	data := url.Values{}
	data.Set("grant_type", grantTypePassword)
	data.Set("client_id", clientID)
	data.Set("scope", scopeOpenID)
	data.Set("username", loginData.Email)
	data.Set("password", loginData.Password)

	req, err := http.NewRequest("POST", openIDConnect, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", contentTypeForm)

	resp, err := rh.client.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var token JWTToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return fmt.Errorf("failed to decode token: %w", err)
	}

	rh.jwtToken = &token
	return nil
}

// getMe retrieves personal information about the authenticated user
func (rh *requestHandler) getMe() (*Me, error) {
	var result Me
	if err := rh.executeGetRequest(me, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// getHouseholds retrieves all households the user belongs to
func (rh *requestHandler) getHouseholds() ([]Household, error) {
	var result []Household
	if err := rh.executeGetRequest(households, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// getCreativeTonies retrieves all Creative-Tonies in a household
func (rh *requestHandler) getCreativeTonies(household *Household) ([]CreativeTonie, error) {
	url := fmt.Sprintf(creativeTonies, household.ID)
	var result []CreativeTonie
	if err := rh.executeGetRequest(url, &result); err != nil {
		return nil, err
	}

	// Set household reference and request handler for each tonie
	for i := range result {
		result[i].household = household
		result[i].requestHandler = rh
	}

	return result, nil
}

// refreshTonie retrieves the latest state of a Creative-Tonie
func (rh *requestHandler) refreshTonie(tonie *CreativeTonie) (*CreativeTonie, error) {
	url := fmt.Sprintf(creativeTonie, tonie.household.ID, tonie.ID)
	var result CreativeTonie
	if err := rh.executeGetRequest(url, &result); err != nil {
		return nil, err
	}

	result.household = tonie.household
	result.requestHandler = rh
	return &result, nil
}

// commitTonie saves changes to a Creative-Tonie
func (rh *requestHandler) commitTonie(tonie *CreativeTonie) error {
	url := fmt.Sprintf(creativeTonie, tonie.household.ID, tonie.ID)

	body, err := json.Marshal(tonie)
	if err != nil {
		return fmt.Errorf("failed to marshal tonie: %w", err)
	}

	return rh.executePatchRequest(url, body)
}

// uploadFile uploads a file to a Creative-Tonie
func (rh *requestHandler) uploadFile(tonie *CreativeTonie, filePath, title string) error {
	// Step 1: Request upload credentials from Toniebox API
	emptyBody := []byte(`{"headers":{}}`)

	req, err := http.NewRequest("POST", fileUpload, bytes.NewReader(emptyBody))
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Authorization", "Bearer "+rh.jwtToken.AccessToken)

	resp, err := rh.client.Do(req)
	if err != nil {
		return fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var amazonBean AmazonBean
	if err := json.NewDecoder(resp.Body).Decode(&amazonBean); err != nil {
		return fmt.Errorf("failed to decode amazon response: %w", err)
	}

	// Step 2: Upload file to Amazon S3
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	fields := amazonBean.Request.Fields
	if err := writer.WriteField("key", fields.Key); err != nil {
		return fmt.Errorf("failed to write key field: %w", err)
	}
	if err := writer.WriteField("x-amz-algorithm", fields.XAmzAlgorithm); err != nil {
		return fmt.Errorf("failed to write x-amz-algorithm field: %w", err)
	}
	if err := writer.WriteField("x-amz-credential", fields.XAmzCredential); err != nil {
		return fmt.Errorf("failed to write x-amz-credential field: %w", err)
	}
	if err := writer.WriteField("x-amz-date", fields.XAmzDate); err != nil {
		return fmt.Errorf("failed to write x-amz-date field: %w", err)
	}
	if err := writer.WriteField("policy", fields.Policy); err != nil {
		return fmt.Errorf("failed to write policy field: %w", err)
	}
	if err := writer.WriteField("x-amz-signature", fields.XAmzSignature); err != nil {
		return fmt.Errorf("failed to write x-amz-signature field: %w", err)
	}
	if err := writer.WriteField("x-amz-security-token", fields.XAmzSecurityToken); err != nil {
		return fmt.Errorf("failed to write x-amz-security-token field: %w", err)
	}

	// Add file
	part, err := writer.CreateFormFile("file", fields.Key)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	// Upload to S3
	s3Req, err := http.NewRequest("POST", fileUploadAmazon, body)
	if err != nil {
		return fmt.Errorf("failed to create S3 request: %w", err)
	}

	s3Req.Header.Set("Content-Type", writer.FormDataContentType())

	s3Resp, err := rh.client.Do(s3Req)
	if err != nil {
		return fmt.Errorf("S3 upload failed: %w", err)
	}
	defer s3Resp.Body.Close()

	if s3Resp.StatusCode != http.StatusNoContent && s3Resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(s3Resp.Body)
		return fmt.Errorf("S3 upload failed with status %d: %s", s3Resp.StatusCode, string(body))
	}

	// Step 3: Add chapter to tonie
	newChapter := Chapter{
		ID:    fields.Key,
		File:  amazonBean.FileID,
		Title: title,
	}

	tonie.Chapters = append(tonie.Chapters, newChapter)

	return nil
}

// disconnect terminates the session
func (rh *requestHandler) disconnect() error {
	req, err := http.NewRequest("DELETE", session, nil)
	if err != nil {
		return fmt.Errorf("failed to create disconnect request: %w", err)
	}

	if rh.jwtToken != nil {
		req.Header.Set("Authorization", "Bearer "+rh.jwtToken.AccessToken)
	}

	resp, err := rh.client.Do(req)
	if err != nil {
		return fmt.Errorf("disconnect request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("disconnect failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// executeGetRequest performs a GET request with authentication
func (rh *requestHandler) executeGetRequest(url string, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if rh.jwtToken != nil {
		req.Header.Set("Authorization", "Bearer "+rh.jwtToken.AccessToken)
	}

	resp, err := rh.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// executePatchRequest performs a PATCH request with authentication
func (rh *requestHandler) executePatchRequest(url string, body []byte) error {
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentTypeJSON)
	if rh.jwtToken != nil {
		req.Header.Set("Authorization", "Bearer "+rh.jwtToken.AccessToken)
	}

	resp, err := rh.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
