package keeper

import (
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"

	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	"github.com/tharsis/evmos/x/ibc/evm/types"
)

// OnRecvPacket processes a cross chain Ethereum transaction.
func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, tx *ethtypes.Transaction) error {
	eip155ChainID := k.evmKeeper.ChainID()

	// validate packet data upon receiving
	if err := types.ValidateEthTx(tx, eip155ChainID); err != nil {
		return err
	}

	if !k.GetParams(ctx).ReceiveEnabled {
		return types.ErrReceiveDisabled
	}

	// call evm
	k.evmKeeper.WithContext(ctx)

	evmParams := k.evmKeeper.GetParams(ctx)
	cfg := evmParams.ChainConfig.EthereumConfig(eip155ChainID)
	signer := ethtypes.MakeSigner(cfg, big.NewInt(ctx.BlockHeight()))

	msg, err := tx.AsMessage(signer, nil)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidEthereumTx, "failed to cast to ethereum core.Message: %s", err.Error())
	}

	nonce := k.evmKeeper.GetNonce(msg.From())

	res, err := k.evmKeeper.ApplyMessage(msg, evmtypes.NewNoOpTracer(), true)
	if err != nil {
		return err
	}

	k.evmKeeper.SetNonce(msg.From(), nonce+1)

	if res.Failed() {
		return sdkerrors.Wrap(evmtypes.ErrVMExecution, res.VmError)
	}

	return nil
}

// OnAcknowledgementPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain. If the acknowledgement
// was a success then nothing occurs. If the acknowledgement failed, then
// the state transition is reverted.
func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, tx *ethtypes.Transaction, ack channeltypes.Acknowledgement) error {
	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return nil // TODO: revert
	default:
		// the acknowledgement succeeded on the receiving chain so nothing
		// needs to be executed and no error needs to be returned
		return nil
	}
}

// OnTimeoutPacket
func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, tx *ethtypes.Transaction) error {
	// TODO: revert
	return nil
}
