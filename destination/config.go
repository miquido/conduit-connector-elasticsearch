package destination

import (
	"fmt"
	"strconv"
)

const (
	ConfigKeyHost     = "host"
	ConfigKeyUsername = "username"
	ConfigKeyPassword = "password"
	ConfigKeyIndex    = "index"
	ConfigKeyBulkSize = "bulkSize"
)

type Config struct {
	Host     string
	Username string
	Password string
	Index    string
	BulkSize uint64
}

func ParseConfig(cfgRaw map[string]string) (Config, error) {
	cfg := Config{
		Host:     cfgRaw[ConfigKeyHost],
		Username: cfgRaw[ConfigKeyUsername],
		Password: cfgRaw[ConfigKeyPassword],
		Index:    cfgRaw[ConfigKeyIndex],
	}

	if cfg.Host == "" {
		return Config{}, requiredConfigErr(ConfigKeyHost)
	}
	if cfg.Password == "" && cfg.Username != "" {
		return Config{}, fmt.Errorf("%q config value must be set when %q is provided", ConfigKeyPassword, ConfigKeyUsername)
	}
	if cfg.Index == "" {
		return Config{}, requiredConfigErr(ConfigKeyIndex)
	}

	if bulkSize, ok := cfgRaw[ConfigKeyBulkSize]; !ok {
		return Config{}, requiredConfigErr(ConfigKeyBulkSize)
	} else if bulkSizeParsed, err := strconv.ParseUint(bulkSize, 10, 64); err != nil {
		return Config{}, fmt.Errorf("failed to parse bulkSize config value: %w", err)
	} else if bulkSizeParsed <= 0 {
		return Config{}, fmt.Errorf("failed to parse bulkSize config value: value must be greated than 0")
	} else {
		cfg.BulkSize = bulkSizeParsed
	}

	return cfg, nil
}

func requiredConfigErr(name string) error {
	return fmt.Errorf("%q config value must be set", name)
}
