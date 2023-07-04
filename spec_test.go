package swaggerui

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpec_SetOpenIdConnectUrl(t *testing.T) {
	var spec Spec

	setup := func() {
		spec = Spec{
			"components": Spec{
				"securitySchemes": Spec{
					"myOIDC": Spec{
						"type": "openIdConnect",
					},
				},
			},
		}
	}

	t.Run("creates nonexistant schene", func(t *testing.T) {
		setup()
		assert.NotPanics(t, func() {
			spec.SetOpenIdConnectUrl("nonexistant", "https://example.com/.well-known/openid-configuration")
		})
	})

	t.Run("creates missing components", func(t *testing.T) {
		spec = Spec(map[string]any{})
		assert.NotPanics(t, func() {
			spec.SetOpenIdConnectUrl("myOIDC", "https://example.com/.well-known/openid-configuration")
		})
	})

	t.Run("creates missing securitySchemes", func(t *testing.T) {
		spec = Spec(map[string]any{
			"components": map[string]any{},
		})
		assert.Panics(t, func() {
			spec.SetOpenIdConnectUrl("myOIDC", "https://example.com/.well-known/openid-configuration")
		})
	})

	t.Run("panics on invalid security scheme type", func(t *testing.T) {
		spec = Spec{
			"components": map[string]any{
				"securitySchemes": map[string]any{
					"myOIDC": map[string]any{
						"type": "oauth2",
					},
				},
			},
		}
		assert.Panics(t, func() {
			spec.SetOpenIdConnectUrl("myOIDC", "https://example.com/.well-known/openid-configuration")
		})
	})

	t.Run("works", func(t *testing.T) {
		setup()
		spec.SetOpenIdConnectUrl("myOIDC", "https://example.com/.well-known/openid-configuration")
		assert.Equal(t, "https://example.com/.well-known/openid-configuration", spec["components"].(Spec)["securitySchemes"].(Spec)["myOIDC"].(Spec)["openIdConnectUrl"])
	})
}

func TestParseSpec(t *testing.T) {
	t.Run("types are as expected", func(t *testing.T) {
		spec, err := ParseSpecYAML(exampleSpec)
		require.NoError(t, err)
		require.NotNil(t, spec)

		assert.IsType(t, Spec{}, spec, "returned object has unexpected type")
		assert.IsType(t, Spec{}, spec["info"], "child properties of spec have unexpected type")
	})
}
