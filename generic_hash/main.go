package main

import (
	"fmt"
	"hash/maphash"
)

type Version struct {
	Major, Minor, Patch int
}

type releaseLine struct{}

func (releaseLine) Hash(h *maphash.Hash, v Version) {
	fmt.Fprintf(h, "%d.%d", v.Major, v.Minor)
}

func (releaseLine) Equal(a, b Version) bool {
	return a.Major == b.Major && a.Minor == b.Minor
}

func main() {
	var h maphash.Hasher[Version] = releaseLine{}

	v1 := Version{Major: 1, Minor: 27, Patch: 0}
	v2 := Version{Major: 1, Minor: 27, Patch: 4}
	fmt.Println(h.Equal(v1, v2))

	seed := maphash.MakeSeed()

	var a, b maphash.Hash
	a.SetSeed(seed)
	b.SetSeed(seed)

	h.Hash(&a, v1)
	h.Hash(&b, v2)
	fmt.Println(a.Sum64() == b.Sum64())
}
