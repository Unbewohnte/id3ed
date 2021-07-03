package util

// Returns found key (int) in provided map by value (string);
// If key does not exist in map - returns -1
func GetKey(mp map[int]string, givenValue string) int {
	for key, value := range mp {
		if value == givenValue {
			return key
		}
	}
	return -1
}
