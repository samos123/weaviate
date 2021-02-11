package modulecapabilities

import (
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/schema"
)

// ClassConfig is a helper type which is passed TO the module to read it's
// per-class config
type ClassConfig interface {
	Class() map[string]interface{}
	Property(propName string) map[string]interface{}
}

// ClassConfigurator is an optional capability interface which a module MAY
// implement. If it is implemented, all methods will be called when the user
// adds or updates a class which has the module set as the vectorizer
type ClassConfigurator interface {

	// ClassDefaults provides the defaults for a per-class module config. The
	// module provider will merge the props into the user-specified config with
	// the user-provided values taking precedence
	ClassConfigDefaults() map[string]interface{}

	// PropertyConfigDefaults provides the defaults for a per-property module
	// config. The module provider will merge the props into the user-specified
	// config with the user-provided values taking precedence. The property's
	// dataType MAY be taken into consideration when deciding defaults.
	// dataType is not guaranteed to be non-nil, it might be nil in the case a
	// user specified an invalid dataType, as some validation only occurs after
	// defaults are set.
	PropertyConfigDefaults(dataType *schema.DataType) map[string]interface{}

	// ValidateClass MAY validate anything about the class, except the config of
	// another module. The specified ClassConfig can be used to easily retrieve
	// the config specific for the module. For example, a module could iterate
	// over class.Properties and call classConfig.Property(prop.Name) to validate
	// the per-property config. A module MUST NOT extract another module's config
	// from class.ModuleConfig["other-modules-name"].
	ValidateClass(class *models.Class, classConfig ClassConfig) error
}
