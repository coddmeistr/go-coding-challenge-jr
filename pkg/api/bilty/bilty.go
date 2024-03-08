// Package bilty prodives API integration with bilty API service
// Package have been tested with unit tests
package bilty

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrInternal = errors.New("internal library error")
	ErrApiError = errors.New("api error")
)

const (
	host = "https://api-ssl.bitly.com/v4"

	shortenUrl = "/shorten"
)

type Bilty struct {
	client *http.Client
	Token  string `json:"token"`
}

func NewBilty(token string, client *http.Client) *Bilty {
	return &Bilty{
		client: client,
		Token:  token,
	}
}

// CreateShortLink creates a short link for the given long URL.
//
// It takes a longUrl string as a parameter and returns a shortened url string and an error.
func (b *Bilty) CreateShortLink(longUrl string) (string, error) {

	bts, err := json.Marshal(CreateLinkRequest{
		Link: longUrl,
	})
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInternal, err)
	}

	req, err := http.NewRequest("POST", host+shortenUrl, bytes.NewReader(bts))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInternal, err)
	}
	req.Header.Set("Authorization", "Bearer "+b.Token)

	resp, err := b.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInternal, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInternal, err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		var apiError ErrorMessage
		if err := json.Unmarshal(body, &apiError); err != nil {
			return "", fmt.Errorf("%w: %v", ErrInternal, err)
		}

		return "", fmt.Errorf("%w: %v", ErrApiError, apiError.Message)
	}

	var response CreateLinkResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return response.ShortLink, nil
}
