package merkel

import (
	"crypto/sha256"
	"encoding/hex"
)

type Raw interface {
	Hash() []byte
	String() string
}

type Tree struct {
	Root *Node
	Raws []Raw
}

type Node struct {
	left, right, parent *Node

	Raw  Raw
	Hash []byte
}

func New(raws []Raw) *Tree {
	root := buildWithNodes(buildWithRaw(raws))
	return &Tree{
		Root: root,
		Raws: raws,
	}
}

func (t *Tree) String() string {
	s := ""
	for _, r := range t.Raws {
		s += r.String() + "\n"
	}

	return s
}

func (t *Tree) Rebuild(raws []Raw) {
	root := buildWithNodes(buildWithRaw(raws))
	t.Root = root
	t.Raws = raws
}

func (n *Node) HashString() string {
	return hex.EncodeToString(n.Hash)
}

func buildWithRaw(raws []Raw) []*Node {
	ns := []*Node{}

	for _, r := range raws {
		n := &Node{
			Hash: r.Hash(),
			Raw:  r,
		}

		ns = append(ns, n)
	}

	if len(ns)%2 == 1 {
		ns = append(ns, nil)
	}

	return ns
}

func buildWithNodes(ns []*Node) *Node {
	var parentNs []*Node

	for i := 0; i < len(ns); i += 2 {
		left := ns[i]
		right := ns[i+1]

		h := sha256.New()
		h.Write(left.Hash)
		if right != nil {
			h.Write(right.Hash)
		}

		parent := &Node{
			left: left,
			Hash: h.Sum(nil),
		}
		if right != nil {
			parent.right = right
		}

		left.parent = parent
		if right != nil {
			right.parent = parent
		}

		parentNs = append(parentNs, parent)
		if len(parentNs) == 2 {
			return parent
		}
	}

	return buildWithNodes(parentNs)
}
