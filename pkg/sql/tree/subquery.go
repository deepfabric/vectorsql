package tree

type Subquery struct {
	Exists bool
	Select *Select
}

func (n *Subquery) String() string {
	var s string

	if n.Exists {
		s += "EXISTS "
	}
	s += "(" + n.Select.String() + ")"
	return s
}
