package internal

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)

func Specification() sdk.Specification {
	return sdk.Specification{
		Name:    "elasticsearch",
		Summary: "An Elasticsearch destination plugin for Conduit.",
		Version: "v0.1.0",
		Author:  "Miquido",
		DestinationParams: map[string]sdk.Parameter{
			"host": {
				Default:     "",
				Required:    true,
				Description: "Server host.",
			},
			"username": {
				Default:     "",
				Required:    false,
				Description: "The username used to authenticate.",
			},
			"password": {
				Default:     "",
				Required:    false,
				Description: "The password used to authenticate. Required when username was provided.",
			},
			"index": {
				Default:     "",
				Required:    true,
				Description: "The name of the target index.",
			},
			"bulkSize": {
				Default:     "1000",
				Required:    true,
				Description: "The number of items stored in bulk in the index.",
			},
		},
		SourceParams: map[string]sdk.Parameter{
			//
		},
	}
}
