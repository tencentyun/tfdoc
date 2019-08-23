package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	// alicloud "github.com/terraform-providers/terraform-provider-alicloud/alicloud"
	// aws "github.com/terraform-providers/terraform-provider-aws/aws"
	// azurerm "github.com/terraform-providers/terraform-provider-azurerm/azurerm"
	// google "github.com/terraform-providers/terraform-provider-google/google"
	tencentcloud "github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud"
)

const LISTENPORT = "127.0.0.1:8080"

var PROVIDERS = []string{
	// "alicloud",
	// "aws",
	// "azurerm",
	// "google",
	"tencentcloud",
}

var TypeMap = map[schema.ValueType]string{
	schema.TypeInvalid: "invalid",
	schema.TypeBool:    "bool",
	schema.TypeInt:     "int",
	schema.TypeFloat:   "float",
	schema.TypeString:  "string",
	schema.TypeList:    "list",
	schema.TypeMap:     "map",
	schema.TypeSet:     "set",
}

type Provider interface {
	Provider() *schema.Provider
}

func NewProvider(name string) (*schema.Provider, error) {
	switch name {
	// case "alicloud":
	// 	return alicloud.Provider().(*schema.Provider), nil
	// case "aws":
	// 	return aws.Provider().(*schema.Provider), nil
	// case "azurerm":
	//	return azurerm.Provider().(*schema.Provider), nil
	// case "google":
	//	return google.Provider().(*schema.Provider), nil
	case "tencentcloud":
		return tencentcloud.Provider(), nil
	default:
		return nil, fmt.Errorf("provider %s is not supported", name)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	http.HandleFunc("/usage", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, USAGE)
	})

	http.HandleFunc("/provider", func(w http.ResponseWriter, r *http.Request) {
		type Result struct {
			Code    int      `json:"code"`
			Message string   `json:"message"`
			Data    []string `json:"data"`
		}

		result := Result{0, "", PROVIDERS}
		data, _ := json.MarshalIndent(result, "", "    ")
		fmt.Fprintf(w, string(data))
	})

	http.HandleFunc("/provider/", func(w http.ResponseWriter, r *http.Request) {

		name := strings.TrimSpace(strings.Trim(r.URL.Path[10:], "/"))
		if name == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var data []byte
		if strings.Contains(name, "/") {
			ns := strings.Split(name, "/")
			name = ns[0]
			res := ns[1]

			type Result struct {
				Code       int                               `json:"code"`
				Message    string                            `json:"message"`
				Arguments  map[string]map[string]interface{} `json:"arguments"`
				Attributes map[string]map[string]interface{} `json:"attributes"`
			}
			result := Result{0, "", map[string]map[string]interface{}{}, map[string]map[string]interface{}{}}
			provider, err := NewProvider(name)
			if err != nil {
				result.Code = 1
				result.Message = err.Error()
			} else {
				found := false
				for k, v := range provider.DataSourcesMap {
					if v.DeprecationMessage != "" {
						continue
					}
					if k == res {
						for kk, vv := range v.Schema {
							if vv.Required || vv.Optional {
								result.Arguments[kk] = getSchema("Arguments", kk, vv)
							} else if vv.Computed {
								result.Attributes[kk] = getSchema("Attributes", kk, vv)
							}
						}
						found = true
						break
					}
				}
				if !found {
					for k, v := range provider.ResourcesMap {
						if v.DeprecationMessage != "" {
							continue
						}
						if k == res {
							for kk, vv := range v.Schema {
								if vv.Required || vv.Optional {
									result.Arguments[kk] = getSchema("Arguments", kk, vv)
								} else if vv.Computed {
									result.Attributes[kk] = getSchema("Attributes", kk, vv)
								}
							}
							found = true
							break
						}
					}
				}
			}
			data, _ = json.MarshalIndent(result, "", "    ")
		} else {
			type Result struct {
				Code        int      `json:"code"`
				Message     string   `json:"message"`
				DataSources []string `json:"data_sources"`
				Resources   []string `json:"resources"`
			}
			result := Result{0, "", []string{}, []string{}}
			provider, err := NewProvider(name)
			if err != nil {
				result.Code = 1
				result.Message = err.Error()
			} else {
				for k, v := range provider.DataSourcesMap {
					if v.DeprecationMessage != "" {
						continue
					}
					result.DataSources = append(result.DataSources, k)
				}

				for k, v := range provider.ResourcesMap {
					if v.DeprecationMessage != "" {
						continue
					}
					result.Resources = append(result.Resources, k)
				}
			}
			data, _ = json.MarshalIndent(result, "", "    ")
		}

		fmt.Fprintf(w, string(data))
	})

	fmt.Printf("HTTP Listen on %s\n", LISTENPORT)
	http.ListenAndServe(LISTENPORT, nil)
}

func getSchema(t string, k string, v *schema.Schema) map[string]interface{} {
	r := map[string]interface{}{}

	if _, ok := v.Elem.(*schema.Resource); ok {
		rr := map[string]map[string]interface{}{}
		for m, n := range v.Elem.(*schema.Resource).Schema {
			vv := getSchema(t, m, n)
			if len(vv) > 0 {
				rr[m] = vv
			}
		}
		if t == "Arguments" {
			r = map[string]interface{}{
				"required":    v.Required,
				"type":        TypeMap[v.Type],
				"description": v.Description,
				"list":        rr,
			}
		} else {
			r = map[string]interface{}{
				"type":        TypeMap[v.Type],
				"description": v.Description,
				"list":        rr,
			}
		}
	} else {
		if t == "Arguments" {
			r = map[string]interface{}{
				"required":    v.Required,
				"type":        TypeMap[v.Type],
				"description": v.Description,
			}
		} else {
			r = map[string]interface{}{
				"type":        TypeMap[v.Type],
				"description": v.Description,
			}
		}
	}

	return r
}
