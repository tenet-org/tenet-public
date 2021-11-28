package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// ModuleCdc defines a protobuf encoding codec used for serialization. The actual codec should be
// provided to x/ibc EVM and defined at the application level.
var ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

// RegisterInterfaces register the ibc transfer module interfaces to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// 	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgTransfer{})

	// 	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
