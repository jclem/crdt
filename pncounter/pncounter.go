package pncounter

type id string

// A PNCounter is a counter that can both grow and shrink
type PNCounter struct {
	id   id
	vals map[id]*[2]int
}

// NewPNCounter creates a new PNCounter
func NewPNCounter(gid id) *PNCounter {
	return &PNCounter{id: gid, vals: make(map[id]*[2]int)}
}

// Increment increments the value at this site for the PNCounter
func (p *PNCounter) Increment() {
	p.vals[p.id][0] = 1
}

// Decrement decrements the value at this site for the PNCounter
func (p *PNCounter) Decrement() {
	p.vals[p.id][1]++
}

// Incorporate incorporates a remote GCounter value
func (p *PNCounter) Incorporate(id id, vals [2]int) {
	if vals[0] > p.vals[id][0] {
		p.vals[id][0] = vals[0]
	}

	if vals[1] > p.vals[id][1] {
		p.vals[id][1] = vals[1]
	}
}

// Value gets the value of the PNCounter
func (p *PNCounter) Value() int {
	sum := 0
	for _, val := range p.vals {
		sum += val[0]
		sum -= val[1]
	}
	return sum
}
