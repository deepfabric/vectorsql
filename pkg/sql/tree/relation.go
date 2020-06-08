package tree

type RelationStatement interface {
	Statement
	relationStatement()
}

func (*TableName) relationStatement()     {}
func (*JoinClause) relationStatement()    {}
func (*UnionClause) relationStatement()   {}
func (*SelectClause) relationStatement()  {}
func (*AliasedSelect) relationStatement() {}

func (*AliasedTable) relationStatement() {}

type AliasedSelect struct {
	Sel *Select
	As  *AliasClause
}

func (n *AliasedSelect) String() string {
	var s string

	s += n.Sel.String()
	if n.As != nil {
		s += " AS " + n.As.String()
	}
	return s
}
