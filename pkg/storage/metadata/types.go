package metadata

type Attribute struct {
	Index bool
	Type  uint32 // type of attribute
	Name  string // name of attribute
}

type Metadata struct {
	IsE   bool
	Attrs []Attribute
}
