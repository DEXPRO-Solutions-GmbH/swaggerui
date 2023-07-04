package swaggerui

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Spec is a custom type used to represent a swagger spec. It is mostly unvalidated and not
// more type-safe than a regular map[string]any, but it allows us to add custom methods.
//
// Nested objects are also expected to be of type Spec.
type Spec map[string]any

// SetOpenIdConnectUrl replaces the url of a given OpenID Connect security scheme.
//
// This method panics if any keys of the underlying Spec don't exist or are not of type Spec.
//
// This method also panics if the given scheme is not of type "openIdConnect".
func (spec Spec) SetOpenIdConnectUrl(securitySchemeName string, url string) {
	comps, has := spec["components"]
	if !has {
		comps = make(Spec)
		spec["components"] = comps
	}

	schemes, has := comps.(Spec)["securitySchemes"]
	if !has {
		schemes = make(Spec)
		spec["components"].(Spec)["securitySchemes"] = schemes
	}

	scheme, has := schemes.(Spec)[securitySchemeName]
	if !has {
		scheme = make(Spec)
		spec["components"].(Spec)["securitySchemes"].(Spec)[securitySchemeName] = scheme
	}

	schemeType, has := scheme.(Spec)["type"]
	if !has {
		schemeType = make(Spec)
		spec["components"].(Spec)["securitySchemes"].(Spec)[securitySchemeName].(Spec)["type"] = "openIdConnect"
	} else if schemeType != "openIdConnect" {
		panic(fmt.Errorf("components.securitySchemes.%s is not of type openIdConnect", securitySchemeName))
	}

	scheme.(Spec)["openIdConnectUrl"] = url
}

// AddServerUrl adds a url element to the servers list.
// The last url added will be the first one to be displayed in the Swagger UI, making it the default server.
func (spec Spec) AddServerUrl(serverUrl string) {
	specServers, has := spec["servers"]
	if !has {
		spec["servers"] = make([]any, 0)
	}

	specServersEntry := make([]any, 0)

	// add serverUrl to servers key
	entry := make(Spec)
	entry["url"] = serverUrl
	specServersEntry = append(specServersEntry, entry)

	// add existing urls to servers key
	specServersList, exist := specServers.([]any)
	if exist {
		for _, value := range specServersList {
			valueList, exist := value.(Spec)
			if exist {
				url, exist := valueList["url"]
				if exist {
					entry := make(Spec)
					entry["url"] = url
					specServersEntry = append(specServersEntry, entry)
				}
			}
		}
	}
	specServers = specServersEntry
	spec["servers"] = specServers
}

func ParseSpecYAML(raw []byte) (Spec, error) {
	var spec Spec
	if err := yaml.Unmarshal(raw, &spec); err != nil {
		return nil, err
	}
	return spec, nil
}
