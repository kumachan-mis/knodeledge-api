package record

type GraphChildEntry struct {
	Name        string
	Relation    string
	Description string
	Children    []GraphChildEntry
}
