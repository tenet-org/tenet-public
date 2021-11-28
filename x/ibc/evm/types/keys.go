package types

const (
	// ModuleName defines the IBC EVM name
	ModuleName = "ibc-evm"

	// Version defines the current version the IBC tranfer
	// module supports
	Version = "evm-1" // TODO: create I

	// PortID is the default port id that EVM module binds to
	PortID = "evm"

	// StoreKey is the store key string for IBC EVM
	StoreKey = ModuleName

	// RouterKey is the message route for IBC EVM
	RouterKey = ModuleName

	// QuerierRoute is the querier route for IBC transfer
	QuerierRoute = ModuleName
)

// PortKey defines the key to store the port ID in store
var PortKey = []byte{0x01}
