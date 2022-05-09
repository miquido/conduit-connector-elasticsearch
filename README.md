# Conduit Connector Elasticsearch

## General
The Elasticsearch plugin is one of [Conduit](https://github.com/ConduitIO/conduit) builtin plugins.
It currently provides only destination Elasticsearch connector, allowing for using it as a destination in a Conduit pipeline.

## How to build it
Run `make`.

# Source

Not supported.

# Destination

The Destination connector stores data in given index.

## Configuration Options

| name                     | description                                                                                                                                                              | required                                | default |
|--------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------|---------|
| `version`                | The version of the Elasticsearch service. One of: `6`, `7`, `8`.                                                                                                         | `true`                                  |         |
| `host`                   | Server host.                                                                                                                                                             | `true`                                  |         |
| `username`               | [v: 6, 7, 8] The username used to authenticate.                                                                                                                          | `false`                                 |         |
| `password`               | [v: 6, 7, 8] The password used to authenticate. Required when username was provided.                                                                                     | `false`                                 |         |
| `cloudId`                | [v: 6, 7, 8] Endpoint for the Elastic Service (https://elastic.co/cloud).                                                                                                | `false`                                 |         |
| `apiKey`                 | [v: 6, 7, 8] Base64-encoded token for authorization; if set, overrides username/password and service token.                                                              | `false`                                 |         |
| `serviceToken`           | [v: 7, 8] Service token for authorization; if set, overrides username/password.                                                                                          | `false`                                 |         |
| `certificateFingerprint` | [v: 7, 8] SHA256 hex fingerprint given by Elasticsearch on first launch.                                                                                                 | `false`                                 |         |
| `index`                  | The name of the target index.                                                                                                                                            | `true`                                  |         |
| `type`                   | [v: 6] The name of the index's type to write the data to.                                                                                                                | `true` for version 6, `false` otherwise |         |
| `bulkSize`               | The number of items stored in bulk in the index. Minimum is `1`, maximum is `10000`. Beware that greater sizes than `1000` may require additional service configuration. | `true`                                  | `1000`  |

# Testing

Run `make test` to run all the unit and integration tests, which require Docker to be installed and running. The command will handle starting and stopping docker containers for you.

# References

- https://github.com/elastic/go-elasticsearch
- https://www.elastic.co/guide/en/elasticsearch/reference/7.17/docs-bulk.html
