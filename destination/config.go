// Copyright © 2022 Meroxa, Inc. and Miquido
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package destination

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/miquido/conduit-connector-elasticsearch/internal/elasticsearch"
)

const (
	ConfigKeyVersion       = "version"
	ConfigKeyConnectionURI = "connectionUri"
	ConfigKeyBulkSize      = "bulkSize"
)

type Config struct {
	Version                elasticsearch.Version
	Host                   string
	Username               string
	Password               string
	CloudID                string
	APIKey                 string
	ServiceToken           string
	CertificateFingerprint string
	Index                  string
	Type                   string
	BulkSize               uint64
}

func (c Config) GetHost() string {
	return c.Host
}

func (c Config) GetUsername() string {
	return c.Username
}

func (c Config) GetPassword() string {
	return c.Password
}

func (c Config) GetCloudID() string {
	return c.CloudID
}

func (c Config) GetAPIKey() string {
	return c.APIKey
}

func (c Config) GetServiceToken() string {
	return c.ServiceToken
}

func (c Config) GetCertificateFingerprint() string {
	return c.CertificateFingerprint
}

func (c Config) GetIndex() string {
	return c.Index
}

func (c Config) GetType() string {
	return c.Type
}

func ParseConfig(cfgRaw map[string]string) (cfg Config, err error) {
	cfg.Version = cfgRaw[ConfigKeyVersion]

	// Version
	if cfg.Version == "" {
		return Config{}, requiredConfigErr(ConfigKeyVersion)
	}
	if cfg.Version != elasticsearch.Version5 &&
		cfg.Version != elasticsearch.Version6 &&
		cfg.Version != elasticsearch.Version7 &&
		cfg.Version != elasticsearch.Version8 {
		return Config{}, fmt.Errorf(
			"%q config value must be one of [%s], %s provided",
			ConfigKeyVersion,
			strings.Join([]elasticsearch.Version{
				elasticsearch.Version5,
				elasticsearch.Version6,
				elasticsearch.Version7,
				elasticsearch.Version8,
			}, ", "),
			cfg.Version,
		)
	}

	// Connection URI
	cfg.Host, cfg.Username, cfg.Password, cfg.Index, cfg.Type, cfg.CloudID, cfg.APIKey, cfg.ServiceToken, cfg.CertificateFingerprint, err = parseConnectionURIConfigValue(cfg, cfgRaw)
	if err != nil {
		return Config{}, err
	}

	// Bulk size
	cfg.BulkSize, err = parseBulkSizeConfigValue(cfgRaw)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func requiredConfigErr(name string) error {
	return fmt.Errorf("%q config value must be set", name)
}

func parseConnectionURIConfigValue(
	cfg Config,
	cfgRaw map[string]string,
) (
	host, username, password, indexName, indexType, cloudID, apiKey, serviceToken, certificateFingerprint string,
	_ error,
) {
	connectionURI, exists := cfgRaw[ConfigKeyConnectionURI]
	if !exists || connectionURI == "" {
		return "", "", "", "", "", "", "", "", "", requiredConfigErr(ConfigKeyConnectionURI)
	}

	connectionURIParsed, err := url.Parse(connectionURI)
	if err != nil {
		return "", "", "", "", "", "", "", "", "", fmt.Errorf("%q config value is invalid: %w", ConfigKeyConnectionURI, err)
	}

	if connectionURIParsed.Hostname() == "" {
		return "", "", "", "", "", "", "", "", "", fmt.Errorf("%q config value is invalid: host is required", ConfigKeyConnectionURI)
	}

	host = fmt.Sprintf("%s://%s", connectionURIParsed.Scheme, connectionURIParsed.Host)

	if connectionURIParsed.Scheme != "http" && connectionURIParsed.Scheme != "https" {
		return "", "", "", "", "", "", "", "", "", fmt.Errorf("%q config value is invalid: URI scheme needs to be one of [http, https], %q provided", ConfigKeyConnectionURI, connectionURIParsed.Scheme)
	}

	if user := connectionURIParsed.User; user != nil {
		username = user.Username()
		password, _ = user.Password()
	}

	index := strings.SplitN(strings.TrimLeft(connectionURIParsed.Path, "/"), "/", 3)
	if cfg.Version == elasticsearch.Version5 || cfg.Version == elasticsearch.Version6 {
		if len(index) < 2 || index[1] == "" {
			return "", "", "", "", "", "", "", "", "", fmt.Errorf("%q config value is invalid: index type needs to be provided in the path", ConfigKeyConnectionURI)
		}

		indexName = index[0]
		indexType = index[1]
	} else {
		indexName = index[0]
	}

	if indexName == "" {
		return "", "", "", "", "", "", "", "", "", fmt.Errorf("%q config value is invalid: index name needs to be provided in the path", ConfigKeyConnectionURI)
	}

	query := connectionURIParsed.Query()

	cloudID = query.Get("cloud_id")
	apiKey = query.Get("api_key")
	serviceToken = query.Get("service_token")
	certificateFingerprint = query.Get("certificate_fingerprint")

	return host, username, password, indexName, indexType, cloudID, apiKey, serviceToken, certificateFingerprint, nil
}

func parseBulkSizeConfigValue(cfgRaw map[string]string) (uint64, error) {
	bulkSize, ok := cfgRaw[ConfigKeyBulkSize]
	if !ok {
		return 0, requiredConfigErr(ConfigKeyBulkSize)
	}

	bulkSizeParsed, err := strconv.ParseUint(bulkSize, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %q config value: %w", ConfigKeyBulkSize, err)
	}
	if bulkSizeParsed <= 0 {
		return 0, fmt.Errorf("failed to parse %q config value: value must be greater than 0", ConfigKeyBulkSize)
	}
	if bulkSizeParsed > 10_000 {
		return 0, fmt.Errorf("failed to parse %q config value: value must be less than 10 000", ConfigKeyBulkSize)
	}

	return bulkSizeParsed, nil
}
