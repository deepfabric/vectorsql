package tree

type TableName struct {
	N ColunmNameList
}

func (n *TableName) String() string {
	return n.N.String()
}
