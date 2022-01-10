package gopcua

import "errors"

// variable config

// const protocol = "opcua"
// const token = ""

var (
	// errNotFoundServerURI = errors.New("route map not found for Server URI")
	// errNotFoundNodeID    = errors.New("route map not found for Node ID")
	// errNotFoundConn      = errors.New("connection not found")

	errFailedConn = errors.New("failed to connect")
	// errFailedRead          = errors.New("failed to read")
	// errFailedParseInterval = errors.New("failed to parse subscription interval")
	// errFailedSub           = errors.New("failed to subscribe")
	// errFailedFindEndpoint  = errors.New("failed to find suitable endpoint")
	// errFailedFetchEndpoint = errors.New("failed to fetch OPC-UA server endpoints")
	// errFailedParseNodeID   = errors.New("failed to parse NodeID")
	// errFailedCreateReq     = errors.New("failed to create request")
	// errResponseStatus      = errors.New("response status not OK")
)
