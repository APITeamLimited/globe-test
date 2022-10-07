package worker

/*
Checks if any stored data is required for execution
*/
/*func loadStoredData(gs *globalState, job libOrch.ChildJob, fs afero.Fs, workerInfo *libWorker.WorkerInfo) error ***REMOVED***
	// Check if we have any stored data
	storedData := job["storedData"]
	if storedData == "" ***REMOVED***
		return nil
	***REMOVED***

	var storedDataMap map[string]string
	err := json.Unmarshal([]byte(storedData), &storedDataMap)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Check if we have any stored data
	if len(storedDataMap) == 0 ***REMOVED***
		return nil
	***REMOVED***

	for key, value := range storedDataMap ***REMOVED***
		// Write the file to the filesystem
		filePath := filepath.Join(gs.workDir, key)
		err := afero.WriteFile(fs, filePath, []byte(value), 0644)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Add the file to the files map
		workerInfo.Files[key] = filePath
	***REMOVED***

	return nil

***REMOVED***
*/
