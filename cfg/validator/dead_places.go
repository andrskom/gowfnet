// nolint:dupl
package validator

import (
	"errors"

	"github.com/andrskom/gowfnet/cfg"
)

type DeadPlaces struct {
	treeBuilder TreeBuilder
}

func NewDeadPlaces(treeBuilder TreeBuilder) *DeadPlaces {
	return &DeadPlaces{treeBuilder: treeBuilder}
}

func (d *DeadPlaces) Validate(c cfg.Interface) error {
	tree, err := d.treeBuilder.Build(c)
	if err != nil {
		return err
	}

	node, err := tree.GetStartNode()
	if err != nil {
		return err
	}

	if err := d.dfsNodeProcessor(node); err != nil {
		return err
	}

	vErr := NewError()

	for _, node := range tree.GetNodeRegistry() {
		switch node.GetColor() {
		case ColorBlack:
		case ColorGray:
			vErr.Addf("unexpected situation for dead places validator in place with id '%s'", node.GetID())
		case ColorWhite:
			vErr.Addf("place with id '%s' is dead place", node.GetID())
		default:
			return errors.New("unexpected color of node")
		}
	}

	return PrepareResultErr(vErr)
}

func (d *DeadPlaces) dfsNodeProcessor(node *TreeNode) error {
	if !node.IsColor(ColorWhite) {
		return nil
	}

	if err := node.SetColor(ColorGray); err != nil {
		return err
	}

	for _, toNode := range node.GetTo() {
		if err := d.dfsNodeProcessor(toNode); err != nil {
			return err
		}
	}

	return node.SetColor(ColorBlack)
}
