package gcounter

type id string

// A GCounter is a grow-only counter
type GCounter struct {
	id   id
	vals map[id]int
}

// NewGCounter creates a new GCounter
func NewGCounter(gid id) *GCounter {
	return &GCounter{id: gid, vals: make(map[id]int)}
}

// Increment increments the value at this site for the GCounter
func (g *GCounter) Increment() {
	g.vals[g.id]++
}

// Incorporate incorporates a remote GCounter value
func (g *GCounter) Incorporate(id id, val int) {
	if existingVal := g.vals[id]; existingVal < val {
		g.vals[id] = val
	}
}

// Value gets the value of the GCounter
func (g *GCounter) Value() int {
	sum := 0
	for _, val := range g.vals {
		sum += val
	}
	return sum
}
