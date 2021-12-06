package types

import (
	"math/big"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestEthereumTxPacket(t *testing.T) {
	txData := &ethtypes.DynamicFeeTx{
		ChainID:   big.NewInt(9000),
		Nonce:     1,
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(1),
		Value:     big.NewInt(0),
	}
	tx := ethtypes.NewTx(txData)
	bz, err := tx.MarshalJSON()
	require.NoError(t, err)
	require.NotEmpty(t, string(bz))
}
