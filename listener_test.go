package gowfnet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAutomaticListenerMiddleware_AfterPlaced_FullIntegrationTest(t *testing.T) {
	r := require.New(t)
	cfg := Cfg{
		Start:  "start",
		Finish: "finish",
		Places: []string{
			"start",
			"firstBranch",
			"firstBranchDone",
			"secondBranch",
			"secondBranchDone",
			"finish",
		},
		Transitions: map[string]CfgTransition{
			"toBranches": {
				From:        []string{"start"},
				To:          []string{"firstBranch", "secondBranch"},
				IsAutomatic: true,
			},
			"firstBranchDone": {
				From:        []string{"firstBranch"},
				To:          []string{"firstBranchDone"},
				IsAutomatic: true,
			},
			"secondBranchDone": {
				From:        []string{"secondBranch"},
				To:          []string{"secondBranchDone"},
				IsAutomatic: true,
			},
			"union": {
				From:        []string{"firstBranchDone", "secondBranchDone"},
				To:          []string{"finish"},
				IsAutomatic: true,
			},
		},
	}

	net, err := BuildFromConfig(cfg)
	r.NoError(err)
	net.SetListener(NewAutomaticListenerMiddleware(nil))

	st := NewState()
	err = net.Start(st, nil)
	r.NoError(err)
	r.True(st.IsFinished())
	r.False(st.IsError())
}
