package app

type Attribute struct {
	AttrName  string
	AttrValue any
}

type AttributeList []Attribute

type Connector interface {
	Name() string
	ReadDir(path string) ([]AttributeList, error)
	ReadFile(path string) ([]byte, error)
}
