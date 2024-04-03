// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/elastic/elastic-agent/pkg/control/v2/client"
	atesting "github.com/elastic/elastic-agent/pkg/testing"
	"github.com/elastic/elastic-agent/pkg/testing/define"
	"github.com/elastic/elastic-agent/pkg/testing/tools/testcontext"
)

var simpleConfig1 = `
outputs:
  default:
    type: fake-action-output
    shipper.enabled: true
inputs:
  - id: fake
    type: fake
    state: 1
    message: Configuring
`

var simpleConfig2 = `
outputs:
  default:
    type: fake-action-output
    shipper.enabled: true
inputs:
  - id: fake
    type: fake
    state: 2
    message: Healthy
`

var simpleNonGroupedConfig = `
outputs:
  default:
    type: fake-action-output
inputs:
  - id: fake-non-grouped
    type: fake-non-grouped
    state: 1
    message: Configuring
`

var complexNonGroupedConfig = `
outputs:
  default:
    type: fake-action-output
inputs:
  - id: fake-non-grouped-0
    type: fake-non-grouped
    state: 2
    message: Healthy
  - id: fake-non-grouped-1
    type: fake-non-grouped
    state: 2
    message: Healthy
`

func TestFakeComponent(t *testing.T) {
	define.Require(t, define.Requirements{
		Group: Default,
		Local: true,
	})

	f, err := define.NewFixture(t, define.Version())
	require.NoError(t, err)

	ctx, cancel := testcontext.WithDeadline(t, context.Background(), time.Now().Add(10*time.Minute))
	defer cancel()
	err = f.Prepare(ctx, fakeComponent, fakeShipper)
	require.NoError(t, err)

	err = f.Run(ctx, atesting.State{
		Configure:  simpleConfig1,
		AgentState: atesting.NewClientState(client.Healthy),
		Components: map[string]atesting.ComponentState{
			"fake-default": {
				State: atesting.NewClientState(client.Healthy),
				Units: map[atesting.ComponentUnitKey]atesting.ComponentUnitState{
					atesting.ComponentUnitKey{UnitType: client.UnitTypeOutput, UnitID: "fake-default"}: {
						State: atesting.NewClientState(client.Healthy),
					},
					atesting.ComponentUnitKey{UnitType: client.UnitTypeInput, UnitID: "fake-default-fake"}: {
						State: atesting.NewClientState(client.Configuring),
					},
				},
			},
		},
	}, atesting.State{
		Configure:  simpleConfig2,
		AgentState: atesting.NewClientState(client.Healthy),
		StrictComponents: map[string]atesting.ComponentState{
			"fake-default": {
				State: atesting.NewClientState(client.Healthy),
				Units: map[atesting.ComponentUnitKey]atesting.ComponentUnitState{
					atesting.ComponentUnitKey{UnitType: client.UnitTypeOutput, UnitID: "fake-default"}: {
						State: atesting.NewClientState(client.Healthy),
					},
					atesting.ComponentUnitKey{UnitType: client.UnitTypeInput, UnitID: "fake-default-fake"}: {
						State: atesting.NewClientState(client.Healthy),
					},
				},
			},
			"fake-shipper-default": {
				State: atesting.NewClientState(client.Healthy),
				Units: map[atesting.ComponentUnitKey]atesting.ComponentUnitState{
					atesting.ComponentUnitKey{UnitType: client.UnitTypeOutput, UnitID: "fake-shipper-default"}: {
						State: atesting.NewClientState(client.Healthy),
					},
					atesting.ComponentUnitKey{UnitType: client.UnitTypeInput, UnitID: "fake-default"}: {
						State: atesting.NewClientState(client.Healthy),
					},
				},
			},
		},
	})
	require.NoError(t, err)
}

func TestFakeNonGroupedComponent(t *testing.T) {
	define.Require(t, define.Requirements{
		Group: Default,
		Local: true,
	})

	f, err := define.NewFixture(t, define.Version())
	require.NoError(t, err)

	ctx, cancel := testcontext.WithDeadline(t, context.Background(), time.Now().Add(10*time.Minute))
	defer cancel()
	err = f.Prepare(ctx, fakeComponent, fakeShipper)
	require.NoError(t, err)

	err = f.Run(ctx, atesting.State{
		Configure:  simpleNonGroupedConfig,
		AgentState: atesting.NewClientState(client.Healthy),
		Components: map[string]atesting.ComponentState{
			"fake-non-grouped-default-fake-non-grouped": {
				State: atesting.NewClientState(client.Healthy),
				Units: map[atesting.ComponentUnitKey]atesting.ComponentUnitState{
					atesting.ComponentUnitKey{UnitType: client.UnitTypeOutput, UnitID: "fake-action-output"}: {
						State: atesting.NewClientState(client.Healthy),
					},
					atesting.ComponentUnitKey{UnitType: client.UnitTypeInput, UnitID: "fake-non-grouped-default-fake-non-grouped-unit"}: {
						State: atesting.NewClientState(client.Configuring),
					},
				},
			},
		},
	}, atesting.State{
		Configure:  complexNonGroupedConfig,
		AgentState: atesting.NewClientState(client.Healthy),
		Components: map[string]atesting.ComponentState{
			"fake-non-grouped-default-fake-non-grouped-0": {
				State: atesting.NewClientState(client.Healthy),
				Units: map[atesting.ComponentUnitKey]atesting.ComponentUnitState{
					atesting.ComponentUnitKey{UnitType: client.UnitTypeOutput, UnitID: "fake-action-output"}: {
						State: atesting.NewClientState(client.Healthy),
					},
					atesting.ComponentUnitKey{UnitType: client.UnitTypeInput, UnitID: "fake-non-grouped-default-fake-non-grouped-0-unit"}: {
						State: atesting.NewClientState(client.Healthy),
					},
				},
			},
			"fake-non-grouped-default-fake-non-grouped-1": {
				State: atesting.NewClientState(client.Healthy),
				Units: map[atesting.ComponentUnitKey]atesting.ComponentUnitState{
					atesting.ComponentUnitKey{UnitType: client.UnitTypeOutput, UnitID: "fake-action-output"}: {
						State: atesting.NewClientState(client.Healthy),
					},
					atesting.ComponentUnitKey{UnitType: client.UnitTypeInput, UnitID: "fake-non-grouped-default-fake-non-grouped-1-unit"}: {
						State: atesting.NewClientState(client.Healthy),
					},
				},
			},
			// "fake-shipper-default": {
			// 	State: atesting.NewClientState(client.Healthy),
			// 	Units: map[atesting.ComponentUnitKey]atesting.ComponentUnitState{
			// 		atesting.ComponentUnitKey{UnitType: client.UnitTypeOutput, UnitID: "fake-shipper-default"}: {
			// 			State: atesting.NewClientState(client.Healthy),
			// 		},
			// 		atesting.ComponentUnitKey{UnitType: client.UnitTypeInput, UnitID: "fake-non-grouped-default-fake-non-grouped-0"}: {
			// 			State: atesting.NewClientState(client.Healthy),
			// 		},
			// 		atesting.ComponentUnitKey{UnitType: client.UnitTypeInput, UnitID: "fake-non-grouped-default-fake-non-grouped-1"}: {
			// 			State: atesting.NewClientState(client.Healthy),
			// 		},
			// 	},
			// },
		},
	})
	require.NoError(t, err)
}
