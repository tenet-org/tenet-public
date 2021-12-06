package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	connectiontypes "github.com/cosmos/ibc-go/v2/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v2/modules/core/exported"

	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, ibcexported.ClientState, error)
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channeltypes.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	SendPacket(ctx sdk.Context, channelCap *capabilitytypes.Capability, packet ibcexported.PacketI) error
}

// ClientKeeper defines the expected IBC client keeper
type ClientKeeper interface {
	GetClientConsensusState(ctx sdk.Context, clientID string) (connection ibcexported.ConsensusState, found bool)
}

// ConnectionKeeper defines the expected IBC connection keeper
type ConnectionKeeper interface {
	GetConnection(ctx sdk.Context, connectionID string) (connection connectiontypes.ConnectionEnd, found bool)
}

// PortKeeper defines the expected IBC port keeper
type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability
}

// EVMKeeper defines the expected EVM keeper
type EVMKeeper interface {
	ChainID() *big.Int
	WithContext(ctx sdk.Context)
	GetParams(ctx sdk.Context) evmtypes.Params
	GetNonce(addr common.Address) uint64
	ApplyMessage(msg core.Message, tracer vm.Tracer, commit bool) (*evmtypes.MsgEthereumTxResponse, error)
	SetNonce(addr common.Address, nonce uint64)
}

// v defines the expected Fee Market keeper
type FeeMarketKeeper interface {
	GetBaseFee(sdk.Context) *big.Int
}
