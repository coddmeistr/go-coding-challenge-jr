package bilty

import (
	"challenge/pkg/config"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

// Checks if provided shortened url actually leads to the same original url
func TestBilty_CreateShortLink(t *testing.T) {

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
