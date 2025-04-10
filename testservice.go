package cicdwf

import (
	"context"
	"errors"
	"io"
	"net/http"
)

type TestService struct {
	client *http.Client
}

func NewTestService() *TestService {
	return &TestService{
		client: &http.Client{},
	}
}

func (testService *TestService) GetServiceResult(ctx context.Context, endpoint string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := testService.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("service returned non-200 status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
