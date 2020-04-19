package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andrskom/gowfnet/cfg"
)

func TestNewTree(t *testing.T) {
	assert.Equal(
		t,
		&Tree{startNodeID: "a", finishNodeID: "z", registry: make(map[string]*TreeNode)},
		NewTree("a", "z"),
	)
}

func TestTree_AddNode_Repeatable(t *testing.T) {
	tree := NewTree("a", "z")

	treeNode := NewTreeNode("b")
	tree.AddNode(treeNode)
	node, err := tree.GetNode("b")
	assert.NoError(t, err)
	assert.Same(t, treeNode, node)

	treeNode = NewTreeNode("b")
	tree.AddNode(treeNode)
	node, err = tree.GetNode("b")
	assert.NoError(t, err)
	assert.Same(t, treeNode, node)
}

func TestTree_GetStartNode_Configured_ExpectedNode(t *testing.T) {
	tree := NewTree("a", "z")

	treeNode := NewTreeNode("a")
	tree.AddNode(treeNode)

	node, err := tree.GetStartNode()
	assert.NoError(t, err)
	assert.Same(t, treeNode, node)
}

func TestTree_GetStartNode_ConfiguredWithoutStart_ExpectedErr(t *testing.T) {
	tree := NewTree("a", "z")

	node, err := tree.GetStartNode()
	assert.Nil(t, node)
	assert.Equal(t, ErrNodeIsNotFound, err)
}

func TestTree_GetFinishNode_Configured_ExpectedNode(t *testing.T) {
	tree := NewTree("a", "z")

	treeNode := NewTreeNode("z")
	tree.AddNode(treeNode)

	node, err := tree.GetFinishNode()
	assert.NoError(t, err)
	assert.Same(t, treeNode, node)
}

func TestTree_GetFinishNode_ConfiguredWithoutFinish_ExpectedErr(t *testing.T) {
	tree := NewTree("a", "z")

	node, err := tree.GetFinishNode()
	assert.Nil(t, node)
	assert.Equal(t, ErrNodeIsNotFound, err)
}

func TestTree_GetNode_Configured_ExpectedNode(t *testing.T) {
	tree := NewTree("a", "z")

	treeNode := NewTreeNode("b")
	tree.AddNode(treeNode)

	node, err := tree.GetNode("b")
	assert.NoError(t, err)
	assert.Same(t, treeNode, node)
}

func TestTree_GetNode_ConfiguredWithoutStart_ExpectedErr(t *testing.T) {
	tree := NewTree("a", "z")

	node, err := tree.GetNode("b")
	assert.Nil(t, node)
	assert.Equal(t, ErrNodeIsNotFound, err)
}

func TestNewTreeNode(t *testing.T) {
	a := assert.New(t)

	node := NewTreeNode("a")

	a.Equal("a", node.GetID())
	a.True(node.IsColor(ColorWhite), "unexpected color of node")
	a.NotNil(node.GetFrom())
	a.Len(node.GetFrom(), 0)
	a.NotNil(node.GetTo())
	a.Len(node.GetTo(), 0)
}

func TestTreeNode_SetColor_CorrectColor_NoErr(t *testing.T) {
	node := NewTreeNode("a")

	require.NoError(t, node.SetColor(ColorGray))
	assert.True(t, node.IsColor(ColorGray), "unexpected color of node")
}

func TestTreeNode_SetColor_IncorrectColor_ExpectedErr(t *testing.T) {
	node := NewTreeNode("a")

	require.Equal(t, ErrUnexpectedColorOfNode, node.SetColor(ColorBlack))
	assert.True(t, node.IsColor(ColorWhite), "unexpected color of node")
}

func TestTreeNode_GetColor(t *testing.T) {
	node := NewTreeNode("a")

	assert.Equal(t, ColorWhite, node.GetColor())
}

func TestTreeNode_IsColor(t *testing.T) {
	node := NewTreeNode("a")

	t.Run("expected", func(t *testing.T) {
		assert.True(t, node.IsColor(ColorWhite))
	})

	t.Run("another one", func(t *testing.T) {
		assert.False(t, node.IsColor(ColorGray))
	})
}

func TestTreeNode_GetID(t *testing.T) {
	node := NewTreeNode("a")

	assert.Equal(t, "a", node.GetID())
}

func TestTreeNode_AddTo_Repeatable(t *testing.T) {
	node := NewTreeNode("a")

	to1 := NewTreeNode("b")
	node.AddTo(to1)

	assert.Equal(
		t,
		map[string]*TreeNode{
			"b": to1,
		},
		node.GetTo(),
	)

	to2 := NewTreeNode("b")
	node.AddTo(to2)

	assert.Equal(
		t,
		map[string]*TreeNode{
			"b": to2,
		},
		node.GetTo(),
	)
}

func TestTreeNode_AddFrom_Repeatable(t *testing.T) {
	node := NewTreeNode("a")

	from1 := NewTreeNode("b")
	node.AddFrom(from1)

	assert.Equal(
		t,
		map[string]*TreeNode{
			"b": from1,
		},
		node.GetFrom(),
	)

	from2 := NewTreeNode("b")
	node.AddFrom(from2)

	assert.Equal(
		t,
		map[string]*TreeNode{
			"b": from2,
		},
		node.GetFrom(),
	)
}

func TestBuildTree_CorrectCfg_ExpectedTree(t *testing.T) {
	minCfg := &cfg.Minimal{
		Start:  "a",
		Finish: "b",
		Places: []cfg.StringID{"a", "b"},
		Transitions: map[string]cfg.MinimalTransition{
			"c": {
				From: []cfg.StringID{"a"},
				To:   []cfg.StringID{"b"},
			},
		},
	}

	tree, err := BuildTree(minCfg)
	require.NoError(t, err)
	require.NotNil(t, tree)

	startNode, err := tree.GetStartNode()
	require.NoError(t, err)

	finishNode, err := tree.GetFinishNode()
	require.NoError(t, err)

	assert.Equal(t, "a", startNode.GetID())
	assert.Equal(t, map[string]*TreeNode{"b": finishNode}, startNode.GetTo())
	assert.Equal(t, map[string]*TreeNode{}, startNode.GetFrom())

	assert.Equal(t, "b", finishNode.GetID())
	assert.Equal(t, map[string]*TreeNode{}, finishNode.GetTo())
	assert.Equal(t, map[string]*TreeNode{"a": startNode}, finishNode.GetFrom())
}

func TestBuildTree_IncorrectCfg_ExpectedErr(t *testing.T) {
	dp := map[string]cfg.Minimal{
		"unexpected in from": {
			Start:  "a",
			Finish: "b",
			Places: []cfg.StringID{"a", "b"},
			Transitions: map[string]cfg.MinimalTransition{
				"c": {
					From: []cfg.StringID{"c"},
					To:   []cfg.StringID{"b"},
				},
			},
		},
		"unexpected in to": {
			Start:  "a",
			Finish: "b",
			Places: []cfg.StringID{"a", "b"},
			Transitions: map[string]cfg.MinimalTransition{
				"c": {
					From: []cfg.StringID{"a"},
					To:   []cfg.StringID{"c"},
				},
			},
		},
	}

	for descr, data := range dp {
		t.Run(descr, func(t *testing.T) {
			tree, err := BuildTree(&data)
			require.Nil(t, tree)
			require.Equal(t, ErrNodeIsNotFound, err)
		})
	}
}

func TestNewNodeStack(t *testing.T) {
	assert.Equal(t, &NodeStack{stack: make([]*TreeNode, 0)}, NewNodeStack())
}

func TestNodeStack_Push(t *testing.T) {
	stack := NewNodeStack()

	{
		node := NewTreeNode("1")
		stack.Push(node)
		assert.Equal(t, 1, stack.Len())
		topNode, err := stack.Peek()
		require.NoError(t, err)
		assert.Same(t, node, topNode)
	}

	{
		node := NewTreeNode("2")
		stack.Push(node)
		assert.Equal(t, 2, stack.Len())
		topNode, err := stack.Peek()
		require.NoError(t, err)
		assert.Same(t, node, topNode)
	}
}

func TestNodeStack_Pop_NotEmpty_ExpectedNode(t *testing.T) {
	stack := NewNodeStack()

	node := NewTreeNode("1")
	stack.Push(node)

	aNode, err := stack.Pop()
	require.NoError(t, err)
	assert.Equal(t, 0, stack.Len())
	assert.Same(t, node, aNode)
}

func TestNodeStack_Pop_Empty_ExpectedError(t *testing.T) {
	stack := NewNodeStack()

	aNode, err := stack.Pop()
	assert.Nil(t, aNode)
	assert.Equal(t, ErrStackIsEmpty, err)
	assert.Equal(t, 0, stack.Len())
}

func TestNodeStack_Peek_NotEmpty_ExpectedNode(t *testing.T) {
	stack := NewNodeStack()

	node := NewTreeNode("1")
	stack.Push(node)

	aNode, err := stack.Peek()
	require.NoError(t, err)
	assert.Equal(t, 1, stack.Len())
	assert.Same(t, node, aNode)
}

func TestNodeStack_Peek_Empty_ExpectedError(t *testing.T) {
	stack := NewNodeStack()

	aNode, err := stack.Peek()
	assert.Nil(t, aNode)
	assert.Equal(t, ErrStackIsEmpty, err)
	assert.Equal(t, 0, stack.Len())
}
