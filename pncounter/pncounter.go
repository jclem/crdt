package pncounter

type id string

// A PNCounter is a counter that can both grow and shrink.
type PNCounter struct {
	id   id
	vals map[id]*[2]int
}

// NewPNCounter creates a new PNCounter.
func NewPNCounter(gid id) *PNCounter {
	return &PNCounter{
		id:   gid,
		vals: map[id]*[2]int{gid: &[2]int{}},
	}
}

// Increment increments the value at this site for the PNCounter.
func (p *PNCounter) Increment() {
	p.vals[p.id][0]++
}

// Decrement decrements the value at this site for the PNCounter.
func (p *PNCounter) Decrement() {
	p.vals[p.id][1]++
}

// Incorporate incorporates a remote GCounter value.
func (p *PNCounter) Incorporate(id id, siteVal [2]int) {
	if _, ok := p.vals[id]; !ok {
		p.vals[id] = &siteVal
		return
	}

	if siteVal[0] > p.vals[id][0] {
		p.vals[id][0] = siteVal[0]
	}

	if siteVal[1] > p.vals[id][1] {
		p.vals[id][1] = siteVal[1]
	}
}

// Value gets the value of the PNCounter.
func (p *PNCounter) Value() int {
	sum := 0
	for _, val := range p.vals {
		sum += val[0]
		sum -= val[1]
	}
	return sum
}
