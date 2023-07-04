package swaggerui

import (
	_ "embed"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"testing"
)

//go:embed example.openapi.yml
var exampleSpec []byte

func TestNewHandlerParsed(t *testing.T) {
	t.Skip("This example is only for manual testing of the handler.")

	engine := gin.Default()
	engine.RedirectTrailingSlash = true

	handler, err := NewHandler(
		exampleSpec,
		WithOIDC("OAuth", "https://keycloak.k8s.staging.squeeze.one/realms/dexpro-dev/.well-known/openid-configuration"),
		WithAddServerUrls())

	require.NoError(t, err)
	handler.Register(engine)

	fmt.Println("Try these URLS")
	fmt.Printf("- openapi: %s/openapi.yml\n", "http://localhost:8044")
	fmt.Printf("- swagger ui: %s/swagger-ui\n", "http://localhost:8044")

	_ = engine.Run(":8044")

	select {}
}
