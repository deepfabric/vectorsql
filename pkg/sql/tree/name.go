package tree

type Name string

func (n Name) String() string {
	return string(n)
}

type NameList []Name

func (n NameList) String() string {
	var s string

	for i := range n {
		if i > 0 {
			s += ", "
		}
		s += n[i].String()
	}
	return s
}

type ColunmName struct {
	Path  Name
	Index ExprStatement
}

func (n ColunmName) String() string {
	var s string

	s += n.Path.String()
	if n.Index != nil {
		s += "[" + n.Index.String() + "]"
	}
	return s
}

type ColunmNameList []ColunmName

func (n ColunmNameList) String() string {
	var s string

	for i := range n {
		if i > 0 {
			s += "."
		}
		s += n[i].String()
	}
	return s
}
