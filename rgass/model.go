package rgass

import (
	"errors"
)

// Model represents all nodes in the RGASS
type Model struct {
	head  *Node        // A sentinel head node
	tail  *Node        // A sentinel tail node
	table map[ID]*Node // A map of node IDs to nodes
}

// NewModel creates a new Model
func NewModel() Model {
	m := Model{table: make(map[ID]*Node)}
	head := &Node{Sentinel: true}
	tail := &Node{Sentinel: true}
	m.table[head.ID] = head
	head.Next = tail
	tail.Prev = head
	m.head = head
	m.tail = tail
	return m
}

// Get returns the node associated with the given ID in the model
func (m *Model) Get(id ID) (*Node, bool) {
	node, ok := m.table[id]
	return node, ok
}

// Head returns the sentinel head node in the model
func (m *Model) Head() *Node {
	return m.head
}

// FindNode finds a node given a root node ID and an offset (Algorithm 5, pp4)
func (m *Model) FindNode(tarID ID, pos int) (tarNode *Node, err error) {
	tarNode, ok := m.Get(tarID)
	if !ok {
		return tarNode, errors.New("Target node not in model")
	}

	if pos > tarNode.Length() {
		return tarNode, errors.New("Position outside of target node")
	}

	for tarNode.Split {
		if pos <= tarNode.List[0].Length() {
			tarNode = tarNode.List[0]
		} else if pos <= tarNode.List[1].ID.Offset+tarNode.List[1].Length() {
			pos -= tarNode.List[0].Length()
			tarNode = tarNode.List[1]
		} else if tarNode.List[2] != nil {
			pos -= (tarNode.List[0].Length() + tarNode.List[1].Length())
			tarNode = tarNode.List[2]
		}
	}

	return tarNode, err
}

// InsertAfter inserts the given node (totally ordered by ID)
func (m *Model) InsertAfter(tarNode *Node, newNodes ...*Node) error {
	if _, ok := m.Get(tarNode.ID); !ok {
		return errors.New("Target node not in model")
	}

	for _, newNode := range newNodes {
		if _, ok := m.Get(newNode.ID); ok {
			return errors.New("Node already in model")
		}

		m.table[newNode.ID] = newNode

		for nextNode := tarNode.Next; nextNode != m.tail; nextNode = nextNode.Next {
			if newNode.ID.Compare(nextNode.ID) == -1 {
				tarNode = nextNode
			} else {
				break
			}
		}
		linkAfter(tarNode, newNode)
		tarNode = newNode
	}
	return nil
}

// Replace replaces a node with new nodes
func (m *Model) Replace(tarNode *Node, newNodes ...*Node) error {
	if _, ok := m.Get(tarNode.ID); !ok {
		return errors.New("Target node not in model")
	}

	firstNewNode := newNodes[0]

	if firstNewNode == nil {
		return nil
	}

	if _, ok := m.Get(firstNewNode.ID); ok {
		return errors.New("Node already in model")
	}

	m.table[firstNewNode.ID] = firstNewNode
	linkAfter(tarNode, firstNewNode)

	tarNode = firstNewNode

	for _, newNode := range newNodes[1:] {
		m.table[newNode.ID] = newNode
		linkAfter(tarNode, newNode)
		tarNode = newNode
	}

	return nil
}

// Iter iterates over all nodes in the model
func (m *Model) Iter() <-chan *Node {
	ch := make(chan *Node)

	go func() {
		for node := m.head; node != m.tail; node = node.Next {
			ch <- node
		}
		close(ch)
	}()

	return ch
}

func linkAfter(tarNode *Node, newNode *Node) {
	newNode.Next = tarNode.Next
	newNode.Prev = tarNode
	newNode.Next.Prev = newNode
	tarNode.Next = newNode
}
