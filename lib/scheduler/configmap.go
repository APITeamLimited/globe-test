package scheduler

import (
	"encoding/json"
	"fmt"
	"sync"
)

// ConfigMap can contain mixed scheduler config types
type ConfigMap map[string]Config

// ConfigConstructor is a simple function that returns a concrete Config instance
// with the specified name and all default values correctly initialized
type ConfigConstructor func(name string, rawJSON []byte) (Config, error)

var (
	configTypesMutex   sync.RWMutex
	configConstructors = make(map[string]ConfigConstructor)
)

// RegisterConfigType adds the supplied ConfigConstructor as the constructor for its
// type in the configConstructors map, in a thread-safe manner
func RegisterConfigType(configType string, constructor ConfigConstructor) ***REMOVED***
	configTypesMutex.Lock()
	defer configTypesMutex.Unlock()

	if constructor == nil ***REMOVED***
		panic("scheduler configs: constructor is nil")
	***REMOVED***
	if _, configTypeExists := configConstructors[configType]; configTypeExists ***REMOVED***
		panic("scheduler configs: RegisterConfigType called twice for  " + configType)
	***REMOVED***

	configConstructors[configType] = constructor
***REMOVED***

// GetParsedConfig returns a struct instance corresponding to the supplied
// config type. It will be fully initialized - with both the default values of
// the type, as well as with whatever the user had specified in the JSON
func GetParsedConfig(name, configType string, rawJSON []byte) (result Config, err error) ***REMOVED***
	configTypesMutex.Lock()
	defer configTypesMutex.Unlock()

	constructor, exists := configConstructors[configType]
	if !exists ***REMOVED***
		return nil, fmt.Errorf("unknown execution scheduler type '%s'", configType)
	***REMOVED***
	return constructor(name, rawJSON)
***REMOVED***

// UnmarshalJSON implements the json.Unmarshaler interface in a two-step manner,
// creating the correct type of configs based on the `type` property.
func (scs *ConfigMap) UnmarshalJSON(data []byte) error ***REMOVED***
	if len(data) == 0 ***REMOVED***
		return nil
	***REMOVED***

	if len(data) == 4 && string(data) == "null" ***REMOVED***
		return nil
	***REMOVED***

	//TODO: use a more sophisticated combination of dec.Token() and dec.More(),
	// which would allow us to support both arrays and maps for this config?
	var protoConfigs map[string]protoConfig
	if err := strictJSONUnmarshal(data, &protoConfigs); err != nil ***REMOVED***
		return err
	***REMOVED***

	result := make(ConfigMap, len(protoConfigs))
	for k, v := range protoConfigs ***REMOVED***
		if v.Type == "" ***REMOVED***
			return fmt.Errorf("execution config '%s' doesn't have a type value", k)
		***REMOVED***
		config, err := GetParsedConfig(k, v.Type, v.rawJSON)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		result[k] = config
	***REMOVED***

	*scs = result

	return nil
***REMOVED***

// Validate checks if all of the specified scheduler options make sense
func (scs ConfigMap) Validate() (errors []error) ***REMOVED***
	for name, scheduler := range scs ***REMOVED***
		if schedErr := scheduler.Validate(); len(schedErr) != 0 ***REMOVED***
			errors = append(errors,
				fmt.Errorf("scheduler %s has errors: %s", name, concatErrors(schedErr, ", ")))
		***REMOVED***
	***REMOVED***
	return errors
***REMOVED***

type protoConfig struct ***REMOVED***
	BaseConfig
	rawJSON json.RawMessage
***REMOVED***

// UnmarshalJSON just reads unmarshals the base config (to get the type), but it also
// stores the unprocessed JSON so we can parse the full config in the next step
func (pc *protoConfig) UnmarshalJSON(b []byte) error ***REMOVED***
	*pc = protoConfig***REMOVED***BaseConfig***REMOVED******REMOVED***, b***REMOVED***
	return json.Unmarshal(b, &pc.BaseConfig)
***REMOVED***
