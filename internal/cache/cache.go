package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CacheEntry represents a cached analysis result
type CacheEntry struct {
	RepoName   string          `json:"repo_name"`
	CachedAt   time.Time       `json:"cached_at"`
	ExpiresAt  time.Time       `json:"expires_at"`
	Analysis   json.RawMessage `json:"analysis"`
}

// CacheIndex stores metadata about all cached repos
type CacheIndex struct {
	Entries     map[string]CacheIndexEntry `json:"entries"`
	LastUpdated time.Time                  `json:"last_updated"`
}

// CacheIndexEntry is a summary of a cached repo
type CacheIndexEntry struct {
	RepoName  string    `json:"repo_name"`
	CachedAt  time.Time `json:"cached_at"`
	ExpiresAt time.Time `json:"expires_at"`
	FileSize  int64     `json:"file_size"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled     bool          `json:"enabled"`
	TTL         time.Duration `json:"ttl"`
	MaxSize     int64         `json:"max_size_mb"`
	AutoCache   bool          `json:"auto_cache"`
}

// Cache manages the local cache for analysis results
type Cache struct {
	cacheDir string
	config   CacheConfig
	index    *CacheIndex
}

// DefaultConfig returns the default cache configuration
func DefaultConfig() CacheConfig {
	return CacheConfig{
		Enabled:   true,
		TTL:       24 * time.Hour, // 24 hours default
		MaxSize:   100,            // 100 MB max
		AutoCache: true,
	}
}

// NewCache creates a new cache instance
func NewCache() (*Cache, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, err
	}

	// Ensure cache directories exist
	reposDir := filepath.Join(cacheDir, "repos")
	if err := os.MkdirAll(reposDir, 0755); err != nil {
		return nil, err
	}

	cache := &Cache{
		cacheDir: cacheDir,
		config:   DefaultConfig(),
	}

	// Load or create index
	if err := cache.loadIndex(); err != nil {
		cache.index = &CacheIndex{
			Entries:     make(map[string]CacheIndexEntry),
			LastUpdated: time.Now(),
		}
	}

	// Load config if exists
	cache.loadConfig()

	return cache, nil
}

// getCacheDir returns the cache directory path
func getCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".repo-lyzer", "cache"), nil
}

// repoToFilename converts repo name to safe filename
func repoToFilename(repoName string) string {
	return strings.ReplaceAll(repoName, "/", "_") + ".json"
}

// loadIndex loads the cache index from disk
func (c *Cache) loadIndex() error {
	indexPath := filepath.Join(c.cacheDir, "cache_index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	c.index = &CacheIndex{}
	return json.Unmarshal(data, c.index)
}

// saveIndex saves the cache index to disk
func (c *Cache) saveIndex() error {
	c.index.LastUpdated = time.Now()
	indexPath := filepath.Join(c.cacheDir, "cache_index.json")
	
	data, err := json.MarshalIndent(c.index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(indexPath, data, 0644)
}

// loadConfig loads cache configuration
func (c *Cache) loadConfig() {
	configPath := filepath.Join(c.cacheDir, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return // Use defaults
	}
	json.Unmarshal(data, &c.config)
}

// SaveConfig saves cache configuration
func (c *Cache) SaveConfig() error {
	configPath := filepath.Join(c.cacheDir, "config.json")
	data, err := json.MarshalIndent(c.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// Get retrieves a cached analysis if available and not expired
func (c *Cache) Get(repoName string) (*CacheEntry, bool) {
	if !c.config.Enabled {
		return nil, false
	}

	entry, exists := c.index.Entries[repoName]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	// Load full entry from file
	filename := repoToFilename(repoName)
	filePath := filepath.Join(c.cacheDir, "repos", filename)
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false
	}

	var cacheEntry CacheEntry
	if err := json.Unmarshal(data, &cacheEntry); err != nil {
		return nil, false
	}

	return &cacheEntry, true
}

// Set stores an analysis result in the cache
func (c *Cache) Set(repoName string, analysis interface{}) error {
	if !c.config.Enabled || !c.config.AutoCache {
		return nil
	}

	// Serialize analysis
	analysisData, err := json.Marshal(analysis)
	if err != nil {
		return err
	}

	now := time.Now()
	entry := CacheEntry{
		RepoName:  repoName,
		CachedAt:  now,
		ExpiresAt: now.Add(c.config.TTL),
		Analysis:  analysisData,
	}

	// Save to file
	filename := repoToFilename(repoName)
	filePath := filepath.Join(c.cacheDir, "repos", filename)
	
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	// Update index
	fileInfo, _ := os.Stat(filePath)
	fileSize := int64(0)
	if fileInfo != nil {
		fileSize = fileInfo.Size()
	}

	c.index.Entries[repoName] = CacheIndexEntry{
		RepoName:  repoName,
		CachedAt:  now,
		ExpiresAt: now.Add(c.config.TTL),
		FileSize:  fileSize,
	}

	return c.saveIndex()
}

// Delete removes a specific repo from cache
func (c *Cache) Delete(repoName string) error {
	filename := repoToFilename(repoName)
	filePath := filepath.Join(c.cacheDir, "repos", filename)
	
	os.Remove(filePath) // Ignore error if file doesn't exist
	
	delete(c.index.Entries, repoName)
	return c.saveIndex()
}

// Clear removes all cached data
func (c *Cache) Clear() error {
	reposDir := filepath.Join(c.cacheDir, "repos")
	
	// Remove all files in repos directory
	entries, err := os.ReadDir(reposDir)
	if err == nil {
		for _, entry := range entries {
			os.Remove(filepath.Join(reposDir, entry.Name()))
		}
	}

	// Clear index
	c.index.Entries = make(map[string]CacheIndexEntry)
	return c.saveIndex()
}

// GetStats returns cache statistics
func (c *Cache) GetStats() CacheStats {
	var totalSize int64
	validCount := 0
	expiredCount := 0
	now := time.Now()

	for _, entry := range c.index.Entries {
		totalSize += entry.FileSize
		if now.After(entry.ExpiresAt) {
			expiredCount++
		} else {
			validCount++
		}
	}

	return CacheStats{
		TotalRepos:   len(c.index.Entries),
		ValidRepos:   validCount,
		ExpiredRepos: expiredCount,
		TotalSizeMB:  float64(totalSize) / (1024 * 1024),
		CacheDir:     c.cacheDir,
	}
}

// CacheStats holds cache statistics
type CacheStats struct {
	TotalRepos   int
	ValidRepos   int
	ExpiredRepos int
	TotalSizeMB  float64
	CacheDir     string
}

// GetCachedRepos returns list of all cached repos
func (c *Cache) GetCachedRepos() []CacheIndexEntry {
	repos := make([]CacheIndexEntry, 0, len(c.index.Entries))
	for _, entry := range c.index.Entries {
		repos = append(repos, entry)
	}
	return repos
}

// IsExpired checks if a cached entry is expired
func (c *Cache) IsExpired(repoName string) bool {
	entry, exists := c.index.Entries[repoName]
	if !exists {
		return true
	}
	return time.Now().After(entry.ExpiresAt)
}

// HasCache checks if a repo is in cache (expired or not)
func (c *Cache) HasCache(repoName string) bool {
	_, exists := c.index.Entries[repoName]
	return exists
}

// GetConfig returns current cache configuration
func (c *Cache) GetConfig() CacheConfig {
	return c.config
}

// SetEnabled enables or disables caching
func (c *Cache) SetEnabled(enabled bool) {
	c.config.Enabled = enabled
	c.SaveConfig()
}

// SetTTL sets the cache time-to-live
func (c *Cache) SetTTL(ttl time.Duration) {
	c.config.TTL = ttl
	c.SaveConfig()
}

// SetAutoCache enables or disables auto-caching
func (c *Cache) SetAutoCache(enabled bool) {
	c.config.AutoCache = enabled
	c.SaveConfig()
}

// CleanExpired removes all expired entries
func (c *Cache) CleanExpired() int {
	removed := 0
	now := time.Now()
	
	for repoName, entry := range c.index.Entries {
		if now.After(entry.ExpiresAt) {
			c.Delete(repoName)
			removed++
		}
	}
	
	return removed
}

// FormatTTL returns a human-readable TTL string
func FormatTTL(d time.Duration) string {
	if d >= 24*time.Hour {
		days := int(d / (24 * time.Hour))
		if days == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%d days", days)
	}
	if d >= time.Hour {
		hours := int(d / time.Hour)
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}
	minutes := int(d / time.Minute)
	if minutes == 1 {
		return "1 minute"
	}
	return fmt.Sprintf("%d minutes", minutes)
}
