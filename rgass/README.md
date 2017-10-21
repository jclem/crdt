# RGASS

A replicated growable array (RGA) is a linked list of elements ordered
consistentently with causality. If two concurrent inserts target the same
element, they are ordered consistently according to their timestamps.

RGAs do not work well for strings, where great overhead is necessary in
representing each individual character in the string as an element in a linked
list. [RGA supporting string (RGASS)][rgass] is an improvement on the RGA for
efficient collaborative editing of strings.

This package is an implementation of RGASS in Go.

[rgass]: http://www.sciencedirect.com/science/article/pii/S1474034616301811
