package bridge

import (
	"fmt"
	"strings"
)

type Module struct ***REMOVED***
	Name    string
	Members map[string]Func
***REMOVED***

func (m *Module) JS() string ***REMOVED***
	jsFuncs := []string***REMOVED******REMOVED***
	for name, fn := range m.Members ***REMOVED***
		jsFuncs = append(jsFuncs, fmt.Sprintf(`'%s': %s`, name, fn.JS(m.Name, name)))
	***REMOVED***
	return fmt.Sprintf("__internal__._register('%s', ***REMOVED***\n%s\n***REMOVED***);", m.Name, strings.Join(jsFuncs, ",\n"))
***REMOVED***

func BridgeModule(name string, members map[string]interface***REMOVED******REMOVED***) Module ***REMOVED***
	mod := Module***REMOVED***
		Name:    name,
		Members: make(map[string]Func),
	***REMOVED***
	for name, mem := range members ***REMOVED***
		mod.Members[name] = BridgeFunc(mem)
	***REMOVED***
	return mod
***REMOVED***
