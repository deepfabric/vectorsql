package tree

import (
	"fmt"
)

type JoinClause struct {
	Type        JoinType
	Cond        JoinCond
	Left, Right RelationStatement
}

type JoinType int

const (
	FullOp JoinType = iota
	LeftOp
	RightOp
	CrossOp
	InnerOp
	NaturalOp
)

var joinTypeName = [...]string{
	FullOp:    "FULL JOIN",
	LeftOp:    "LEFT jOIN",
	RightOp:   "RIGHT JOIN",
	CrossOp:   "CROSS JOIN",
	InnerOp:   "INNER JOIN",
	NaturalOp: "NATURAL JOIN",
}

func (i JoinType) String() string {
	if i < 0 || i > JoinType(len(joinTypeName)-1) {
		return fmt.Sprintf("UnionType(%d)", i)
	}
	return joinTypeName[i]
}

// JoinCond represents a join condition.
type JoinCond interface {
	Statement
	joinCond()
}

func (*OnJoinCond) joinCond()  {}
func (*NonJoinCond) joinCond() {}

// OnJoinCond represents an ON join condition.
type OnJoinCond struct {
	E ExprStatement
}

type NonJoinCond struct {
}

func (n *OnJoinCond) String() string {
	return " ON " + n.E.String()
}

func (n *NonJoinCond) String() string {
	return ""
}

func (n *JoinClause) String() string {
	var s string

	s += n.Left.String()
	s += " " + n.Type.String()
	s += " " + n.Right.String()
	if n.Cond != nil {
		s += n.Cond.String()
	}
	return s
}
