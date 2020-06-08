package tree

type AliasClause struct {
	Alias Name
	Cols  NameList
}

func (n *AliasClause) String() string {
	var s string

	s += n.Alias.String()
	if len(n.Cols) > 0 {
		s += "(" + n.Cols.String() + ")"
	}
	return s
}
