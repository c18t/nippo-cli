package core

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	homeDir   string
	configDir string
	dataDir   string
	cacheDir  string
}

var Cfg *Config

func (c *Config) GetConfigDir() string {
	if c.configDir != "" {
		return c.configDir
	}

	defaultConfigDir := path.Join(c.homeDir, ".config")
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" || !path.IsAbs(configDir) {
		configDir = defaultConfigDir
	}
	c.configDir = path.Join(configDir, "nippo")
	return c.configDir
}

func (c *Config) GetDataDir() string {
	if c.dataDir != "" {
		return c.dataDir
	}

	defaultDataDir := path.Join(c.homeDir, ".local", "share")
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" || !path.IsAbs(dataDir) {
		dataDir = defaultDataDir
	}
	c.dataDir = path.Join(dataDir, "nippo")
	return c.dataDir
}

func (c *Config) GetCacheDir() string {
	if c.cacheDir != "" {
		return c.cacheDir
	}

	defaultCacheDir := path.Join(c.homeDir, ".cache")
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" || !path.IsAbs(cacheDir) {
		cacheDir = defaultCacheDir
	}
	c.cacheDir = path.Join(cacheDir, "nippo")
	return c.cacheDir
}

func (c *Config) LoadConfig(filePath string) error {
	if filePath != "" {
		viper.SetConfigFile(filePath)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			home = ""
		}
		c.homeDir = home
		viper.AddConfigPath(c.GetDataDir())
		viper.SetConfigType("toml")
		viper.SetConfigName(".nippo")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return viper.Unmarshal(&c)
}
