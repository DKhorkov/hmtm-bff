package usecases

func validateFileExtension(extension string, allowedExtensions []string) bool {
	for _, allowedExtension := range allowedExtensions {
		if allowedExtension == extension {
			return true
		}
	}

	return false
}

func validateFileSize(size int64, maxSize int64) bool {
	return size <= maxSize
}