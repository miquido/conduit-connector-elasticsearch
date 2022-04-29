package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/miquido/conduit-connector-elasticsearch"
	"github.com/miquido/conduit-connector-elasticsearch/destination"
)

func main() {
	sdk.Serve(elasticsearch.Specification, nil, destination.NewDestination)
}
