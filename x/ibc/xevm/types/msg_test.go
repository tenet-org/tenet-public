package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"github.com/tharsis/ethermint/tests"

	"github.com/cosmos/ibc-go/v2/modules/core/02-client/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type MsgsTestSuite struct {
	suite.Suite
}

func TestMsgsTestSuite(t *testing.T) {
	suite.Run(t, new(MsgsTestSuite))
}

func (suite *MsgsTestSuite) TestMsgXEVMGetters() {
	msgInvalid := MsgXEVM{}
	msg := NewMsgXEVM(
		PortID, "", types.NewHeight(1, 10), 0, ethtypes.NewTx(&ethtypes.DynamicFeeTx{}),
	)
	suite.Require().Equal(RouterKey, msg.Route())
	suite.Require().Equal(TypeIBCEthereumTx, msg.Type())
	suite.Require().Panics(func() { msgInvalid.GetSignBytes() })
	// suite.Require().Nil(msgInvalid.GetSigners())
	// suite.Require().NotNil(msg.GetSigners())
}

func (suite *MsgsTestSuite) TestMsgIBCEthereum_GetSigners() {
}

func (suite *MsgsTestSuite) TestMsgXEVM() {
	invalidBz, err := ethtypes.NewTx(&ethtypes.DynamicFeeTx{}).MarshalBinary()
	suite.Require().NoError(err)

	txData := &ethtypes.DynamicFeeTx{
		ChainID:    big.NewInt(9000),
		Nonce:      1,
		GasTipCap:  big.NewInt(1),
		GasFeeCap:  big.NewInt(10),
		Gas:        21_000,
		To:         nil,
		Value:      big.NewInt(1),
		AccessList: ethtypes.AccessList{},
	}

	_, pk := tests.NewAddrKey()
	privkey, ok := pk.(*ethsecp256k1.PrivKey)
	suite.Require().True(ok)

	key, err := privkey.ToECDSA()
	suite.Require().NoError(err)

	signer := ethtypes.LatestSignerForChainID(txData.ChainID)

	tx, err := ethtypes.SignNewTx(key, signer, txData)
	suite.Require().NoError(err)
	bz, err := tx.MarshalBinary()
	suite.Require().NoError(err)

	testCases := []struct {
		name       string
		msg        MsgXEVM
		expectPass bool
	}{
		{
			"invalid source port",
			MsgXEVM{SourcePort: ""},
			false,
		},
		{
			"invalid source channel",
			MsgXEVM{SourcePort: PortID, SourceChannel: ""},
			false,
		},
		{
			"empty ethereum tx",
			MsgXEVM{
				SourcePort:    PortID,
				SourceChannel: "channel-10",
			},
			false,
		},
		{
			"empty ethereum tx",
			MsgXEVM{
				SourcePort:    PortID,
				SourceChannel: "channel-10",
				EthereumTx:    invalidBz,
			},
			false,
		},
		{
			"valid tx",
			MsgXEVM{
				SourcePort:    PortID,
				SourceChannel: "channel-10",
				EthereumTx:    bz,
			},
			true,
		},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s", i, tc.name)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s", i, tc.name)
		}
	}
}
