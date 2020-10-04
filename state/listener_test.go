package state

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStubListener_BeforeMove_AlwaycNoErr(t *testing.T) {
	ctrl := gomock.NewController(t)

	assert.NoError(t, NewStubListener().BeforeMove(context.Background(), NewMockOpInterface(ctrl), nil, nil))
}
