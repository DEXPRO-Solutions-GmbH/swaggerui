package swaggerui

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

//go:embed example.openapi.yml
var exampleSpec []byte

func TestHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.Default()
	engine.RedirectTrailingSlash = true

	handler, err := NewHandler(
		exampleSpec,
		WithOIDC("OAuth", "https://your.idp.domain/realms/some-realm/.well-known/openid-configuration"),
		WithAddServerUrls())

	require.NoError(t, err)
	handler.Register(engine)

	testserver := httptest.NewServer(engine)
	defer testserver.Close()

	t.Run("responds with YML", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/openapi.yml", testserver.URL))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/x-yaml; charset=utf-8", res.Header.Get("Content-Type"))

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(body), len(exampleSpec), "handler should respond with spec which is at least as long as the example spec")
	})
}

func TestGetRedirectPath(t *testing.T) {
	require.Equal(t, "/swagger-ui", getRedirectPath(nil))
}
