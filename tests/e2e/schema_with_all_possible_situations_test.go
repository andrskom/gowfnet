// nolint:gochecknoglobals,funlen
package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/andrskom/gowfnet"
	"github.com/andrskom/gowfnet/cfg"
	"github.com/andrskom/gowfnet/listener/channel"
	"github.com/andrskom/gowfnet/state"
)

// configWithAllPossibleSituation cover workflow patterns.
// http://mlwiki.org/index.php/Workflow_Patterns
var configWithAllPossibleSituation = cfg.Minimal{
	Start:  "start",
	Finish: "finish",
	Places: []cfg.StringID{
		"start",
		"finish",
		"sequence",
		"parallelStart",
		"parallelTread1.position1",
		"parallelTread1.position2",
		"parallelTread2.position1",
		"parallelSync",
		"choice1.position1",
		"choice1.position2",
		"choice2.position1",
		"merge",
	},
	Transitions: cfg.MinimalTransitionRegistry{
		"toSeq1": {
			From: []cfg.StringID{"start"},
			To:   []cfg.StringID{"sequence"},
		},
		"toParallelStart": {
			From: []cfg.StringID{"sequence"},
			To:   []cfg.StringID{"parallelStart"},
		},
		"toParallel": {
			From: []cfg.StringID{"parallelStart"},
			To:   []cfg.StringID{"parallelTread2.position1", "parallelTread1.position1"},
		},
		"toParallelTread1.position2": {
			From: []cfg.StringID{"parallelTread1.position1"},
			To:   []cfg.StringID{"parallelTread1.position2"},
		},
		"toParallelSync": {
			From: []cfg.StringID{"parallelTread1.position2", "parallelTread2.position1"},
			To:   []cfg.StringID{"parallelSync"},
		},
		"toChoice1": {
			From: []cfg.StringID{"parallelSync"},
			To:   []cfg.StringID{"choice1.position1"},
		},
		"toChoice1.Position2": {
			From: []cfg.StringID{"choice1.position1"},
			To:   []cfg.StringID{"choice1.position2"},
		},
		"choice1.Position2ToMerge": {
			From: []cfg.StringID{"choice1.position2"},
			To:   []cfg.StringID{"merge"},
		},
		"toChoice2": {
			From: []cfg.StringID{"parallelSync"},
			To:   []cfg.StringID{"choice2.position1"},
		},
		"choice2ToMerge": {
			From: []cfg.StringID{"choice2.position1"},
			To:   []cfg.StringID{"merge"},
		},
		"toFinish": {
			From: []cfg.StringID{"merge"},
			To:   []cfg.StringID{"finish"},
		},
	},
}

func TestRunByAllNetWithoutErr(t *testing.T) {
	r := require.New(t)
	listener := channel.New(100)
	net := gowfnet.NewNet(configWithAllPossibleSituation)
	st := state.NewState()

	net.WithListener(listener)

	r.NoError(net.Start(context.Background(), st))
	r.True(
		state.ErrorIs(
			state.ErrCodeStateHasNotTokenInPlace,
			net.Transit(context.Background(), st, "toParallelStart"),
		),
	)
	r.NoError(net.Transit(context.Background(), st, "toSeq1"))
	r.NoError(net.Transit(context.Background(), st, "toParallelStart"))
	r.NoError(net.Transit(context.Background(), st, "toParallel"))
	r.True(
		state.ErrorIs(
			state.ErrCodeStateHasNotTokenInPlace,
			net.Transit(context.Background(), st, "toParallelSync"),
		),
	)
	r.NoError(net.Transit(context.Background(), st, "toParallelTread1.position2"))
	r.NoError(net.Transit(context.Background(), st, "toParallelSync"))
	r.NoError(net.Transit(context.Background(), st, "toChoice1"))
	r.True(
		state.ErrorIs(
			state.ErrCodeStateHasNotTokenInPlace,
			net.Transit(context.Background(), st, "toChoice2"),
		),
	)
	r.NoError(net.Transit(context.Background(), st, "toChoice1.Position2"))
	r.NoError(net.Transit(context.Background(), st, "choice1.Position2ToMerge"))
	r.NoError(net.Transit(context.Background(), st, "toFinish"))
	r.True(st.IsFinished())

	r.Equal("start", listener.ReadEvt())
	r.Equal("move_FROM:[]_TO:[start]", listener.ReadEvt())
	r.Equal("moved_FROM:[]_TO:[start]", listener.ReadEvt())
	r.Equal("started", listener.ReadEvt())

	r.Equal("transit_toParallelStart", listener.ReadEvt())

	r.Equal("transit_toSeq1", listener.ReadEvt())
	r.Equal("move_FROM:[start]_TO:[sequence]", listener.ReadEvt())
	r.Equal("moved_FROM:[start]_TO:[sequence]", listener.ReadEvt())
	r.Equal("toSeq1_transited", listener.ReadEvt())

	r.Equal("transit_toParallelStart", listener.ReadEvt())
	r.Equal("move_FROM:[sequence]_TO:[parallelStart]", listener.ReadEvt())
	r.Equal("moved_FROM:[sequence]_TO:[parallelStart]", listener.ReadEvt())
	r.Equal("toParallelStart_transited", listener.ReadEvt())

	r.Equal("transit_toParallel", listener.ReadEvt())
	r.Equal("move_FROM:[parallelStart]_TO:[parallelTread2.position1 parallelTread1.position1]", listener.ReadEvt())
	r.Equal("moved_FROM:[parallelStart]_TO:[parallelTread2.position1 parallelTread1.position1]", listener.ReadEvt())
	r.Equal("toParallel_transited", listener.ReadEvt())

	r.Equal("transit_toParallelSync", listener.ReadEvt())

	r.Equal("transit_toParallelTread1.position2", listener.ReadEvt())
	r.Equal("move_FROM:[parallelTread1.position1]_TO:[parallelTread1.position2]", listener.ReadEvt())
	r.Equal("moved_FROM:[parallelTread1.position1]_TO:[parallelTread1.position2]", listener.ReadEvt())
	r.Equal("toParallelTread1.position2_transited", listener.ReadEvt())

	r.Equal("transit_toParallelSync", listener.ReadEvt())
	r.Equal("move_FROM:[parallelTread1.position2 parallelTread2.position1]_TO:[parallelSync]", listener.ReadEvt())
	r.Equal("moved_FROM:[parallelTread1.position2 parallelTread2.position1]_TO:[parallelSync]", listener.ReadEvt())
	r.Equal("toParallelSync_transited", listener.ReadEvt())

	r.Equal("transit_toChoice1", listener.ReadEvt())
	r.Equal("move_FROM:[parallelSync]_TO:[choice1.position1]", listener.ReadEvt())
	r.Equal("moved_FROM:[parallelSync]_TO:[choice1.position1]", listener.ReadEvt())
	r.Equal("toChoice1_transited", listener.ReadEvt())

	r.Equal("transit_toChoice2", listener.ReadEvt())

	r.Equal("transit_toChoice1.Position2", listener.ReadEvt())
	r.Equal("move_FROM:[choice1.position1]_TO:[choice1.position2]", listener.ReadEvt())
	r.Equal("moved_FROM:[choice1.position1]_TO:[choice1.position2]", listener.ReadEvt())
	r.Equal("toChoice1.Position2_transited", listener.ReadEvt())

	r.Equal("transit_choice1.Position2ToMerge", listener.ReadEvt())
	r.Equal("move_FROM:[choice1.position2]_TO:[merge]", listener.ReadEvt())
	r.Equal("moved_FROM:[choice1.position2]_TO:[merge]", listener.ReadEvt())
	r.Equal("choice1.Position2ToMerge_transited", listener.ReadEvt())

	r.Equal("transit_toFinish", listener.ReadEvt())
	r.Equal("move_FROM:[merge]_TO:[finish]", listener.ReadEvt())
	r.Equal("moved_FROM:[merge]_TO:[finish]", listener.ReadEvt())
	r.Equal("finished", listener.ReadEvt())
	r.Equal("toFinish_transited", listener.ReadEvt())
}
