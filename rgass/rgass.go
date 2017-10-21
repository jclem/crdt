package rgass

import (
	"errors"
)

// RGASS (replicated growable array supporting string) is a CRDT for efficient string-based
// collaborative editing.
type RGASS struct {
	Model Model
}

// NewRGASS creates a new RGASS.
func NewRGASS() RGASS {
	rgass := RGASS{Model: NewModel()}
	return rgass
}

// MustGet returns the node with the given ID, or panics if the node is not found.
func (r RGASS) MustGet(id ID) *Node {
	node, ok := r.Model.Get(id)
	if !ok {
		panic(errors.New("Node not found"))
	}
	return node
}

// Head returns the head sentinel node of the RGASS's internal model. This is the node that should
// be the target of operations wishing to insert at the start of the RGASS.
func (r RGASS) Head() *Node {
	return r.Model.Head()
}

// LocalInsert incorporates a locally-generated insert operation. (Algorithm 3, pp4)
func (r *RGASS) LocalInsert(tarID ID, pos int, str string, id ID) error {
	tarNode, ok := r.Model.Get(tarID)

	if !ok {
		return errors.New("Node not found")
	}

	return r.doInsert(tarNode, pos, str, id)
}

// LocalDelete incorporates a locally-generated delete operation. (Algorithm 6, pp4)
func (r *RGASS) LocalDelete(tarID ID, pos int, delLen int) ([]*Node, int, error) {
	tarNode, ok := r.Model.Get(tarID)
	nodeList := []*Node{tarNode}
	effectiveLen := delLen

	if !ok {
		return nodeList, effectiveLen, errors.New("Node not found in model")
	}

	if pos == 0 && delLen == tarNode.Length() {
		tarNode.DeleteWhole()
	}

	if pos == 0 && delLen < tarNode.Length() {
		fNode, lNode := tarNode.DeletePrior(delLen)
		r.Model.Replace(tarNode, fNode, lNode)
	}

	if pos > 0 && pos+delLen == tarNode.Length() {
		fNode, lNode := tarNode.DeleteLast(delLen)
		r.Model.Replace(tarNode, fNode, lNode)
	}

	if pos > 0 && pos+delLen < tarNode.Length() {
		fNode, mNode, lNode := tarNode.DeleteMiddle(pos, delLen)
		r.Model.Replace(tarNode, fNode, mNode, lNode)
	}

	if pos > 0 && pos+delLen > tarNode.Length() {
		remainingLen := delLen - (tarNode.Length() - pos)
		fNode, lNode := tarNode.DeleteLast(pos)
		r.Model.Replace(tarNode, fNode, lNode)

		node := lNode.Next

		for remainingLen > 0 {
			if remainingLen > node.Length() {
				remainingLen -= node.Length()
				node.DeleteWhole()
			} else {
				fNode, lNode := node.DeletePrior(remainingLen)
				r.Model.Replace(node, fNode, lNode)
				remainingLen = 0
			}

			if remainingLen > 0 {
				node = node.Next
				for node.Hidden {
					effectiveLen += node.Length()
					node = node.Next
				}
				nodeList = append(nodeList, node)
			}
		}
	}

	return nodeList, effectiveLen, nil
}

// RemoteInsert incorporates an insert from a remote site (Algorithm 4, pp4)
func (r *RGASS) RemoteInsert(tarID ID, pos int, str string, id ID) error {
	tarNode, err := r.Model.FindNode(tarID, pos)
	pos -= tarNode.AncestorOffset

	if err != nil {
		return err
	}

	return r.doInsert(tarNode, pos, str, id)
}

func (r *RGASS) doInsert(tarNode *Node, pos int, str string, id ID) error {
	newNode := &Node{ID: id, Str: str}

	if tarNode.Sentinel { // If we are targeting the head of the model
		return r.Model.InsertAfter(tarNode, newNode)
	}

	if pos == tarNode.Length() {
		return r.Model.InsertAfter(tarNode, newNode)
	}

	fNode, lNode, err := tarNode.SplitTwo(pos)
	if err != nil {
		return err
	}

	err = r.Model.Replace(tarNode, fNode, lNode)
	if err != nil {
		return err
	}

	return r.Model.InsertAfter(fNode, newNode)
}

// RemoteDelete incorporates a delete from a remote site (Algorithm 7, pp5)
func (r *RGASS) RemoteDelete(tarIDList []ID, pos int, delLen int) error {
	count := len(tarIDList)
	if len(tarIDList) == 0 {
		return nil
	}
	tarNode, ok := r.Model.Get(tarIDList[0])
	if !ok {
		return errors.New("Node not found in model")
	}

	if count == 1 {
		return r.doDelete(tarNode, pos, delLen)
	}

	r.doDelete(tarNode, pos, tarNode.Length()-pos)
	sumLen := tarNode.Length() - pos

	for _, tmpNodeID := range tarIDList[1 : len(tarIDList)-1] {
		tmpNode, ok := r.Model.Get(tmpNodeID)
		if !ok {
			return errors.New("Node not found in model")
		}
		r.doDelete(tmpNode, 0, tmpNode.Length())
		sumLen += tmpNode.Length()
	}

	lastLen := delLen - sumLen
	lastNode, ok := r.Model.Get(tarIDList[len(tarIDList)-1])
	if !ok {
		return errors.New("Node not found in model")
	}
	r.doDelete(lastNode, 0, lastLen)

	return nil
}

// Algorithm 8 (pp5)
func (r *RGASS) doDelete(node *Node, pos int, delLen int) error {
	if !node.Split {
		nodeLen := node.Length()

		if pos == 0 && delLen == nodeLen {
			node.DeleteWhole()
			return nil
		} else if pos == 0 && delLen < nodeLen {
			fNode, lNode := node.DeletePrior(delLen)
			return r.Model.Replace(node, fNode, lNode)
		} else if pos > 0 && pos+delLen == nodeLen {
			fNode, lNode := node.DeleteLast(pos)
			return r.Model.Replace(node, fNode, lNode)
		} else if pos > 0 && pos+delLen < nodeLen {
			fNode, mNode, lNode := node.DeleteMiddle(pos, delLen)
			return r.Model.Replace(node, fNode, mNode, lNode)
		} else {
			return errors.New("Delete length longer than node")
		}
	}

	c0 := node.List[0]
	c1 := node.List[1]

	var c2 *Node
	if len(node.List) == 3 {
		c2 = node.List[2]
	}

	c0Len := c0.Length()
	c1Len := c1.Length()

	if pos <= c0Len && pos+delLen <= c0Len { // |../.,...,...
		return r.doDelete(c0, pos, delLen)
	} else if pos <= c0Len && delLen-(c0Len-pos) <= c1Len { // |...,../.,...
		if err := r.doDelete(c0, pos, c0Len-pos); err != nil {
			return err
		}
		return r.doDelete(c1, 0, delLen-(c0Len-pos))
	} else if pos <= c0Len && delLen-(c0Len-pos) >= c1Len { // |...,...,../.
		firstLen := c0Len - pos
		if err := r.doDelete(c0, pos, firstLen); err != nil {
			return err
		}
		if err := r.doDelete(c1, 0, c1.Length()); err != nil {
			return err
		}
		return r.doDelete(c2, 0, delLen-firstLen-c1Len)
	} else if pos > c0Len && pos-c0Len <= c1Len && pos-c0Len+delLen <= c1Len { // ...,|../.,...
		return r.doDelete(c1, pos-c0Len, delLen)
	} else if pos > c0Len && pos-c0Len <= c1Len && pos-c0Len+delLen >= c1Len { // ...,|...,../.
		firstPos := pos - c0Len
		if err := r.doDelete(c1, firstPos, c1Len-firstPos); err != nil {
			return err
		}
		return r.doDelete(c2, 0, delLen-(c1Len-firstPos))
	} else if c2 != nil && pos > c0Len+c1Len && pos-c0Len-c1Len+delLen <= c2.Length() { // ...,...,|../.
		return r.doDelete(c2, pos-c0Len-c1Len, delLen)
	}

	return errors.New("Delete length longer than node")
}

// Text returns the current visible text of the RGASS
func (r RGASS) Text() string {
	str := ""

	for node := range r.Model.Iter() {
		if node.Hidden {
			continue
		}
		str += node.Str
	}

	return str
}
