package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v2/modules/core/24-host"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

// msg types
const (
	TypeIBCEthereumTx = "ibc-ethereum-tx"
)

// MsgIBCEthereumTx creates a new MsgIBCEthereumTx instance
func NewMsgIBCEthereumTx(
	sourcePort, sourceChannel string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64,
	tx *ethtypes.Transaction,
) *MsgIBCEthereumTx {
	bz, err := tx.MarshalBinary()
	if err != nil {
		return nil
	}

	return &MsgIBCEthereumTx{
		SourcePort:       sourcePort,
		SourceChannel:    sourceChannel,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
		EthereumTx:       bz,
	}
}

// Route implements sdk.Msg
func (MsgIBCEthereumTx) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgIBCEthereumTx) Type() string {
	return TypeIBCEthereumTx
}

// ValidateBasic performs a basic check of the MsgIBCEthereumTx fields.
// NOTE: timeout height or timestamp values can be 0 to disable the timeout.
// NOTE: The recipient addresses format is not validated as the format defined by
// the chain is not known to IBC.
func (msg MsgIBCEthereumTx) ValidateBasic() error {
	if err := host.PortIdentifierValidator(msg.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(msg.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}

	tx := new(ethtypes.Transaction)
	err := tx.UnmarshalBinary(msg.EthereumTx)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidEthereumTx, "failed to unmarshal binary: %s", err.Error())
	}

	ethtx := new(evmtypes.MsgEthereumTx)
	if err := ethtx.FromEthereumTx(tx); err != nil {
		return err
	}

	return ethtx.ValidateBasic()
}

// GetSignBytes implements sdk.Msg.
func (msg MsgIBCEthereumTx) GetSignBytes() []byte {
	panic("amino encoding not supported")
}

// GetSigners implements sdk.Msg
func (msg MsgIBCEthereumTx) GetSigners() []sdk.AccAddress {
	// FIXME:
	// signer, err := sdk.AccAddressFromBech32(msg.Sender)
	// if err != nil {
	// 	panic(err)
	// }
	return []sdk.AccAddress{}
}
