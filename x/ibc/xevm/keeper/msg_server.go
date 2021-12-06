package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/tharsis/evmos/x/ibc/xevm/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) IBCEthereumTx(goCtx context.Context, msg *types.MsgXEVM) (*types.MsgXEVMResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tx := new(ethtypes.Transaction)
	if err := tx.UnmarshalBinary(msg.EthereumTx); err != nil {
		return nil, err
	}

	if err := k.SendIBCEthereumTx(
		ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, tx,
	); err != nil {
		return nil, err
	}

	k.Logger(ctx).Info("IBC cross EVM transaction", "txhash", tx.Hash().String())

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

	return &types.MsgXEVMResponse{}, nil
}
