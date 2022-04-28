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

	bulkSizeParsed, err := parseBulkSizeConfigValue(cfgRaw)
	if err != nil {
		return Config{}, err
	}

	cfg.BulkSize = bulkSizeParsed

	return cfg, nil
}

func requiredConfigErr(name string) error {
	return fmt.Errorf("%q config value must be set", name)
}

func parseBulkSizeConfigValue(cfgRaw map[string]string) (uint64, error) {
	bulkSize, ok := cfgRaw[ConfigKeyBulkSize]
	if !ok {
		return 0, requiredConfigErr(ConfigKeyBulkSize)
	}

	bulkSizeParsed, err := strconv.ParseUint(bulkSize, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse bulkSize config value: %w", err)
	}
	if bulkSizeParsed <= 0 {
		return 0, fmt.Errorf("failed to parse bulkSize config value: value must be greated than 0")
	}

	return bulkSizeParsed, nil
}
