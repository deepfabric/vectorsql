package tree

import "fmt"

type OrderStatement interface {
	Statement
	orderStatement()
}

func (*Top) orderStatement()    {}
func (*Ftop) orderStatement()   {}
func (OrderBy) orderStatement() {}

type Top struct {
	N ExprStatement
}

type Ftop struct {
	N ExprStatement
}

func (t *Top) String() string {
	return fmt.Sprintf("TOP %s", t.N.String())
}

func (t *Ftop) String() string {
	return fmt.Sprintf("FTOP %s", t.N.String())
}

type OrderBy []*Order

type Order struct {
	Type Direction
	E    ExprStatement
}

// Direction for ordering results.
type Direction int8

// Direction values.
const (
	DefaultDirection Direction = iota
	Ascending
	Descending
)

var directionName = [...]string{
	DefaultDirection: "",
	Ascending:        "ASC",
	Descending:       "DESC",
}

func (i Direction) String() string {
	if i < 0 || i > Direction(len(directionName)-1) {
		return fmt.Sprintf("Direction(%d)", i)
	}
	return directionName[i]
}

func (n *Order) String() string {
	var s string

	s += n.E.String()
	if n.Type != DefaultDirection {
		s = " " + n.Type.String()
	}
	return s
}

func (n OrderBy) String() string {
	var s string

	s += "ORDER BY "
	for i := range n {
		if i > 0 {
			s += ", "
		}
		s += n[i].String()
	}
	return s
}
