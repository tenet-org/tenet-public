package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/tharsis/evmos/x/ibc/evm/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) IBCEthereumTx(goCtx context.Context, msg *types.MsgIBCEthereumTx) (*types.MsgIBCEthereumTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tx := new(ethtypes.Transaction)
	if err := tx.UnmarshalJSON(msg.EthereumTx); err != nil {
		return nil, err
	}

	if err := k.SendIBCEthereumTx(
		ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, tx,
	); err != nil {
		return nil, err
	}

	// k.Logger(ctx).Info("IBC fungible token transfer", "token", msg.Token.Denom, "amount", msg.Token.Amount.String(), "sender", msg.Sender, "receiver", msg.Receiver)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEthereumTx,
			// sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			// sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	return &types.MsgIBCEthereumTxResponse{}, nil
}
