package utils

// contains checks if a given dataSrc string is in the dataSrcRepeat.
// Return true if dataSrc is in dataSrcRepeat, else false.
func Contains(dataSrcRepeat []string, dataSrc string) bool {
	for _, data := range dataSrcRepeat {
		if data == dataSrc {
			return true
		}
	}
	return false
}
