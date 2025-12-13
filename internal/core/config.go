package core

import (
	"encoding"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var Cfg *Config

// ErrConfigNotFound is returned when the configuration file does not exist.
// This indicates the user should run 'nippo init' first.
type ErrConfigNotFound struct {
	Path string
}

func (e *ErrConfigNotFound) Error() string {
	return fmt.Sprintf("configuration file not found: %s", e.Path)
}

type Config struct {
	configDir string // resolved config directory (not persisted)
	dataDir   string // cached resolved data directory
	cacheDir  string // cached resolved cache directory

	LastUpdateCheckTimestamp time.Time     `mapstructure:"last_update_check_timestamp"`
	LastFormatTimestamp      time.Time     `mapstructure:"last_format_timestamp"`
	Project                  ConfigProject `mapstructure:"project"`
	Paths                    ConfigPaths   `mapstructure:"path"`
}

type ConfigProject struct {
	Url           string `mapstructure:"url"`
	DriveFolderId string `mapstructure:"drive_folder_id"`
	SiteUrl       string `mapstructure:"site_url"`
	Branch        string `mapstructure:"branch"`
	TemplatePath  string `mapstructure:"template_path"`
	AssetPath     string `mapstructure:"asset_path"`
}

type ConfigPaths struct {
	DataDir  string `mapstructure:"data_dir"`
	CacheDir string `mapstructure:"cache_dir"`
}

// InitConfig initializes the global configuration.
// If the config file is not found, Cfg is initialized with defaults
// and ErrConfigNotFound is returned. The caller can check for this
// error type to decide whether to proceed without a config file.
func InitConfig(configFile string) error {
	cfg := &Config{}
	err := cfg.LoadConfig(configFile)
	if err != nil {
		// Always set Cfg so commands can access default paths
		Cfg = cfg
		return err
	}
	Cfg = cfg
	return nil
}

var (
	timeType          = reflect.TypeOf((*time.Time)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func (c *Config) GetConfigDir() string {
	if c.configDir != "" {
		return c.configDir
	}
	c.configDir = ResolveConfigDir()
	return c.configDir
}

// GetConfigFilePath returns the path to the config file.
// If viper has a config file set, returns that. Otherwise returns default path.
func GetConfigFilePath() string {
	if cfgFile := viper.ConfigFileUsed(); cfgFile != "" {
		return cfgFile
	}
	return filepath.Join(ResolveConfigDir(), "nippo.toml")
}

func (c *Config) GetDataDir() string {
	if c.dataDir != "" {
		return c.dataDir
	}

	// Priority: configured path > fallback chain
	if c.Paths.DataDir != "" {
		expanded := ExpandPath(c.Paths.DataDir)
		if filepath.IsAbs(expanded) {
			c.dataDir = expanded
		} else {
			// Relative paths are resolved relative to config directory
			c.dataDir = filepath.Join(c.GetConfigDir(), expanded)
		}
		return c.dataDir
	}

	c.dataDir = ResolveDataDir()
	return c.dataDir
}

func (c *Config) GetCacheDir() string {
	if c.cacheDir != "" {
		return c.cacheDir
	}

	// Priority: configured path > fallback chain
	if c.Paths.CacheDir != "" {
		expanded := ExpandPath(c.Paths.CacheDir)
		if filepath.IsAbs(expanded) {
			c.cacheDir = expanded
		} else {
			// Relative paths are resolved relative to config directory
			c.cacheDir = filepath.Join(c.GetConfigDir(), expanded)
		}
		return c.cacheDir
	}

	c.cacheDir = ResolveCacheDir()
	return c.cacheDir
}

func (c *Config) ResetLastUpdateCheckTimestamp() {
	c.LastUpdateCheckTimestamp = c.getDefaultLastUpdateCheckTimestamp()
}

func (c *Config) LoadConfig(filePath string) error {
	if filePath != "" {
		viper.SetConfigFile(filePath)
		// Set configDir to the directory containing the specified config file
		absPath, err := filepath.Abs(filePath)
		if err == nil {
			c.configDir = filepath.Dir(absPath)
		}
	} else {
		viper.AddConfigPath(c.GetConfigDir())
		viper.SetConfigType("toml")
		viper.SetConfigName("nippo")
	}

	// set default value
	viper.SetDefault("last_update_check_timestamp", c.getDefaultLastUpdateCheckTimestamp())

	viper.SetEnvPrefix("NIPPO")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// Do not auto-create config file. Return ErrConfigNotFound.
			configPath := filepath.Join(c.GetConfigDir(), "nippo.toml")
			if filePath != "" {
				configPath = filePath
			}
			return &ErrConfigNotFound{Path: configPath}
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

	// Check if paths are configured (non-empty)
	pathsConfigured := c.Paths.DataDir != "" || c.Paths.CacheDir != ""

	for key := range cMaps {
		// Skip empty path values - we'll add them as comments later
		if !pathsConfigured && (key == "path.data_dir" || key == "path.cache_dir") {
			continue
		}
		viper.Set(key, cMaps[key])
	}
	// Use WriteConfigAs to handle both new and existing config files
	configPath := GetConfigFilePath()

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	err = viper.WriteConfigAs(configPath)
	if err != nil {
		return err
	}

	// If paths are not configured, append commented [path] section
	if !pathsConfigured {
		err = c.appendCommentedPathsSection()
		if err != nil {
			return err
		}
	}

	return nil
}

// appendCommentedPathsSection adds a commented [path] section to the config file
func (c *Config) appendCommentedPathsSection() error {
	configPath := GetConfigFilePath()

	// Read current content
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Check if [path] section already exists
	if strings.Contains(string(content), "[path]") {
		return nil
	}

	// Append commented path section
	pathsSection := fmt.Sprintf(`
[path]
# Uncomment and modify to customize file locations.
# data_dir = %q
# cache_dir = %q
`, ResolveDataDir(), ResolveCacheDir())

	return os.WriteFile(configPath, append(content, []byte(pathsSection)...), 0644)
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
