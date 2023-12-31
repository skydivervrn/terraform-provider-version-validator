package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	versionValidatorDataSourceResourceName = "version_validator"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{},
			DataSourcesMap: map[string]*schema.Resource{
				versionValidatorDataSourceResourceName: dataSourceVersionValidator(),
			},
			ResourcesMap: map[string]*schema.Resource{},
		}

		//p.ConfigureContextFunc = configure(version, p)

		return p
	}
}
