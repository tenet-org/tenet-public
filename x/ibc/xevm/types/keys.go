package types

const (
	// ModuleName defines the IBC EVM name
	ModuleName = "xevm"

	// Version defines the current version the IBC tranfer
	// module supports
	Version = "xevm-1"

	// PortID is the default port id that EVM module binds to
	PortID = "xevm"

	// StoreKey is the store key string for IBC EVM
	StoreKey = ModuleName

	// RouterKey is the message route for IBC EVM
	RouterKey = ModuleName

	// QuerierRoute is the querier route for IBC transfer
	QuerierRoute = ModuleName
)

// PortKey defines the key to store the port ID in store
var PortKey = []byte{0x01}
