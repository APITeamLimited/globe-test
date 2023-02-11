package libWorker

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/dop251/goja"
)

const (
	StandaloneScriptVariant = "standaloneScript"
	HTTPRequestVariant      = "httpRequest"
	GroupVariant            = "group"

	// This is only to be used in paths, the actual script name is the same as the file name
	StandaloneScriptAnonymousName = "standalone.js"
)

type (
	Node interface {
		GetVariant() string // standaloneScript, httpRequest, group
		GetId() string
		GetName() string
		SetSourceId(scriptName string, sourceId string)
		// Map of script name to goja.Callable
		RegisterExports(filename string, exports map[string]goja.Callable) error
		GetScripts() map[string]map[string]goja.Callable
	}

	SourceScript struct {
		// If used in StandaloneScriptNode this is StandaloneScriptAnonymousName
		Name     string `json:"name"`
		Contents string `json:"contents"`
		// Exports is an aray of goja callables
		Exports  map[string]goja.Callable `json:"-"`
		SourceId string                   `json:"-"`
	}

	StandaloneScriptNode struct {
		Variant string       `json:"variant"` // standaloneScript
		Id      string       `json:"id"`
		Name    string       `json:"name"`
		Script  SourceScript `json:"script"`
	}

	HTTPRequestNode struct {
		Variant      string                 `json:"variant"` // httpRequest
		Id           string                 `json:"id"`
		Name         string                 `json:"name"`
		FinalRequest map[string]interface{} `json:"finalRequest"`
		Scripts      []SourceScript         `json:"scripts"`
	}

	GroupNode struct {
		Variant  string         `json:"variant"` // group
		Id       string         `json:"id"`
		Name     string         `json:"name"`
		Scripts  []SourceScript `json:"script"`
		Children []Node         `json:"children"`
	}

	TestData struct {
		RootNode   Node         `json:"rootNode"`
		RootScript SourceScript `json:"rootScript"`
	}
)

func ExtractTestData(rawTestData map[string]interface{}) (*TestData, error) {
	testData := TestData{}

	if rootNode, ok := rawTestData["rootNode"]; ok {
		if node, err := extractNode(rootNode.(map[string]interface{})); err != nil {
			return &testData, err
		} else {
			testData.RootNode = node
		}
	}

	if rootScript, ok := rawTestData["rootScript"]; ok {
		testData.RootScript = SourceScript{
			Name:     rootScript.(map[string]interface{})["name"].(string),
			Contents: rootScript.(map[string]interface{})["contents"].(string),
		}
	} else {
		return &testData, errors.New("rootScript not found")
	}

	return &testData, nil
}

func extractNode(rawNode map[string]interface{}) (Node, error) {
	if variant, ok := rawNode["variant"]; ok {
		switch variant {
		case "standaloneScript":
			standaloneScriptNode := StandaloneScriptNode{
				Variant: "standaloneScript",
				Id:      rawNode["id"].(string),
				Name:    rawNode["name"].(string),
				Script: SourceScript{
					Name:     rawNode["script"].(map[string]interface{})["name"].(string),
					Contents: rawNode["script"].(map[string]interface{})["contents"].(string),
				},
			}
			return &standaloneScriptNode, nil

		case "httpRequest":
			parsedScripts := []SourceScript{}

			if rawScripts, ok := rawNode["scripts"]; ok {
				for _, rawScript := range rawScripts.([]interface{}) {
					parsedScripts = append(parsedScripts, SourceScript{
						Name:     rawScript.(map[string]interface{})["name"].(string),
						Contents: rawScript.(map[string]interface{})["contents"].(string),
					})
				}
			} else {
				return nil, errors.New("no scripts found")
			}

			httpRequestNode := HTTPRequestNode{
				Variant:      "httpRequest",
				Id:           rawNode["id"].(string),
				Name:         rawNode["name"].(string),
				FinalRequest: rawNode["finalRequest"].(map[string]interface{}),
				Scripts:      parsedScripts,
			}

			return &httpRequestNode, nil
		case "group":
			parsedScripts := []SourceScript{}

			if rawScripts, ok := rawNode["scripts"]; ok {
				for _, rawScript := range rawScripts.([]interface{}) {
					parsedScripts = append(parsedScripts, SourceScript{
						Name:     rawScript.(map[string]interface{})["name"].(string),
						Contents: rawScript.(map[string]interface{})["contents"].(string),
					})
				}
			} else {
				return nil, errors.New("no scripts found")
			}

			children := []Node{}

			if rawChildren, ok := rawNode["children"]; ok {
				for _, rawChild := range rawChildren.([]interface{}) {
					if child, err := extractNode(rawChild.(map[string]interface{})); err != nil {
						return nil, err
					} else {
						children = append(children, child)
					}
				}
			} else {
				return nil, errors.New("no children found")
			}

			groupNode := GroupNode{
				Variant:  "group",
				Id:       rawNode["id"].(string),
				Name:     rawNode["name"].(string),
				Scripts:  parsedScripts,
				Children: children,
			}

			return &groupNode, nil
		}
	}

	return nil, errors.New("unknown node type")
}

func (n *StandaloneScriptNode) GetVariant() string {
	return n.Variant
}

func (n *StandaloneScriptNode) GetId() string {
	return n.Id
}

func (n *StandaloneScriptNode) GetName() string {
	return n.Name
}

func (n *StandaloneScriptNode) SetSourceId(scriptName string, sourceId string) {
	if n.Script.Name == scriptName {
		n.Script.SourceId = sourceId
	}
}

func (n *StandaloneScriptNode) RegisterExports(filename string, exports map[string]goja.Callable) error {
	parts := strings.Split(filename, "/")

	urlEncodedLastPart := parts[len(parts)-1]
	urlDecodedLastPart, err := url.QueryUnescape(urlEncodedLastPart)
	if err != nil {
		return err
	}

	if n.Script.Name == urlDecodedLastPart {
		n.Script.Exports = exports
		return nil
	}

	return fmt.Errorf("script name %s does not match filename %s", n.Script.Name, urlDecodedLastPart)
}

func (n *StandaloneScriptNode) GetScripts() map[string]map[string]goja.Callable {
	exports := map[string]map[string]goja.Callable{}

	exports[n.Script.Name] = n.Script.Exports

	return exports
}

func (n *HTTPRequestNode) GetVariant() string {
	return n.Variant
}

func (n *HTTPRequestNode) GetId() string {
	return n.Id
}

func (n *HTTPRequestNode) GetName() string {
	return n.Name
}

func (n *HTTPRequestNode) SetSourceId(scriptName string, sourceId string) {
	for i, script := range n.Scripts {
		if script.Name == scriptName {
			n.Scripts[i].SourceId = sourceId
		}
	}
}

func (n *HTTPRequestNode) RegisterExports(filename string, exports map[string]goja.Callable) error {
	parts := strings.Split(filename, "/")

	urlEncodedLastPart := parts[len(parts)-1]
	urlDecodedLastPart, err := url.QueryUnescape(urlEncodedLastPart)
	if err != nil {
		return err
	}

	for i, script := range n.Scripts {
		if script.Name == urlDecodedLastPart {
			n.Scripts[i].Exports = exports
			return nil
		}
	}

	return fmt.Errorf("could not register exports for %s", filename)
}

func (n *HTTPRequestNode) GetScripts() map[string]map[string]goja.Callable {
	exports := map[string]map[string]goja.Callable{}

	for _, script := range n.Scripts {
		exports[script.Name] = script.Exports
	}

	return exports
}

func (n *GroupNode) GetVariant() string {
	return n.Variant
}

func (n *GroupNode) GetId() string {
	return n.Id
}

func (n *GroupNode) GetName() string {
	return n.Name
}

func (n *GroupNode) SetSourceId(scriptName string, sourceId string) {
	for i, script := range n.Scripts {
		if script.Name == scriptName {
			n.Scripts[i].SourceId = sourceId
		}
	}
}

func (n *GroupNode) RegisterExports(filename string, exports map[string]goja.Callable) error {
	parts := strings.Split(filename, "/")

	urlEncodedLastPart := parts[len(parts)-1]
	urlDecodedLastPart, err := url.QueryUnescape(urlEncodedLastPart)
	if err != nil {
		return err
	}

	for i, script := range n.Scripts {
		if script.Name == urlDecodedLastPart {
			n.Scripts[i].Exports = exports
			return nil
		}
	}

	return fmt.Errorf("could not register exports for %s", filename)
}

func (n *GroupNode) GetScripts() map[string]map[string]goja.Callable {
	exports := map[string]map[string]goja.Callable{}

	for _, script := range n.Scripts {
		exports[script.Name] = script.Exports
	}

	return exports
}

func GetInnerNode(parentNode Node, path string) (Node, error) {
	// Split path at slash
	parts := strings.Split(path, "/")

	// Find first part in node

	parentNodeVariant := parentNode.GetVariant()

	if parentNodeVariant == HTTPRequestVariant {
		if len(parts) > 2 {
			return nil, fmt.Errorf("path %s not found on node %s, variant %s", path, parentNode.GetId(), parentNode.GetVariant())
		}

		return parentNode, nil
	} else if parentNodeVariant == StandaloneScriptVariant {
		if len(parts) > 2 {
			return nil, fmt.Errorf("path %s not found on node %s, variant %s", path, parentNode.GetId(), parentNode.GetVariant())
		}

		return parentNode, nil
	} else if parentNodeVariant == GroupVariant {
		if len(parts) == 2 {
			return parentNode, nil
		}

		// Find child node
		for _, childNode := range parentNode.(*GroupNode).Children {
			if childNode.GetName() == parts[1] {
				return GetInnerNode(childNode, strings.Join(parts[1:], "/"))
			}
		}

		return nil, fmt.Errorf("path %s not found on node %s, variant %s", path, parentNode.GetId(), parentNode.GetVariant())
	}
	return nil, fmt.Errorf("unknown node variant: %s", parentNodeVariant)

}
