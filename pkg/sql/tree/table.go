package tree

type TableStatement interface {
	Statement
	tableStatement()
}

func (*Subquery) tableStatement()     {}
func (*TableName) tableStatement()    {}
func (*AliasedTable) tableStatement() {}

type AliasedTable struct {
	As  *AliasClause
	Tbl TableStatement
}

func (n *AliasedTable) String() string {
	var s string

	s += n.Tbl.String()
	if n.As != nil {
		s += " AS " + n.As.String()
	}
	return s
}
