package worker

/*
Checks if any stored data is required for execution
*/
/*func loadStoredData(gs *globalState, job libOrch.ChildJob, fs afero.Fs, workerInfo *libWorker.WorkerInfo) error {
	// Check if we have any stored data
	storedData := job["storedData"]
	if storedData == "" {
		return nil
	}

	var storedDataMap map[string]string
	err := json.Unmarshal([]byte(storedData), &storedDataMap)
	if err != nil {
		return err
	}

	// Check if we have any stored data
	if len(storedDataMap) == 0 {
		return nil
	}

	for key, value := range storedDataMap {
		// Write the file to the filesystem
		filePath := filepath.Join(gs.workDir, key)
		err := afero.WriteFile(fs, filePath, []byte(value), 0644)
		if err != nil {
			return err
		}

		// Add the file to the files map
		workerInfo.Files[key] = filePath
	}

	return nil

}
*/
