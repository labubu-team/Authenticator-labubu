package app

import "os"

func SaveToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

func ReadFromFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
