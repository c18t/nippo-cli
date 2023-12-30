package core

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	configDir string
	dataDir   string
	cacheDir  string
}

var Cfg *Config

func (c *Config) GetConfigDir() string {
	if c.configDir != "" {
		return c.configDir
	}

	defaultConfigDir := filepath.Join(c.homeDir(), ".config")
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" || !filepath.IsAbs(configDir) {
		configDir = defaultConfigDir
	}
	c.configDir = filepath.Join(configDir, "nippo")
	return c.configDir
}

func (c *Config) GetDataDir() string {
	if c.dataDir != "" {
		return c.dataDir
	}

	defaultDataDir := filepath.Join(c.homeDir(), ".local", "share")
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" || !filepath.IsAbs(dataDir) {
		dataDir = defaultDataDir
	}
	c.dataDir = filepath.Join(dataDir, "nippo")
	return c.dataDir
}

func (c *Config) GetCacheDir() string {
	if c.cacheDir != "" {
		return c.cacheDir
	}

	defaultCacheDir := filepath.Join(c.homeDir(), ".cache")
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" || !filepath.IsAbs(cacheDir) {
		cacheDir = defaultCacheDir
	}
	c.cacheDir = filepath.Join(cacheDir, "nippo")
	return c.cacheDir
}

func (c *Config) LoadConfig(filePath string) error {
	if filePath != "" {
		viper.SetConfigFile(filePath)
	} else {
		viper.AddConfigPath(c.GetDataDir())
		viper.SetConfigType("toml")
		viper.SetConfigName("nippo")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return viper.Unmarshal(&c)
}

func (c *Config) homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}
	return home
}
