package eam

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const (
	JsonDataFormat = "json"
	TextDataFormat = "text"

	defaultTimeout = 30 * time.Second
)

type EAMClient struct {
	baseURL    string
	dataFormat string
	httpClient *http.Client
}

func NewEAMClient(baseURL string, userName, passwd string) (*EAMClient, error) {
	jarCookie, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}

	apiClient := &EAMClient{
		baseURL:    baseURL,
		dataFormat: JsonDataFormat,
		httpClient: &http.Client{
			Jar:     jarCookie,
			Timeout: defaultTimeout,
		},
	}

	if err := apiClient.Login(userName, passwd); err != nil {
		return nil, err
	}

	return apiClient, nil
}

func (c *EAMClient) Init(opts ...Option) error {
	for i := range opts {
		opt := opts[i]
		if err := opt(c); err != nil {
			return fmt.Errorf("api client initializing: %w", err)
		}
	}

	return nil
}

func (c EAMClient) Login(userName, passwd string) error {
	requestURL := ""
	requestData := url.Values{}
	requestData.Set("username", userName)
	requestData.Set("password", passwd)

	response, err := c.httpClient.PostForm(requestURL, requestData)
	if err != nil {
		return fmt.Errorf("login request to api service: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ErrAPIServiceLogin
	}

	return nil
}
