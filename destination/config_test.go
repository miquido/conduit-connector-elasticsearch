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
	"testing"

	"github.com/jaswdr/faker"
	"github.com/miquido/conduit-connector-elasticsearch/internal/elasticsearch"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	fakerInstance := faker.New()

	t.Run("fails when Version is empty", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			"nonExistentKey": "value",
		})

		require.EqualError(t, err, fmt.Sprintf("%q config value must be set", ConfigKeyVersion))
	})

	t.Run("fails when Version is unsupported", func(t *testing.T) {
		var version = "invalid-version"

		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion: version,
			"nonExistentKey": "value",
		})

		require.EqualError(t, err, fmt.Sprintf(
			"%q config value must be one of [%s, %s, %s, %s], %s provided",
			ConfigKeyVersion,
			elasticsearch.Version5,
			elasticsearch.Version6,
			elasticsearch.Version7,
			elasticsearch.Version8,
			version,
		))
	})

	t.Run("fails when connection URI is empty", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:       elasticsearch.Version8,
			ConfigKeyConnectionURI: "",
			"nonExistentKey":       "value",
		})

		require.EqualError(t, err, fmt.Sprintf(`%q config value must be set`, ConfigKeyConnectionURI))
	})

	t.Run("fails when connection URI has invalid format", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:       elasticsearch.Version8,
			ConfigKeyConnectionURI: "https://Username:Password@example.com:port",
			"nonExistentKey":       "value",
		})

		require.EqualError(t, err, fmt.Sprintf(`%q config value is invalid: parse "https://Username:Password@example.com:port": invalid port ":port" after host`, ConfigKeyConnectionURI))
	})

	t.Run("fails when Host is empty", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:       elasticsearch.Version8,
			ConfigKeyConnectionURI: "https://Username:Password@:1234/index/type?cloud_id=A&api_key=B&service_token=C&certificate_fingerprint=D",
			"nonExistentKey":       "value",
		})

		require.EqualError(t, err, fmt.Sprintf("%q config value is invalid: host is required", ConfigKeyConnectionURI))
	})

	t.Run("fails when URI scheme is unsupported", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:       elasticsearch.Version8,
			ConfigKeyConnectionURI: "custom://Username:Password@example.com:1234/index/type?cloud_id=A&api_key=B&service_token=C&certificate_fingerprint=D",
			"nonExistentKey":       "value",
		})

		require.EqualError(t, err, fmt.Sprintf("%q config value is invalid: URI scheme needs to be one of [http, https], \"custom\" provided", ConfigKeyConnectionURI))
	})

	t.Run("fails when index name was not provided", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:       elasticsearch.Version8,
			ConfigKeyConnectionURI: "https://Username:Password@example.com:1234?cloud_id=A&api_key=B&service_token=C&certificate_fingerprint=D",
			"nonExistentKey":       "value",
		})

		require.EqualError(t, err, fmt.Sprintf("%q config value is invalid: index name needs to be provided in the path", ConfigKeyConnectionURI))
	})

	for _, version := range []elasticsearch.Version{
		elasticsearch.Version5,
		elasticsearch.Version6,
	} {
		t.Run(fmt.Sprintf("fails when index type was not provided for version=%s", version), func(t *testing.T) {
			_, err := ParseConfig(map[string]string{
				ConfigKeyVersion:       version,
				ConfigKeyConnectionURI: "https://Username:Password@example.com:1234/index?cloud_id=A&api_key=B&service_token=C&certificate_fingerprint=D",
				"nonExistentKey":       "value",
			})

			require.EqualError(t, err, fmt.Sprintf("%q config value is invalid: index type needs to be provided in the path", ConfigKeyConnectionURI))
		})
	}

	t.Run("fails when Bulk Size is empty", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:       elasticsearch.Version8,
			ConfigKeyConnectionURI: "https://example.com/index",
			"nonExistentKey":       "value",
		})

		require.EqualError(t, err, fmt.Sprintf("%q config value must be set", ConfigKeyBulkSize))
	})

	t.Run("fails when Bulk Size is an invalid positive integer", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:       elasticsearch.Version8,
			ConfigKeyConnectionURI: "https://example.com/index",
			ConfigKeyBulkSize:      "-1",
			"nonExistentKey":       "value",
		})

		require.EqualError(t, err, fmt.Sprintf(`failed to parse %q config value: strconv.ParseUint: parsing "-1": invalid syntax`, ConfigKeyBulkSize))
	})

	t.Run("fails when Bulk Size is less than 1", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:       elasticsearch.Version8,
			ConfigKeyConnectionURI: "https://example.com/index",
			ConfigKeyBulkSize:      "0",
			"nonExistentKey":       "value",
		})

		require.EqualError(t, err, fmt.Sprintf("failed to parse %q config value: value must be greater than 0", ConfigKeyBulkSize))
	})

	t.Run("fails when Bulk Size is greater than 10 000", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:       elasticsearch.Version8,
			ConfigKeyConnectionURI: "https://example.com/index",
			ConfigKeyBulkSize:      "10001",
			"nonExistentKey":       "value",
		})

		require.EqualError(t, err, fmt.Sprintf("failed to parse %q config value: value must be less than 10 000", ConfigKeyBulkSize))
	})

	t.Run("returns config when all required config values were provided", func(t *testing.T) {
		var (
			host  = "example.com"
			index = fakerInstance.Lorem().Word()
		)

		var cfgRaw = map[string]string{
			ConfigKeyVersion:       elasticsearch.Version8,
			ConfigKeyConnectionURI: fmt.Sprintf("https://%s/%s", host, index),
			ConfigKeyBulkSize:      "1",
			"nonExistentKey":       "value",
		}

		config, err := ParseConfig(cfgRaw)

		require.NoError(t, err)
		require.Equal(t, cfgRaw[ConfigKeyVersion], config.Version)
		require.Equal(t, fmt.Sprintf("https://%s", host), config.Host)
		require.Equal(t, index, config.Index)
		require.Equal(t, cfgRaw[ConfigKeyBulkSize], fmt.Sprintf("%d", config.BulkSize))
		require.Empty(t, "", config.Username)
		require.Empty(t, "", config.Password)
		require.Empty(t, "", config.Type)
		require.Empty(t, "", config.CloudID)
		require.Empty(t, "", config.APIKey)
		require.Empty(t, "", config.ServiceToken)
		require.Empty(t, "", config.CertificateFingerprint)
	})

	for _, version := range []elasticsearch.Version{
		elasticsearch.Version5,
		elasticsearch.Version6,
	} {
		t.Run(fmt.Sprintf("returns config when all config values were provided for version=%s", version), func(t *testing.T) {
			var (
				username               = fakerInstance.Internet().User()
				password               = fakerInstance.Hash().MD5()
				host                   = "example.com"
				port                   = fakerInstance.IntBetween(0, 65535)
				indexName              = fakerInstance.Lorem().Word()
				indexType              = fakerInstance.Lorem().Word()
				cloudID                = fakerInstance.RandomStringWithLength(32)
				apiKey                 = fakerInstance.RandomStringWithLength(32)
				serviceToken           = fakerInstance.RandomStringWithLength(32)
				certificateFingerprint = fakerInstance.Hash().SHA256()
			)

			var cfgRaw = map[string]string{
				ConfigKeyVersion: version,
				ConfigKeyConnectionURI: fmt.Sprintf(
					"https://%s:%s@%s:%d/%s/%s?cloud_id=%s&api_key=%s&service_token=%s&certificate_fingerprint=%s",
					username,
					password,
					host,
					port,
					indexName,
					indexType,
					cloudID,
					apiKey,
					serviceToken,
					certificateFingerprint,
				),
				ConfigKeyBulkSize: fmt.Sprintf("%d", fakerInstance.Int32Between(1, 10_000)),
				"nonExistentKey":  "value",
			}

			config, err := ParseConfig(cfgRaw)

			require.NoError(t, err)
			require.Equal(t, cfgRaw[ConfigKeyVersion], config.Version)
			require.Equal(t, fmt.Sprintf("https://%s:%d", host, port), config.Host)
			require.Equal(t, indexName, config.Index)
			require.Equal(t, indexType, config.Type)
			require.Equal(t, cfgRaw[ConfigKeyBulkSize], fmt.Sprintf("%d", config.BulkSize))
			require.Equal(t, username, config.Username)
			require.Equal(t, password, config.Password)
			require.Equal(t, cloudID, config.CloudID)
			require.Equal(t, apiKey, config.APIKey)
			require.Equal(t, serviceToken, config.ServiceToken)
			require.Equal(t, certificateFingerprint, config.CertificateFingerprint)
		})
	}

	for _, version := range []elasticsearch.Version{
		elasticsearch.Version7,
		elasticsearch.Version8,
	} {
		t.Run(fmt.Sprintf("returns config when all config values were provided for version=%s", version), func(t *testing.T) {
			var (
				username               = fakerInstance.Internet().User()
				password               = fakerInstance.Hash().MD5()
				host                   = "example.com"
				port                   = fakerInstance.IntBetween(0, 65535)
				indexName              = fakerInstance.Lorem().Word()
				cloudID                = fakerInstance.RandomStringWithLength(32)
				apiKey                 = fakerInstance.RandomStringWithLength(32)
				serviceToken           = fakerInstance.RandomStringWithLength(32)
				certificateFingerprint = fakerInstance.Hash().SHA256()
			)

			var cfgRaw = map[string]string{
				ConfigKeyVersion: version,
				ConfigKeyConnectionURI: fmt.Sprintf(
					"https://%s:%s@%s:%d/%s?cloud_id=%s&api_key=%s&service_token=%s&certificate_fingerprint=%s",
					username,
					password,
					host,
					port,
					indexName,
					cloudID,
					apiKey,
					serviceToken,
					certificateFingerprint,
				),
				ConfigKeyBulkSize: fmt.Sprintf("%d", fakerInstance.Int32Between(1, 10_000)),
				"nonExistentKey":  "value",
			}

			config, err := ParseConfig(cfgRaw)

			require.NoError(t, err)
			require.Equal(t, cfgRaw[ConfigKeyVersion], config.Version)
			require.Equal(t, fmt.Sprintf("https://%s:%d", host, port), config.Host)
			require.Equal(t, indexName, config.Index)
			require.Equal(t, "", config.Type)
			require.Equal(t, cfgRaw[ConfigKeyBulkSize], fmt.Sprintf("%d", config.BulkSize))
			require.Equal(t, username, config.Username)
			require.Equal(t, password, config.Password)
			require.Equal(t, cloudID, config.CloudID)
			require.Equal(t, apiKey, config.APIKey)
			require.Equal(t, serviceToken, config.ServiceToken)
			require.Equal(t, certificateFingerprint, config.CertificateFingerprint)
		})
	}
}

func TestConfig_Getters(t *testing.T) {
	fakerInstance := faker.New()

	var (
		host                   = fakerInstance.Internet().URL()
		username               = fakerInstance.Internet().Email()
		password               = fakerInstance.Internet().Password()
		cloudID                = fakerInstance.RandomStringWithLength(32)
		apiKey                 = fakerInstance.RandomStringWithLength(32)
		serviceToken           = fakerInstance.RandomStringWithLength(32)
		certificateFingerprint = fakerInstance.Hash().SHA256()
		indexName              = fakerInstance.Lorem().Word()
		indexType              = fakerInstance.Lorem().Word()
	)

	config := Config{
		Host:                   host,
		Username:               username,
		Password:               password,
		CloudID:                cloudID,
		APIKey:                 apiKey,
		ServiceToken:           serviceToken,
		CertificateFingerprint: certificateFingerprint,
		Index:                  indexName,
		Type:                   indexType,
	}

	require.Equal(t, host, config.GetHost())
	require.Equal(t, username, config.GetUsername())
	require.Equal(t, password, config.GetPassword())
	require.Equal(t, cloudID, config.GetCloudID())
	require.Equal(t, apiKey, config.GetAPIKey())
	require.Equal(t, serviceToken, config.GetServiceToken())
	require.Equal(t, certificateFingerprint, config.GetCertificateFingerprint())
	require.Equal(t, indexName, config.GetIndex())
	require.Equal(t, indexType, config.GetType())
}
