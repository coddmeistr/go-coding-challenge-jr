package bilty

import (
	"bytes"
	"challenge/pkg/config"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (s roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return s(r)
}

func newClientMock(t *testing.T, statusCode int, path string, response any) *Bilty {
	return &Bilty{
		client: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, path, r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				respBody, err := json.Marshal(response)
				if err != nil {
					assert.Fail(t, "Cannot read bytes")
				}
				return &http.Response{
					StatusCode: statusCode,
					Body:       io.NopCloser(bytes.NewReader(respBody)),
				}, nil
			}),
		},
	}
}

func TestCreateShortLink_TestCases(t *testing.T) {

	tc := []struct {
		name string

		origin string

		expectedStatusCode int
		expectedResponse   any
		wantErr            bool
		wantErrMsg         string
	}{
		{
			name:               "ok",
			origin:             "https://www.google.com/",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   CreateLinkResponse{ShortLink: "bit.ly/something"},
			wantErr:            false,
			wantErrMsg:         "",
		},
		{
			name:               "some error",
			origin:             "https://www.google.com/",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorMessage{Message: "something goes wrong"},
			wantErr:            true,
			wantErrMsg:         "something goes wrong",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			bil := newClientMock(t, tt.expectedStatusCode, shortenUrl, tt.expectedResponse)

			got, err := bil.CreateShortLink(tt.origin)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
				resp, ok := tt.expectedResponse.(CreateLinkResponse)
				require.True(t, ok)
				assert.Equal(t, resp.ShortLink, got)
			}
		})
	}
}

// API test
func TestBitly_TestCases(t *testing.T) {

	// Load environment variables in viper from context and from file
	viper.AutomaticEnv()
	envPath := "../../../.env"
	if err := config.ReadAndParseFromFile(envPath, nil); err != nil {
		fmt.Printf(".env file was not found in %s\n", envPath)
	}
	token := viper.GetString("BITLY_OAUTH_TOKEN")
	require.NotEqual(t, "", token)

	// Test contains only ONE testcase, because we dont want to exceed url creation limit so fast
	tc := []struct {
		name       string
		origin     string
		httpClient *http.Client
	}{
		{
			name:       "ok",
			origin:     "https://www.google.com/",
			httpClient: http.DefaultClient,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			bil := NewBilty(token, tt.httpClient)

			shortened, err := bil.CreateShortLink(tt.origin)
			require.NoError(t, err)

			req, err := http.NewRequest("GET", shortened, nil)
			require.NoError(t, err)
			clientWithNoRedirect := http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			resp, err := clientWithNoRedirect.Do(req)
			require.NoError(t, err)

			locHeader := resp.Header.Get("Location")
			assert.Equal(t, resp.StatusCode, 301)
			assert.Equal(t, tt.origin, locHeader)
		})
	}
}
