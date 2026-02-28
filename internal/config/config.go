package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	moodleerrors "github.com/nathanfredericks/moodle-cli/internal/errors"
)

const (
	configDirName   = ".moodle"
	configFileName  = "config.json"
	dirPermissions  = 0700
	filePermissions = 0600
)

// Site represents the configured Moodle instance.
type Site struct {
	URL      string `json:"url"`
	Username string `json:"username,omitempty"`
}

// Config holds the CLI configuration.
type Config struct {
	Site     *Site             `json:"site,omitempty"`
	Settings map[string]string `json:"settings,omitempty"`
}

// ConfigManager interface for accessing configuration.
type ConfigManager interface {
	Site() (Site, error)
	SaveSite(site Site) error
	DeleteSite() error
	Get(key string) string
	Set(key string, value string) error
	ConfigDir() string
}

// FileConfigManager implements ConfigManager using the filesystem.
type FileConfigManager struct {
	dir string
}

// NewConfigManager creates a new FileConfigManager.
func NewConfigManager() (*FileConfigManager, error) {
	dir, err := configDir()
	if err != nil {
		return nil, &moodleerrors.ConfigError{Msg: "unable to determine config directory", Err: err}
	}
	if err := os.MkdirAll(dir, dirPermissions); err != nil {
		return nil, &moodleerrors.ConfigError{Msg: "unable to create config directory", Err: err}
	}
	return &FileConfigManager{dir: dir}, nil
}

// NewConfigManagerWithDir creates a FileConfigManager using a specific directory.
func NewConfigManagerWithDir(dir string) (*FileConfigManager, error) {
	if err := os.MkdirAll(dir, dirPermissions); err != nil {
		return nil, &moodleerrors.ConfigError{Msg: "unable to create config directory", Err: err}
	}
	return &FileConfigManager{dir: dir}, nil
}

func (m *FileConfigManager) ConfigDir() string {
	return m.dir
}

func (m *FileConfigManager) Site() (Site, error) {
	// Check env vars first
	if url := os.Getenv("MOODLE_URL"); url != "" {
		return Site{URL: url}, nil
	}

	cfg, err := m.loadConfig()
	if err != nil {
		return Site{}, err
	}
	if cfg.Site == nil {
		return Site{}, &moodleerrors.ConfigError{Msg: "not logged in; run 'moodle auth login'"}
	}
	return *cfg.Site, nil
}

func (m *FileConfigManager) SaveSite(site Site) error {
	cfg, err := m.loadConfig()
	if err != nil {
		return err
	}
	cfg.Site = &site
	return m.saveConfig(cfg)
}

func (m *FileConfigManager) DeleteSite() error {
	cfg, err := m.loadConfig()
	if err != nil {
		return err
	}
	cfg.Site = nil
	return m.saveConfig(cfg)
}

func (m *FileConfigManager) Get(key string) string {
	cfg, err := m.loadConfig()
	if err != nil {
		return ""
	}
	return cfg.Settings[key]
}

// AllSettings returns all configuration settings.
func (m *FileConfigManager) AllSettings() map[string]string {
	cfg, err := m.loadConfig()
	if err != nil {
		return nil
	}
	return cfg.Settings
}

func (m *FileConfigManager) Set(key string, value string) error {
	cfg, err := m.loadConfig()
	if err != nil {
		return err
	}
	if cfg.Settings == nil {
		cfg.Settings = make(map[string]string)
	}
	cfg.Settings[key] = value
	return m.saveConfig(cfg)
}

func (m *FileConfigManager) loadConfig() (Config, error) {
	path := filepath.Join(m.dir, configFileName)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Config{Settings: make(map[string]string)}, nil
	}
	if err != nil {
		return Config{}, &moodleerrors.ConfigError{Msg: "unable to read config", Err: err}
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, &moodleerrors.ConfigError{Msg: "unable to parse config", Err: err}
	}
	if cfg.Settings == nil {
		cfg.Settings = make(map[string]string)
	}
	return cfg, nil
}

func (m *FileConfigManager) saveConfig(cfg Config) error {
	path := filepath.Join(m.dir, configFileName)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return &moodleerrors.ConfigError{Msg: "unable to marshal config", Err: err}
	}
	return os.WriteFile(path, data, filePermissions)
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDirName), nil
}
