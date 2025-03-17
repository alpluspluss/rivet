package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (bc *BuildCache) Save() error {
	for sourcePath, entry := range bc.entries {
		cacheFileName := fmt.Sprintf("%s.cache", filepath.Base(sourcePath))
		cacheFilePath := filepath.Join(bc.cacheDir, cacheFileName)

		jsonData, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("failed to serialize cache: %v", err)
		}

		if err := os.WriteFile(cacheFilePath, jsonData, 0644); err != nil {
			return fmt.Errorf("failed to write cache file: %v", err)
		}
	}
	return nil
}

func (bc *BuildCache) Load() error {
	entries, err := os.ReadDir(bc.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %v", err)
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".cache" {
			cacheFilePath := filepath.Join(bc.cacheDir, entry.Name())

			jsonData, err := os.ReadFile(cacheFilePath)
			if err != nil {
				return fmt.Errorf("failed to read cache file: %v", err)
			}

			var cacheEntry CacheEntry
			if err := json.Unmarshal(jsonData, &cacheEntry); err != nil {
				return fmt.Errorf("failed to parse cache: %v", err)
			}

			sourceName := filepath.Base(filepath.Base(entry.Name()))
			bc.entries[sourceName] = cacheEntry
		}
	}
	return nil
}

func (bc *BuildCache) Clean() error {
	if _, err := os.Stat(bc.cacheDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(bc.cacheDir); err != nil {
			return fmt.Errorf("failed to remove cache directory: %v", err)
		}
	}

	return os.MkdirAll(bc.cacheDir, 0755)
}
