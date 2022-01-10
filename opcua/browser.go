package opcua

type BrowsedNode struct {
	NodeID      string
	DataType    string
	Description string
	Unit        string
	Scale       string
	BrowseName  string
}

type Browser interface {
	Browse(string, string) ([]BrowsedNode, error)
}
