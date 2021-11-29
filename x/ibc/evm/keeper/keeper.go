package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	clientkeeper "github.com/cosmos/ibc-go/v2/modules/core/02-client/keeper"
	host "github.com/cosmos/ibc-go/v2/modules/core/24-host"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	evmkeeper "github.com/tharsis/ethermint/x/evm/keeper"
	feemarketkeeper "github.com/tharsis/ethermint/x/feemarket/keeper"
	"github.com/tharsis/evmos/x/ibc/evm/types"
)

// Keeper defines the IBC EVM keeper
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	clientkeeper    clientkeeper.Keeper
	channelKeeper   types.ChannelKeeper
	portKeeper      types.PortKeeper
	evmKeeper       evmkeeper.Keeper
	feemarketKeeper feemarketkeeper.Keeper
	scopedKeeper    capabilitykeeper.ScopedKeeper
}

// NewKeeper creates a new IBC EMV Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	// TODO: client keeper
	channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper,
	evmKeeper evmkeeper.Keeper, feemarketKeeper feemarketkeeper.Keeper,
	scopedKeeper capabilitykeeper.ScopedKeeper,
) Keeper {

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:             cdc,
		storeKey:        key,
		paramstore:      paramSpace,
		channelKeeper:   channelKeeper,
		portKeeper:      portKeeper,
		evmKeeper:       evmKeeper,
		feemarketKeeper: feemarketKeeper,
		scopedKeeper:    scopedKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// IsBound checks if the IBC EVM module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// GetPort returns the portID for the IBC EVM module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the IBC EVM module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey, []byte(portID))
}

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the IBC EVM module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}
