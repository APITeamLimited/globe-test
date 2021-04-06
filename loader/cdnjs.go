/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

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
