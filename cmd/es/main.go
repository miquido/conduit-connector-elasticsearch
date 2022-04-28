package main

import (
	"github.com/conduitio/conduit-connector-elasticsearch/destination"
	"github.com/conduitio/conduit-connector-elasticsearch/internal"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

func main() {
	sdk.Serve(internal.Specification, nil, destination.NewDestination)
}
