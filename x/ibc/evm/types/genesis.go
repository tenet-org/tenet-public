package types

import (
	host "github.com/cosmos/ibc-go/v2/modules/core/24-host"
)

// NewGenesisState creates a new ibc-evm GenesisState instance.
func NewGenesisState(portID string, params Params) *GenesisState {
	return &GenesisState{
		PortId: portID,
		Params: params,
	}
}

// DefaultGenesisState returns a GenesisState with "evm" as the default PortID.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		PortId: PortID,
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}
	return gs.Params.Validate()
}
