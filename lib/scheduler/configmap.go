package scheduler

import (
	"encoding/json"
	"fmt"
)

// GetParsedConfig returns a struct instance corresponding to the supplied
// config type. It will be fully initialized - with both the default values of
// the type, as well as with whatever the user had specified in the JSON
func GetParsedConfig(name, configType string, rawJSON []byte) (result Config, err error) ***REMOVED***
	switch configType ***REMOVED***
	case constantLoopingVUsType:
		config := NewConstantLoopingVUsConfig(name)
		err = json.Unmarshal(rawJSON, &config)
		result = config
	case variableLoopingVUsType:
		config := NewVariableLoopingVUsConfig(name)
		err = json.Unmarshal(rawJSON, &config)
		result = config
	case sharedIterationsType:
		config := NewSharedIterationsConfig(name)
		err = json.Unmarshal(rawJSON, &config)
		result = config
	case perVUIterationsType:
		config := NewPerVUIterationsConfig(name)
		err = json.Unmarshal(rawJSON, &config)
		result = config
	case constantArrivalRateType:
		config := NewConstantArrivalRateConfig(name)
		err = json.Unmarshal(rawJSON, &config)
		result = config
	case variableArrivalRateType:
		config := NewVariableArrivalRateConfig(name)
		err = json.Unmarshal(rawJSON, &config)
		result = config
	default:
		return nil, fmt.Errorf("unknown execution scheduler type '%s'", configType)
	***REMOVED***
	return
***REMOVED***

// ConfigMap can contain mixed scheduler config types
type ConfigMap map[string]Config

// UnmarshalJSON implements the json.Unmarshaler interface in a two-step manner,
// creating the correct type of configs based on the `type` property.
func (scs *ConfigMap) UnmarshalJSON(b []byte) error ***REMOVED***
	var protoConfigs map[string]protoConfig
	if err := json.Unmarshal(b, &protoConfigs); err != nil ***REMOVED***
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
