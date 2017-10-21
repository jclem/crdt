package rgass

// ID identifies a text node in an RGASS site.
type ID struct {
	Session int // The session identifier of the node's inserting site
	Vector  int // The vector clock value of the node's inserting site at the time it was inserted
	Site    int // The site identifier of the node's inserting site
	Offset  int // The offset of the node's content within its parent
	Length  int // The length of the node's content
}

// Compare compares an ID to another ID. (Definition 4, pp11)
func (i *ID) Compare(other ID) int {
	if i.Session < other.Session {
		return -1
	}

	if i.Session > other.Session {
		return 1
	}

	if i.Vector < other.Vector {
		return -1
	}

	if i.Vector > other.Vector {
		return 1
	}

	if i.Site < other.Site {
		return -1
	}

	if i.Site > other.Site {
		return 1
	}

	if i.Offset > other.Offset {
		return -1
	}

	if i.Offset < other.Offset {
		return 1
	}

	return 0
}
