package example

import (
	"errors"

	"github.com/jclem/crdt/rgass"
)

// Op is an operation sent to a site
type Op struct {
	Type       string
	Target     rgass.ID
	TargetList []rgass.ID
	Pos        int
	Len        int
	Str        string
	ID         rgass.ID
}

// Site is an individual editor of an RGASS.
type Site struct {
	session   int
	id        int
	vec       int
	rg        *rgass.RGASS
	OutStream chan Op
	open      bool
}

const bufferSize = 100

// NewSite creates a new Site.
func NewSite(session int, id int) Site {
	rg := rgass.NewRGASS()
	return Site{
		session:   session,
		id:        id,
		rg:        &rg,
		OutStream: make(chan Op, bufferSize),
		open:      true,
	}
}

// Close closes the stream's channel and stops it from accepting input
func (s *Site) Close() {
	s.open = false
	close(s.OutStream)
}

// Receive processes an incoming operation.
func (s *Site) Receive(op Op) error {
	if op.Type == "insert" {
		return s.rg.RemoteInsert(op.Target, op.Pos, op.Str, op.ID)
	}

	return s.rg.RemoteDelete(op.TargetList, op.Pos, op.Len)
}

// Insert inserts a string into the site at `pos` (a position in the visible text)
func (s *Site) Insert(pos int, str string) error {
	if !s.open {
		return errors.New("Site is not open")
	}

	count := 0

	var node *rgass.Node
	for node = range s.rg.Model.Iter() {
		if node.Hidden {
			continue
		}

		if count+len(node.Str) >= pos {
			pos -= count
			break
		}

		count += len(node.Str)
	}

	id := s.idFor(pos, len(str))
	if err := s.rg.LocalInsert(node.ID, pos, str, id); err != nil {
		return err
	}

	tarNode := node.GetAncestor()
	return s.broadcastInsert(tarNode.ID, node.AncestorOffset+pos, str, id)
}

// Delete deletes a string from the site at `pos` of length `len` (a position in the visible text)
func (s *Site) Delete(pos int, delLen int) error {
	if !s.open {
		return errors.New("Site is not open")
	}

	count := 0

	var node *rgass.Node
	for node = range s.rg.Model.Iter() {
		if node.Hidden {
			continue
		}

		if count+len(node.Str) >= pos {
			pos -= count
			break
		}

		count += len(node.Str)
	}

	nodeList, effectiveLen, err := s.rg.LocalDelete(node.ID, pos, delLen)
	if err != nil {
		return err
	}

	delIDList := make([]rgass.ID, len(nodeList))
	effectivePos := pos
	for i, node := range nodeList {
		if i == 0 {
			effectivePos += node.AncestorOffset
		}
		delIDList[i] = node.GetAncestor().ID
	}
	return s.broadcastDelete(delIDList, effectivePos, effectiveLen)
}

// Text returns the site's text
func (s *Site) Text() string {
	return s.rg.Text()
}

func (s *Site) broadcastDelete(targetList []rgass.ID, pos int, delLen int) error {
	return s.broadcast(Op{
		Type:       "delete",
		Pos:        pos,
		Len:        delLen,
		TargetList: targetList,
	})
}

func (s *Site) broadcastInsert(tarID rgass.ID, pos int, str string, id rgass.ID) error {
	return s.broadcast(Op{
		Type:   "insert",
		Target: tarID,
		Pos:    pos,
		Str:    str,
		ID:     id,
	})
}

func (s *Site) broadcast(op Op) error {
	select {
	case s.OutStream <- op:
		return nil
	default:
		return errors.New("Site buffer full")
	}
}

func (s *Site) idFor(pos int, len int) rgass.ID {
	id := rgass.ID{
		Session: s.session,
		Vector:  s.vec,
		Site:    s.id,
		Length:  len,
	}
	s.vec++
	return id
}
