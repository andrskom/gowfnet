// nolint:dupl
package validator

import (
	"errors"

	"github.com/andrskom/gowfnet/cfg"
)

type NonFinishPlaces struct {
	treeBuilder TreeBuilder
}

func NewNonFinishPlaces(treeBuilder TreeBuilder) *NonFinishPlaces {
	return &NonFinishPlaces{treeBuilder: treeBuilder}
}

func (n *NonFinishPlaces) Validate(c cfg.Interface) error {
	tree, err := n.treeBuilder.Build(c)
	if err != nil {
		return err
	}

	node, err := tree.GetFinishNode()
	if err != nil {
		return err
	}

	if err := n.dfsNodeProcessor(node); err != nil {
		return err
	}

	vErr := NewError()

	for _, node := range tree.GetNodeRegistry() {
		switch node.GetColor() {
		case ColorBlack:
		case ColorGray:
			vErr.Addf("unexpected situation for non-finish places validator in place with id '%s'", node.GetID())
		case ColorWhite:
			vErr.Addf("place with id '%s' is non-finish place", node.GetID())
		default:
			return errors.New("unexpected color of node")
		}
	}

	return PrepareResultErr(vErr)
}

func (n *NonFinishPlaces) dfsNodeProcessor(node *TreeNode) error {
	if !node.IsColor(ColorWhite) {
		return nil
	}

	if err := node.SetColor(ColorGray); err != nil {
		return err
	}

	for _, toNode := range node.GetFrom() {
		if err := n.dfsNodeProcessor(toNode); err != nil {
			return err
		}
	}

	return node.SetColor(ColorBlack)
}
