package tree

import "fmt"

type UnionClause struct {
	All         bool
	Type        UnionType
	Left, Right RelationStatement
}

type UnionType int

const (
	UnionOp UnionType = iota
	IntersectOp
	ExceptOp
)

var unionTypeName = [...]string{
	UnionOp:     "UNION",
	IntersectOp: "INTERSECT",
	ExceptOp:    "EXCEPT",
}

func (i UnionType) String() string {
	if i < 0 || i > UnionType(len(unionTypeName)-1) {
		return fmt.Sprintf("UnionType(%d)", i)
	}
	return unionTypeName[i]
}

func (n *UnionClause) String() string {
	var s string

	s += n.Left.String()
	s += " " + n.Type.String()
	if n.All {
		s += " ALL"
	}
	s += " " + n.Right.String()
	return s
}
