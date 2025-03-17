package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/alpluspluss/rivet/internal/utils"
)

type CacheEntry struct {
	Hash          string              `json:"hash"`
	Includes      map[string]FileInfo `json:"includes"`
	CompilerFlags []string            `json:"compiler_flags"`
	Target        string              `json:"target"`
	Profile       string              `json:"profile"`
	Timestamp     int64               `json:"timestamp"`
}

type FileInfo struct {
	Hash    string `json:"hash"`
	ModTime int64  `json:"mtime"`
	Size    int64  `json:"size"`
}

type BuildCache struct {
	cacheDir   string
	entries    map[string]CacheEntry
	quickCheck bool
}

func NewBuildCache(workspaceRoot string) *BuildCache {
	cacheDir := filepath.Join(workspaceRoot, ".rivet_cache")
	os.MkdirAll(cacheDir, 0755)

	return &BuildCache{
		cacheDir:   cacheDir,
		entries:    make(map[string]CacheEntry),
		quickCheck: true,
	}
}

func (bc *BuildCache) NeedsRebuild(
	source, object string,
	includes []string,
	compilerFlags []string,
	target, profile string,
) bool {
	log.Printf("Checking if %s needs rebuild...", source)
	if _, err := os.Stat(object); os.IsNotExist(err) {
		log.Println("Object file doesn't exist")
		return true
	}

	entry, exists := bc.entries[source]
	if !exists {
		log.Println("No cache entry found")
		return true
	}

	if entry.Target != target ||
		entry.Profile != profile ||
		!utils.StringSlicesEqual(entry.CompilerFlags, compilerFlags) {
		log.Println("Build configuration changed")
		return true
	}

	if bc.fileChanged(source, entry.Hash) {
		log.Println("Source file changed")
		return true
	}

	for _, include := range includes {
		includeInfo, exists := entry.Includes[include]
		if !exists {
			log.Printf("New include file %s", include)
			return true
		}

		if bc.fileChangedWithInfo(include, &includeInfo) {
			log.Printf("Include file %s changed", include)
			return true
		}
	}

	if len(entry.Includes) != len(includes) {
		log.Println("Number of includes changed")
		return true
	}

	return false
}

func (bc *BuildCache) Update(source string, includes []string,
	compilerFlags []string,
	target, profile string,
) error {
	includeInfos := make(map[string]FileInfo)
	for _, include := range includes {
		info, err := bc.getFileInfo(include)
		if err != nil {
			return fmt.Errorf("failed to get include file info: %v", err)
		}
		includeInfos[include] = *info
	}

	sourceInfo, err := bc.getFileInfo(source)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %v", err)
	}

	bc.entries[source] = CacheEntry{
		Hash:          sourceInfo.Hash,
		Includes:      includeInfos,
		CompilerFlags: compilerFlags,
		Target:        target,
		Profile:       profile,
		Timestamp:     time.Now().Unix(),
	}

	return nil
}

func (bc *BuildCache) getFileInfo(path string) (*FileInfo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for %s: %v", path, err)
	}

	var hash string
	if bc.quickCheck {
		hash = "quick_check"
	} else {
		fileHash, err := bc.hashFile(path)
		if err != nil {
			return nil, err
		}
		hash = fileHash
	}

	return &FileInfo{
		Hash:    hash,
		ModTime: stat.ModTime().Unix(),
		Size:    stat.Size(),
	}, nil
}

func (bc *BuildCache) fileChanged(path string, oldHash string) bool {
	info, err := bc.getFileInfo(path)
	if err != nil {
		return true
	}

	if bc.quickCheck {
		return false
	}

	return info.Hash != oldHash
}

func (bc *BuildCache) fileChangedWithInfo(path string, oldInfo *FileInfo) bool {
	newInfo, err := bc.getFileInfo(path)
	if err != nil {
		return true
	}

	if bc.quickCheck {
		return newInfo.ModTime != oldInfo.ModTime || newInfo.Size != oldInfo.Size
	}

	return newInfo.Hash != oldInfo.Hash
}

func (bc *BuildCache) hashFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %v", path, err)
	}

	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:]), nil
}

func (bc *BuildCache) SetQuickCheck(enable bool) {
	bc.quickCheck = enable
}
