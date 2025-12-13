package repofile

func (fr *FileRecordsRepo) Reset() {
	if fr == nil {
		return
	}

	fr.resetTest = 0

}
