package tests

import (
	"challenge/pkg/proto"
	suits "challenge/tests/suit"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

// This test testing entire gRPC flow, otherwise the test in bilty package
// testing the same thing, but isolated
func TestUrlShorten_OkWithValidRedirect(t *testing.T) {
	_, s := suits.NewDefault(t)

	validUrl := "https://github.com/maxik12233"

	link, err := s.Client.MakeShortLink(context.Background(), &proto.Link{Data: validUrl})
	require.NoError(t, err)

	// Make a request on gotten shortened url
	// And make sure that it redirects to original url
	req, err := http.NewRequest("GET", link.GetData(), nil)
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
	assert.Equal(t, validUrl, locHeader)
}
