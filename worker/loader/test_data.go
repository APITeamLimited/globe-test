package loader

import (
	"fmt"
	"net/url"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/google/uuid"
)

func LoadTestData(testData *libWorker.TestData) (*[]*SourceData, error) {
	sourceData, err := recursiveAddSourceData(testData.RootNode, "", testData.RootScript.Name)
	if err != nil {
		return nil, err
	}

	// Ensure exactly one root source
	rootSourceCount := 0
	for _, data := range sourceData {
		if data.RootSource {
			rootSourceCount++
		}
	}

	if rootSourceCount != 1 {
		return nil, fmt.Errorf("expected exactly one root source, got %d, this can occur if multiple scripts in the root node have the same name", rootSourceCount)
	}

	// Ensure root source is first
	for i, data := range sourceData {
		if data.RootSource {
			sourceData[0], sourceData[i] = sourceData[i], sourceData[0]
			break
		}
	}

	return &sourceData, nil
}

func recursiveAddSourceData(node libWorker.Node, existingPath string, rootScriptName string) ([]*SourceData, error) {
	variant := node.GetVariant()

	sourceData := make([]*SourceData, 0)

	switch variant {
	case libWorker.StandaloneScriptVariant:
		standaloneScriptNode := node.(*libWorker.StandaloneScriptNode)

		newData := SourceData{
			Data:     []byte(standaloneScriptNode.Script.Contents),
			URL:      &url.URL{Path: fmt.Sprintf("%s%s/%s", existingPath, standaloneScriptNode.Script.Name, libWorker.StandaloneScriptAnonymousName)},
			SourceId: uuid.New().String(),
		}

		if existingPath == "" && standaloneScriptNode.Script.Name == rootScriptName {
			newData.RootSource = true
		}

		node.SetSourceId(standaloneScriptNode.Script.Name, newData.SourceId)

		sourceData = append(sourceData, &newData)

		return sourceData, nil
	case libWorker.HTTPRequestVariant:
		httpRequestNode := node.(*libWorker.HTTPRequestNode)

		for _, script := range httpRequestNode.Scripts {
			newData := SourceData{
				Data:     []byte(script.Contents),
				URL:      &url.URL{Path: fmt.Sprintf("%s%s/%s", existingPath, node.GetName(), script.Name)},
				SourceId: uuid.New().String(),
			}

			if existingPath == "" && script.Name == rootScriptName {
				newData.RootSource = true
			}

			node.SetSourceId(script.Name, newData.SourceId)

			sourceData = append(sourceData, &newData)
		}

		return sourceData, nil
	case libWorker.GroupVariant:
		groupNode := node.(*libWorker.GroupNode)

		for _, script := range groupNode.Scripts {
			newData := SourceData{
				Data:     []byte(script.Contents),
				URL:      &url.URL{Path: fmt.Sprintf("%s%s/%s", existingPath, node.GetName(), script.Name)},
				SourceId: uuid.New().String(),
			}

			if existingPath == "" && script.Name == rootScriptName {
				newData.RootSource = true
			}

			node.SetSourceId(script.Name, newData.SourceId)

			sourceData = append(sourceData, &newData)
		}

		for _, child := range groupNode.Children {
			newPath := fmt.Sprintf("%s%s/", existingPath, groupNode.Name)

			childSourceData, err := recursiveAddSourceData(child, newPath, rootScriptName)
			if err != nil {
				return nil, err
			}

			sourceData = append(sourceData, childSourceData...)
		}

		return sourceData, nil
	default:
		return nil, fmt.Errorf("unknown node variant: %s", variant)
	}
}
