package types

import (
	"math/big"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

// ValidateEthTx converts an ethereum transaction to an Ethermint MsgEthereumTx
// and performs a stateless validation of the fields.
func ValidateEthTx(tx *ethtypes.Transaction, chainID *big.Int) error {
	if tx == nil {
		return sdkerrors.Wrap(ErrInvalidEthereumTx, "tx cannot be nil")
	}

	if chainID == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidChainID, "chain ID cannot be nil")
	}

	if tx.ChainId().Cmp(chainID) != 0 {
		return sdkerrors.Wrapf(ErrInvalidEthereumTx, "chain ID mismatch, expected %d, got %d", chainID, tx.ChainId())
	}

	ethtx := new(evmtypes.MsgEthereumTx)
	if err := ethtx.FromEthereumTx(tx); err != nil {
		return err
	}

	return ethtx.ValidateBasic()
}
