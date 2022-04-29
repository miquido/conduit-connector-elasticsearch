package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/miquido/conduit-connector-elasticsearch/destination"
	"github.com/miquido/conduit-connector-elasticsearch/internal"
)

func main() {
	sdk.Serve(internal.Specification, nil, destination.NewDestination)
}
