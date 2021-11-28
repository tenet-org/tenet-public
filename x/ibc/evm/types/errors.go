package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IBC EVM sentinel errors
var (
	ErrInvalidPacketTimeout = sdkerrors.Register(ModuleName, 2, "invalid packet timeout")
	ErrInvalidDenomForEVM   = sdkerrors.Register(ModuleName, 3, "invalid denomination for EVM state transition")
	ErrInvalidVersion       = sdkerrors.Register(ModuleName, 4, "invalid IBC EVM version")
	ErrSendDisabled         = sdkerrors.Register(ModuleName, 5, "IBC EVM packets sent from this chain are disabled")
	ErrReceiveDisabled      = sdkerrors.Register(ModuleName, 6, "IBC EVM packets received to this chain are disabled")
	ErrMaxEVMChannels       = sdkerrors.Register(ModuleName, 7, "max IBC EVM channels")
	ErrInvalidEthereumTx    = sdkerrors.Register(ModuleName, 8, "invalid ethereum tx")
)
