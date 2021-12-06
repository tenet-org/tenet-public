package keeper_test

import "github.com/tharsis/evmos/x/ibc/xevm/types"

func (suite *KeeperTestSuite) TestParams() {
	params := suite.app.XEvmKeeper.GetParams(suite.ctx)
	suite.Require().Equal(types.DefaultParams(), params)
	params.SendEnabled = false
	suite.app.XEvmKeeper.SetParams(suite.ctx, params)
	newParams := suite.app.XEvmKeeper.GetParams(suite.ctx)
	suite.Require().Equal(newParams, params)
}
