package swaggerui

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

//go:embed dist
var distFs embed.FS

// Handler is responsible for both serving the SwaggerUI and the associated OpenAPI
// spec.
type Handler struct {
	specYML []byte
	FS      http.FileSystem

	middlewares []HandlerMiddleware
}

// NewHandler returns a new Handler based on the given spec. It is expected that
// the given data contains valid YAML.
func NewHandler(specYML []byte, opts ...Option) (*Handler, error) {
	// parse yml spec - this validates that the given string is valid yaml
	var spec Spec
	err := yaml.Unmarshal(specYML, &spec)
	if err != nil {
		return nil, err
	}

	// setup virtual fs for swagger-ui
	subFs, err := fs.Sub(distFs, "dist")
	if err != nil {
		panic(err)
	}

	h := &Handler{
		specYML: specYML,
		FS:      http.FS(subFs),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h, nil
}

// Register is used to register the Handler on a gin router.
func (handler *Handler) Register(router gin.IRoutes) {
	router.GET("/openapi.yml", handler.GetSpec)
	router.StaticFS("/swagger-ui", handler.FS)
	router.GET("/swaggerui", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, getRedirectPath(ctx))
	})
}

// Option objects used to construct Handler objects.
type Option func(handler *Handler)

// WithOIDC adds a HandlerMiddleware which replaces the OpenIDConnect URL of a given security scheme
// with the given url. This is useful if your API is secured by an OIDC provider and you want
// to make the provider your API trusts configurable.
func WithOIDC(securitySchemeName, oidcURL string) Option {
	return func(handler *Handler) {
		handler.middlewares = append(handler.middlewares, func(ctx *gin.Context, spec Spec) {
			spec.SetOpenIdConnectUrl(securitySchemeName, oidcURL)
		})
	}
}

// WithAddServerUrls adds a HandlerMiddleware which adds the requests host to the OpenAPI specs
// server list, once with HTTP, once with HTTPS.
//
// This is meant to be used on APIs which could be served under any hostname and where
// you don't want to hardcode the hostname in your spec.
func WithAddServerUrls() Option {
	return func(handler *Handler) {
		mw := newServerURLMiddleware(false)
		handler.middlewares = append(handler.middlewares, mw)
	}
}

// WithReplaceServerUrls is the same as WithAddServerUrls but it will replace all previously
// defined server urls.
func WithReplaceServerUrls() Option {
	return func(handler *Handler) {
		mw := newServerURLMiddleware(true)
		handler.middlewares = append(handler.middlewares, mw)
	}
}

// newServerURLMiddleware returns a new HandlerMiddleware which adds the requests host to the OpenAPI specs.
//
// The replace parameter controls whether the middleware should replace all previously defined server urls or simply add
// to the list.
func newServerURLMiddleware(replace bool) HandlerMiddleware {
	mw := func(ctx *gin.Context, spec Spec) {
		path := ctx.Request.URL.Path
		path = strings.TrimSuffix(path, "/openapi.yml")
		httpsUrl := fmt.Sprintf("https://%s%s", ctx.Request.Host, path)
		httpUrl := fmt.Sprintf("http://%s%s", ctx.Request.Host, path)
		if replace {
			spec.RemoveServerURLs()
		}
		spec.AddServerUrl(httpUrl)
		spec.AddServerUrl(httpsUrl)
	}
	return mw
}

// WithMiddleware applies the given HandlerMiddleware to every request.
func WithMiddleware(mw HandlerMiddleware) Option {
	return func(handler *Handler) {
		handler.middlewares = append(handler.middlewares, mw)
	}
}

// A HandlerMiddleware can modify the given Spec. You may use the gin.Context if you require any
// information from the request.
type HandlerMiddleware func(ctx *gin.Context, spec Spec)

// GetSpec is the gin handler function used to serve the OpenAPI spec
// of this handler.
func (handler *Handler) GetSpec(ctx *gin.Context) {
	// Note: The spec is unmarshaled from it's yml source on every request.
	// This is on purpose because it allows us to have an unmodified, request-scoped copy of the spec
	// which can be modified based on the incoming request. This also avoids any concurrency issues related to that.
	var spec Spec
	err := yaml.Unmarshal(handler.specYML, &spec)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "unexpected, do not forget error message"})
		return
	}

	for _, middleware := range handler.middlewares {
		middleware(ctx, spec)
	}

	ctx.YAML(200, spec)
}

func getRedirectPath(_ *gin.Context) string {
	return "/swagger-ui"
}
