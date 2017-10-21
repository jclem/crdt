package rgass

import (
	"errors"
)

// Node represents a node of text inside an RGASS
type Node struct {
	ID             ID      // The identifier of this node
	List           []*Node // A list of nodes this node has been split into
	Str            string  // The string contents of this node
	Split          bool    // Whether the node has been split
	Sentinel       bool    // Whether the node is a sentinel (just a marker for the head)
	Hidden         bool    // Whether the node is hidden
	Next           *Node   // A pointer to the next node in the model
	Prev           *Node   // A pointer to the previous node in the model
	Ancestor       *Node   // A pointer to a child node's most distant ancestor
	AncestorOffset int     // The offset of this node from its most distant ancestor
}

// GetAncestor gets the node ancestor, or the node itself if it is an ancestor
func (n *Node) GetAncestor() *Node {
	if n.Ancestor != nil {
		return n.Ancestor
	}
	return n
}

// DeleteLast deletes the last part of the node.
func (n *Node) DeleteLast(pos int) (*Node, *Node) {
	fNode, lNode, _ := n.SplitTwo(pos)
	lNode.Hidden = true
	return fNode, lNode
}

// DeleteMiddle deletes the middle part of the node.
func (n *Node) DeleteMiddle(pos int, len int) (*Node, *Node, *Node) {
	fNode, mNode, lNode, _ := n.SplitThree(pos, len)
	mNode.Hidden = true
	return fNode, mNode, lNode
}

// DeletePrior deletes the prior part of the node.
func (n *Node) DeletePrior(pos int) (*Node, *Node) {
	fNode, lNode, _ := n.SplitTwo(pos)
	fNode.Hidden = true
	return fNode, lNode
}

// DeleteWhole deletes an entire node.
func (n *Node) DeleteWhole() *Node {
	n.Hidden = true
	return n
}

// Length gets the length of the node.
func (n *Node) Length() int {
	return n.ID.Length
}

// SplitThree splits a node in three with a given position and length (Algorithm 2, pp3)
func (n *Node) SplitThree(pos int, delLen int) (*Node, *Node, *Node, error) {
	var fNode Node
	var mNode Node
	var lNode Node

	if err := n.checkPos(pos); err != nil {
		return &fNode, &mNode, &lNode, err
	}

	fNode = *n
	fNode.ID.Length = pos
	fNode.Str = n.Str[0:pos]
	fNode.Ancestor = n
	fNode.AncestorOffset = n.AncestorOffset

	mNode = *n
	mNode.ID.Length = delLen
	mNode.ID.Offset = fNode.ID.Offset + pos
	mNode.Str = n.Str[pos : pos+delLen]
	mNode.Ancestor = n
	mNode.AncestorOffset = n.AncestorOffset + mNode.ID.Offset

	lNode = *n
	lNode.ID.Offset = mNode.ID.Offset + delLen
	lNode.ID.Length = n.Length() - fNode.Length() - mNode.Length()
	lNode.Str = n.Str[pos+delLen:]
	lNode.Ancestor = n
	lNode.AncestorOffset = n.AncestorOffset + lNode.ID.Offset

	n.Hidden = true
	n.Split = true
	n.List = []*Node{&fNode, &mNode, &lNode}

	return &fNode, &mNode, &lNode, nil
}

// SplitTwo splits a node in two at a given position (Algorithm 1, pp3)
func (n *Node) SplitTwo(pos int) (*Node, *Node, error) {
	var fNode Node
	var lNode Node

	if err := n.checkPos(pos); err != nil {
		return &fNode, &lNode, err
	}

	fNode = *n
	fNode.ID.Length = pos
	fNode.Str = n.Str[0:pos]
	fNode.Ancestor = n
	fNode.AncestorOffset = n.AncestorOffset

	lNode = *n
	lNode.ID.Offset = n.ID.Offset + pos
	lNode.ID.Length = n.ID.Length - pos
	lNode.Str = n.Str[pos:]
	lNode.Ancestor = n
	lNode.AncestorOffset = n.AncestorOffset + lNode.ID.Offset

	n.Hidden = true
	n.Split = true
	n.List = []*Node{&fNode, &lNode}

	return &fNode, &lNode, nil
}

func (n Node) checkPos(pos int) error {
	var err error

	if pos < 0 {
		err = errors.New("Position in node can not be less than 0")
	}

	if pos > n.ID.Length {
		err = errors.New("Position in node can not be greater than node length")
	}

	return err
}
