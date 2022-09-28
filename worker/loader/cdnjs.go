package loader

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

type cdnjsEnvelope struct ***REMOVED***
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Version  string `json:"version"`
	Assets   []struct ***REMOVED***
		Version string   `json:"version"`
		Files   []string `json:"files"`
	***REMOVED***
***REMOVED***

func cdnjs(logger logrus.FieldLogger, path string, parts []string) (string, error) ***REMOVED***
	name := parts[0]
	version := parts[1]
	filename := parts[2]

	data, err := fetch(logger, "https://api.cdnjs.com/libraries/"+name)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	var envelope cdnjsEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	// CDNJS doesn't actually send 404s, nonexistent libs' data is just *empty*.
	if envelope.Name == "" ***REMOVED***
		return "", fmt.Errorf("cdnjs: no such library: %s", name)
	***REMOVED***

	// If no version is specified, use the default/latest one.
	if version == "" ***REMOVED***
		version = envelope.Version
	***REMOVED***

	// If no filename is specified, use the default one, but make sure it actually exists in the
	// chosen version (it may have changed name over the years). If not, the first listed file
	// that does exist in that version is a pretty safe guess.
	if filename == "" ***REMOVED***
		filename = envelope.Filename

		backupFilename := filename
		filenameExistsInVersion := false
		for _, ver := range envelope.Assets ***REMOVED***
			if ver.Version != version ***REMOVED***
				continue
			***REMOVED***
			if len(ver.Files) == 0 ***REMOVED***
				return "",
					fmt.Errorf("cdnjs: no files for version %s of %s, this is a problem with the library or cdnjs not k6",
						version, path)
			***REMOVED***
			backupFilename = ver.Files[0]
			for _, file := range ver.Files ***REMOVED***
				if file == filename ***REMOVED***
					filenameExistsInVersion = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if !filenameExistsInVersion ***REMOVED***
			filename = backupFilename
		***REMOVED***
	***REMOVED***

	return "https://cdnjs.cloudflare.com/ajax/libs/" + name + "/" + version + "/" + filename, nil
***REMOVED***
