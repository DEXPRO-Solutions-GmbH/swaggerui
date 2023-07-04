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

func (handler *Handler) Register(router gin.IRoutes) {
	//router.GET("/openapi.yml", handler.GetSpec)
	router.GET("/openapi.yml", handler.GetSpec)
	router.StaticFS("/swagger-ui", handler.FS)
}

type Option func(handler *Handler)

func WithOIDC(securitySchemeName, oidcURL string) Option {
	return func(handler *Handler) {
		handler.middlewares = append(handler.middlewares, func(ctx *gin.Context, spec Spec) {
			spec.SetOpenIdConnectUrl(securitySchemeName, oidcURL)
		})
	}
}

func WithAddServerUrls() Option {
	return func(handler *Handler) {
		handler.middlewares = append(handler.middlewares, func(ctx *gin.Context, spec Spec) {
			path := ctx.Request.URL.Path
			path = strings.TrimSuffix(path, "/openapi.yml")
			httpsUrl := fmt.Sprintf("https://%s%s", ctx.Request.Host, path)
			httpUrl := fmt.Sprintf("http://%s%s", ctx.Request.Host, path)
			spec.AddServerUrl(httpUrl)
			spec.AddServerUrl(httpsUrl)
		})
	}
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
