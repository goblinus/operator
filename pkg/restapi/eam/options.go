package eam

import "time"

type Option func(api *EAMClient) error

func WithBaseURL(apiBaseURL string) Option {
	return func(api *EAMClient) error {
		api.baseURL = apiBaseURL
		return nil
	}
}

func WithUserPasswd(user, passwd string) Option {
	return func(api *EAMClient) error {
		return api.Login(user, passwd)
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(api *EAMClient) error {
		api.httpClient.Timeout = timeout
		return nil
	}
}

func WithDataFormat(expectedDataFormat string) Option {
	return func(api *EAMClient) error {
		if expectedDataFormat == JsonDataFormat || expectedDataFormat == TextDataFormat {
			api.dataFormat = expectedDataFormat
			return nil
		}

		return ErrUnexpectedDataFormat
	}
}
