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

	t.Run("fails when Host is empty", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion: elasticsearch.Version6,
			"nonExistentKey": "value",
		})

		require.EqualError(t, err, fmt.Sprintf("%q config value must be set", ConfigKeyHost))
	})

	t.Run("fails when Username is provided but Password is empty", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:  elasticsearch.Version6,
			ConfigKeyHost:     fakerInstance.Internet().URL(),
			ConfigKeyUsername: fakerInstance.Internet().Email(),
			"nonExistentKey":  "value",
		})

		require.EqualError(t, err, fmt.Sprintf("%q config value must be set when %q is provided", ConfigKeyPassword, ConfigKeyUsername))
	})

	t.Run("fails when Password is provided but Username is empty", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:  elasticsearch.Version6,
			ConfigKeyHost:     fakerInstance.Internet().URL(),
			ConfigKeyPassword: fakerInstance.Internet().Password(),
			"nonExistentKey":  "value",
		})

		require.EqualError(t, err, fmt.Sprintf("%q config value must be set when %q is provided", ConfigKeyUsername, ConfigKeyPassword))
	})

	t.Run("fails when Index is empty", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion: elasticsearch.Version6,
			ConfigKeyHost:    fakerInstance.Internet().URL(),
			"nonExistentKey": "value",
		})

		require.EqualError(t, err, fmt.Sprintf("%q config value must be set", ConfigKeyIndex))
	})

	for _, version := range []elasticsearch.Version{
		elasticsearch.Version6,
	} {
		t.Run(fmt.Sprintf("fails when Type is empty for Version=%s", version), func(t *testing.T) {
			_, err := ParseConfig(map[string]string{
				ConfigKeyVersion: version,
				ConfigKeyHost:    fakerInstance.Internet().URL(),
				ConfigKeyIndex:   fakerInstance.Lorem().Word(),
				"nonExistentKey": "value",
			})

			require.EqualError(t, err, fmt.Sprintf("%q config value must be set", ConfigKeyType))
		})
	}

	t.Run("fails when Bulk Size is empty", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion: elasticsearch.Version8,
			ConfigKeyHost:    fakerInstance.Internet().URL(),
			ConfigKeyIndex:   fakerInstance.Lorem().Word(),
			"nonExistentKey": "value",
		})

		require.EqualError(t, err, fmt.Sprintf("%q config value must be set", ConfigKeyBulkSize))
	})

	t.Run("fails when Bulk Size is an invalid positive integer", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:  elasticsearch.Version8,
			ConfigKeyHost:     fakerInstance.Internet().URL(),
			ConfigKeyIndex:    fakerInstance.Lorem().Word(),
			ConfigKeyBulkSize: "-1",
			"nonExistentKey":  "value",
		})

		require.EqualError(t, err, fmt.Sprintf(`failed to parse %q config value: strconv.ParseUint: parsing "-1": invalid syntax`, ConfigKeyBulkSize))
	})

	t.Run("fails when Bulk Size is less than 1", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:  elasticsearch.Version8,
			ConfigKeyHost:     fakerInstance.Internet().URL(),
			ConfigKeyIndex:    fakerInstance.Lorem().Word(),
			ConfigKeyBulkSize: "0",
			"nonExistentKey":  "value",
		})

		require.EqualError(t, err, fmt.Sprintf("failed to parse %q config value: value must be greater than 0", ConfigKeyBulkSize))
	})

	t.Run("fails when Bulk Size is greater than 10 000", func(t *testing.T) {
		_, err := ParseConfig(map[string]string{
			ConfigKeyVersion:  elasticsearch.Version8,
			ConfigKeyHost:     fakerInstance.Internet().URL(),
			ConfigKeyIndex:    fakerInstance.Lorem().Word(),
			ConfigKeyBulkSize: "10001",
			"nonExistentKey":  "value",
		})

		require.EqualError(t, err, fmt.Sprintf("failed to parse %q config value: value must be less than 10 000", ConfigKeyBulkSize))
	})

	t.Run("returns config when all required config values were provided", func(t *testing.T) {
		var cfgRaw = map[string]string{
			ConfigKeyVersion:  elasticsearch.Version8,
			ConfigKeyHost:     fakerInstance.Internet().URL(),
			ConfigKeyIndex:    fakerInstance.Lorem().Word(),
			ConfigKeyBulkSize: "1",
			"nonExistentKey":  "value",
		}

		config, err := ParseConfig(cfgRaw)

		require.NoError(t, err)
		require.Equal(t, cfgRaw[ConfigKeyVersion], config.Version)
		require.Equal(t, cfgRaw[ConfigKeyHost], config.Host)
		require.Equal(t, cfgRaw[ConfigKeyIndex], config.Index)
		require.Equal(t, cfgRaw[ConfigKeyBulkSize], fmt.Sprintf("%d", config.BulkSize))
		require.Empty(t, "", config.Username)
		require.Empty(t, "", config.Password)
		require.Empty(t, "", config.Type)
		require.Empty(t, "", config.CloudID)
		require.Empty(t, "", config.APIKey)
		require.Empty(t, "", config.ServiceToken)
		require.Empty(t, "", config.CertificateFingerprint)
	})

	t.Run("returns config when all config values were provided", func(t *testing.T) {
		var cfgRaw = map[string]string{
			ConfigKeyVersion:                elasticsearch.Version8,
			ConfigKeyHost:                   fakerInstance.Internet().URL(),
			ConfigKeyIndex:                  fakerInstance.Lorem().Word(),
			ConfigKeyType:                   fakerInstance.Lorem().Word(),
			ConfigKeyBulkSize:               fmt.Sprintf("%d", fakerInstance.Int32Between(1, 10_000)),
			ConfigKeyUsername:               fakerInstance.Internet().Email(),
			ConfigKeyPassword:               fakerInstance.Internet().Password(),
			ConfigKeyCloudID:                fakerInstance.RandomStringWithLength(32),
			ConfigKeyAPIKey:                 fakerInstance.RandomStringWithLength(32),
			ConfigKeyServiceToken:           fakerInstance.RandomStringWithLength(32),
			ConfigKeyCertificateFingerprint: fakerInstance.Hash().SHA256(),
			"nonExistentKey":                "value",
		}

		config, err := ParseConfig(cfgRaw)

		require.NoError(t, err)
		require.Equal(t, cfgRaw[ConfigKeyVersion], config.Version)
		require.Equal(t, cfgRaw[ConfigKeyHost], config.Host)
		require.Equal(t, cfgRaw[ConfigKeyIndex], config.Index)
		require.Equal(t, cfgRaw[ConfigKeyType], config.Type)
		require.Equal(t, cfgRaw[ConfigKeyBulkSize], fmt.Sprintf("%d", config.BulkSize))
		require.Equal(t, cfgRaw[ConfigKeyUsername], config.Username)
		require.Equal(t, cfgRaw[ConfigKeyPassword], config.Password)
		require.Equal(t, cfgRaw[ConfigKeyCloudID], config.CloudID)
		require.Equal(t, cfgRaw[ConfigKeyAPIKey], config.APIKey)
		require.Equal(t, cfgRaw[ConfigKeyServiceToken], config.ServiceToken)
		require.Equal(t, cfgRaw[ConfigKeyCertificateFingerprint], config.CertificateFingerprint)
	})
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
	)

	config := Config{
		Host:                   host,
		Username:               username,
		Password:               password,
		CloudID:                cloudID,
		APIKey:                 apiKey,
		ServiceToken:           serviceToken,
		CertificateFingerprint: certificateFingerprint,
	}

	require.Equal(t, host, config.GetHost())
	require.Equal(t, username, config.GetUsername())
	require.Equal(t, password, config.GetPassword())
	require.Equal(t, cloudID, config.GetCloudID())
	require.Equal(t, apiKey, config.GetAPIKey())
	require.Equal(t, serviceToken, config.GetServiceToken())
	require.Equal(t, certificateFingerprint, config.GetCertificateFingerprint())
}
