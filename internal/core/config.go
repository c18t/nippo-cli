package core

import (
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	configDir string
	dataDir   string
	cacheDir  string

	LastUpdateCheckTimestamp time.Time `mapstructure:"last_update_check_timestamp"`
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

	// set default value
	time, _ := time.Parse(time.RFC3339, "2001-01-01T00:00:00Z")
	viper.SetDefault("last_update_check_timestamp", time)

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			viper.SafeWriteConfig()
		default:
			return err
		}
	}
	return viper.Unmarshal(c)
}

func (c *Config) SaveConfig() error {
	cMaps := c.configFieldMap()
	vMaps := viper.AllSettings()
	for key := range vMaps {
		viper.Set(key, cMaps[key])
	}
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}
	return home
}

func (c *Config) configFieldMap() map[string]any {
	var cMap = map[string]any{}
	t := reflect.TypeOf(*c)
	v := reflect.ValueOf(*c)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("mapstructure")
		if tag == "" {
			continue
		}
		cMap[tag] = v.Field(i).Interface()
	}
	return cMap
}
