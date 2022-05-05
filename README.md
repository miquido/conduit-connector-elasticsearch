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

| name                     | description                                                                                    | required | default |
|--------------------------|------------------------------------------------------------------------------------------------|----------|---------|
| `version`                | The version of the Elasticsearch service. One of: `7`, `8`.                                    | true     |         |
| `host`                   | Server host.                                                                                   | true     |         |
| `username`               | The username used to authenticate.                                                             | false    |         |
| `password`               | The password used to authenticate. Required when username was provided.                        | false    |         |
| `cloudId`                | Endpoint for the Elastic Service (https://elastic.co/cloud).                                   | false    |         |
| `apiKey`                 | Base64-encoded token for authorization; if set, overrides username/password and service token. | false    |         |
| `serviceToken`           | Service token for authorization; if set, overrides username/password.                          | false    |         |
| `certificateFingerprint` | SHA256 hex fingerprint given by Elasticsearch on first launch.                                 | false    |         |
| `index`                  | The name of the target index.                                                                  | true     |         |
| `bulkSize`               | The number of items stored in bulk in the index.                                               | true     | `1000`  |

# Testing

Run `make test` to run all the unit and integration tests, which require Docker to be installed and running. The command will handle starting and stopping docker containers for you.

# References

- https://github.com/elastic/go-elasticsearch
- https://www.elastic.co/guide/en/elasticsearch/reference/7.17/docs-bulk.html
