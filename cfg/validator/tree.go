package validator

import (
	"errors"

	"github.com/andrskom/gowfnet/cfg"
)

var (
	ErrNodeIsNotFound        = errors.New("node is not found by id")
	ErrUnexpectedColorOfNode = errors.New("unexpected color of node")
	ErrStackIsEmpty          = errors.New("stack is empty")
)

type Tree struct {
	startNodeID  string
	finishNodeID string
	registry     map[string]*TreeNode
}

func NewTree(startNodeID string, finishNodeID string) *Tree {
	return &Tree{
		startNodeID:  startNodeID,
		finishNodeID: finishNodeID,
		registry:     make(map[string]*TreeNode),
	}
}

func (t *Tree) AddNode(n *TreeNode) {
	t.registry[n.GetID()] = n
}

func (t *Tree) GetStartNode() (*TreeNode, error) {
	return t.GetNode(t.startNodeID)
}

func (t *Tree) GetFinishNode() (*TreeNode, error) {
	return t.GetNode(t.finishNodeID)
}

func (t *Tree) GetNode(nodeID string) (*TreeNode, error) {
	res, ok := t.registry[nodeID]
	if !ok {
		return nil, ErrNodeIsNotFound
	}

	return res, nil
}

func (t *Tree) GetNodeRegistry() map[string]*TreeNode {
	return t.registry
}

type Color int

const (
	ColorWhite Color = iota
	ColorGray
	ColorBlack
)

type TreeNode struct {
	id    string
	color Color
	to    map[string]*TreeNode
	from  map[string]*TreeNode
}

func NewTreeNode(id string) *TreeNode {
	return &TreeNode{
		id:    id,
		color: ColorWhite,
		to:    make(map[string]*TreeNode),
		from:  make(map[string]*TreeNode),
	}
}

func (n *TreeNode) SetColor(c Color) error {
	switch c {

	}

	if n.color != c-1 {
		return ErrUnexpectedColorOfNode
	}

	n.color = c

	return nil
}

func (n *TreeNode) GetColor() Color {
	return n.color
}

func (n *TreeNode) IsColor(c Color) bool {
	return n.color == c
}

func (n *TreeNode) GetID() string {
	return n.id
}

func (n *TreeNode) AddTo(node *TreeNode) {
	n.to[node.GetID()] = node
}

func (n *TreeNode) GetTo() map[string]*TreeNode {
	return n.to
}

func (n *TreeNode) AddFrom(node *TreeNode) {
	if _, ok := n.from[node.GetID()]; ok {
		return
	}

	n.from[node.GetID()] = node
}

func (n *TreeNode) GetFrom() map[string]*TreeNode {
	return n.from
}

type TreeBuilder interface {
	Build(c cfg.Interface) (*Tree, error)
}

// CfgTreeBuilder from config of net and return tree.
type CfgTreeBuilder struct {
}

func NewCfgTreeBuilder() *CfgTreeBuilder {
	return &CfgTreeBuilder{}
}

func (b *CfgTreeBuilder) Build(c cfg.Interface) (*Tree, error) {
	tree := NewTree(c.GetStart().GetID(), c.GetFinish().GetID())

	for _, place := range c.GetPlaces() {
		tree.AddNode(NewTreeNode(place.GetID()))
	}

	for _, tr := range c.GetTransitions().GetAsMap() {
		for _, f := range tr.GetFrom() {
			fromNode, err := tree.GetNode(f.GetID())
			if err != nil {
				return nil, err
			}

			for _, t := range tr.GetTo() {
				toNode, err := tree.GetNode(t.GetID())
				if err != nil {
					return nil, err
				}

				fromNode.AddTo(toNode)
				toNode.AddFrom(fromNode)
			}
		}
	}

	return tree, nil
}

type NodeStack struct {
	stack []*TreeNode
}

func NewNodeStack() *NodeStack {
	return &NodeStack{stack: make([]*TreeNode, 0)}
}

func (s *NodeStack) Len() int {
	return len(s.stack)
}

func (s *NodeStack) Push(node *TreeNode) {
	s.stack = append(s.stack, node)
}

func (s *NodeStack) Peek() (*TreeNode, error) {
	if s.Len() == 0 {
		return nil, ErrStackIsEmpty
	}

	return s.stack[s.Len()-1], nil
}

func (s *NodeStack) Pop() (*TreeNode, error) {
	if s.Len() == 0 {
		return nil, ErrStackIsEmpty
	}

	res := s.stack[s.Len()-1]

	s.stack = s.stack[:s.Len()-1]

	return res, nil
}
