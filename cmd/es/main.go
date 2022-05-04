package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
	es "github.com/miquido/conduit-connector-elasticsearch"
	esDestination "github.com/miquido/conduit-connector-elasticsearch/destination"
)

func main() {
	sdk.Serve(es.Specification, nil, esDestination.NewDestination)
}
