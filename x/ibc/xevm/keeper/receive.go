package keeper

import (
	"encoding/json"
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"

	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	"github.com/tharsis/evmos/x/ibc/xevm/types"
)

// OnRecvPacket processes a cross chain Ethereum transaction.
func (k Keeper) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	tx *ethtypes.Transaction,
) ([]byte, error) {
	if !k.IsEVMChain() {
		return nil, types.ErrNonEVMChain
	}

	if !k.GetParams(ctx).ReceiveEnabled {
		return nil, types.ErrReceiveDisabled
	}

	eip155ChainID := k.evmKeeper.ChainID()

	// validate packet data (ethereum tx) upon receiving
	if err := types.ValidateEthTx(tx, eip155ChainID); err != nil {
		return nil, err
	}

	// call evm
	k.evmKeeper.WithContext(ctx)

	evmParams := k.evmKeeper.GetParams(ctx)
	cfg := evmParams.ChainConfig.EthereumConfig(eip155ChainID)
	signer := ethtypes.MakeSigner(cfg, big.NewInt(ctx.BlockHeight()))

	baseFee := k.feeMarketKeeper.GetBaseFee(ctx)

	msg, err := tx.AsMessage(signer, baseFee)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidEthereumTx, "failed to cast to ethereum core.Message: %s", err.Error())
	}

	nonce := k.evmKeeper.GetNonce(msg.From())

	res, err := k.evmKeeper.ApplyMessage(msg, evmtypes.NewNoOpTracer(), true)
	if err != nil {
		return nil, err
	}

	k.evmKeeper.SetNonce(msg.From(), nonce+1)

	if res.Failed() {
		return nil, sdkerrors.Wrap(evmtypes.ErrVMExecution, res.VmError)
	}

	// return the JSON bytes of the response and send them as the
	// acknowledgement result
	return json.Marshal(res)
}

// OnAcknowledgementPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain. If the acknowledgement
// was a success then nothing occurs. If the acknowledgement failed, then
// the state transition is reverted.
func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, tx *ethtypes.Transaction, ack channeltypes.Acknowledgement) error {
	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		k.Logger(ctx).Info(
			"EVM state transition failed",
			"txhash", tx.Hash().String(),
			"counterparty-chain-ID", tx.ChainId().String(),
			"error", resp.Error,
		)
		return nil
	case *channeltypes.Acknowledgement_Result:
		if len(resp.Result) == 0 {
			return sdkerrors.Wrap(channeltypes.ErrInvalidAcknowledgement, "result ack bytes cannnot be empty")
		}

		var res evmtypes.MsgEthereumTxResponse
		err := json.Unmarshal(resp.Result, &res)
		if err != nil {
			return sdkerrors.Wrap(channeltypes.ErrInvalidAcknowledgement, err.Error())
		}

		if res.Failed() {
			// // TODO: refund gas?
			return sdkerrors.Wrapf(channeltypes.ErrInvalidAcknowledgement, "state transition execution failed: %s", res.VmError)
		}

		// TODO: use receiver handler so that other module can use the
		// returned bytes and logs from the state transition
		return nil
	default:
		return channeltypes.ErrInvalidAcknowledgement
	}
}

// OnTimeoutPacket returns an error
func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, tx *ethtypes.Transaction) error {
	return sdkerrors.Wrapf(channeltypes.ErrPacketTimeout, "IBC ethereum tx timeout: %s", tx.Hash())
}
