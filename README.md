# crdt

A [conflict-free replicated data type][crdt] (CRDT) is a data structure that can
be updated concurrently (and in different orders) at different replicas. As long
as all replicas receive all states or operations on the CRDT, the values of the
replicas converge.

This package implements several CRDTs in Go. The goal is not necessarily to
create production-ready CRDTs (for example, many of them do not support safe
atomic operations, nor are they thread-safe), but rather as an exercise for me
to get more comfortable with Go and to learn about CRDTs.

## Included CRDTs

- [gcounter](gcounter/) A grow-only counter
- [LWW Register](lwwregister/) A last-write wins register
- [pncounter](pncounter/) A counter which can increment or decrement
- [rgass](rgass/) A CRDT for efficient string-based collaborative editing

[crdt]: https://en.wikipedia.org/wiki/Conflict-free_replicated_data_type
