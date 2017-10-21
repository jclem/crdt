package lwwregister

// An ID identifies a site updating an LWW register
type ID int64

// A Timestamp is a totally-orderable timestamp for a register update
type Timestamp struct {
	ID  ID
	Vec int64
}

// An LWWRegister is a last-write wins register.
type LWWRegister struct {
	id  ID
	vec int64
	ts  Timestamp
	Val interface{}
}

// NewRegister creates a new LWW register.
func NewRegister(id ID) *LWWRegister {
	return &LWWRegister{id, 0, Timestamp{id, 0}, nil}
}

// Update updates the LWW register value.
func (r *LWWRegister) Update(val interface{}) {
	r.vec++
	r.Val = val
	r.ts = Timestamp{r.id, r.vec}
}

// Incorporate incorporates a remote LWW update.
func (r *LWWRegister) Incorporate(ts Timestamp, val interface{}) {
	if r.ts.Compare(ts) == -1 {
		r.vec = ts.Vec + 1
		r.ts = ts
		r.Val = val
		return
	}
}

// Compare compares two timestamps.
func (t Timestamp) Compare(o Timestamp) int {
	if t.Vec < o.Vec {
		return -1
	}

	if t.Vec > o.Vec {
		return 1
	}

	if t.ID < o.ID {
		return -1
	}

	if t.ID > o.ID {
		return 1
	}

	return 0
}
