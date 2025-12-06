package app

type (
	AttrName  string
	AttrValue any
)

type AttributeList map[AttrName]AttrValue

type Connector interface {
	Name() string
	ReadDir(path string) ([]AttributeList, error)
	ReadFile(path string) ([]byte, error)
}

func (a AttrName) String() string {
	return string(a)
}
