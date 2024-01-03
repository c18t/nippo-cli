package core

import (
	"encoding"
	"fmt"
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

	LastUpdateCheckTimestamp time.Time     `mapstructure:"last_update_check_timestamp"`
	Project                  ConfigProject `mapstructure:"project"`
}
type ConfigProject struct {
	Url          string `mapstructure:"url"`
	TemplatePath string `mapstructure:"template_path"`
	AssetPath    string `mapstructure:"asset_path"`
}

var Cfg *Config

var timeType = reflect.TypeOf((*time.Time)(nil)).Elem()
var textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

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

func (c *Config) ResetLastUpdateCheckTimestamp() {
	c.LastUpdateCheckTimestamp = c.getDefaultLastUpdateCheckTimestamp()
}

func (c *Config) LoadConfig(filePath string) error {
	if filePath != "" {
		viper.SetConfigFile(filePath)
	} else {
		viper.AddConfigPath(c.GetConfigDir())
		viper.SetConfigType("toml")
		viper.SetConfigName("nippo")
	}

	// set default value
	viper.SetDefault("last_update_check_timestamp", c.getDefaultLastUpdateCheckTimestamp())

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
	cMaps, err := c.configFieldMap(map[string]any{}, *c, "")
	if err != nil {
		return err
	}
	for key := range cMaps {
		viper.Set(key, cMaps[key])
	}
	err = viper.WriteConfig()
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

func (c *Config) configFieldMap(cMap map[string]any, i interface{}, prefex string) (map[string]any, error) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	for i := 0; i < t.NumField(); i++ {
		var err error
		tag := t.Field(i).Tag.Get("mapstructure")
		if tag == "" {
			continue
		}
		tag = fmt.Sprintf("%s%s", prefex, tag)
		prefix := tag + "."
		isTime := t.Field(i).Type == timeType
		hasTextMarshaler := t.Field(i).Type.Implements(textMarshalerType)
		if isTime || hasTextMarshaler {
			cMap[tag] = v.Field(i).Interface()
		} else {
			switch v.Field(i).Kind() {
			case reflect.Interface:
				if v.IsNil() {
					return nil, fmt.Errorf("encoding a nil interface is not supported")
				}
				cMap, err = c.configFieldMap(cMap, v.Field(i).Elem(), prefix)
			case reflect.Ptr:
				var el reflect.Value
				if v.IsNil() {
					el = reflect.Zero(v.Type().Elem())
				} else {
					el = v.Field(i).Elem()
				}
				cMap, err = c.configFieldMap(cMap, el, prefix)
			case reflect.Struct:
				cMap, err = c.configFieldMap(cMap, v.Field(i).Interface(), prefix)
			default:
				cMap[tag] = v.Field(i).Interface()
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return cMap, nil
}

func (c *Config) getDefaultLastUpdateCheckTimestamp() time.Time {
	time, _ := time.Parse(time.RFC3339, "2001-01-01T00:00:00Z")
	return time
}
