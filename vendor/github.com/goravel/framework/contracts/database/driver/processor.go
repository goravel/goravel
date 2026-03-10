package driver

type Processor interface {
	ProcessColumns(dbColumns []DBColumn) []Column
	ProcessForeignKeys(dbIndexes []DBForeignKey) []ForeignKey
	ProcessIndexes(dbIndexes []DBIndex) []Index
	ProcessTypes(types []Type) []Type
}

type DBColumn struct {
	Collation     string
	Comment       string
	Default       string
	Extra         string
	Name          string
	Nullable      string
	Type          string
	TypeName      string
	Length        int
	Places        int
	Precision     int
	Autoincrement bool
	Primary       bool
}

type DBForeignKey struct {
	Name           string
	Columns        string
	ForeignSchema  string
	ForeignTable   string
	ForeignColumns string
	OnUpdate       string
	OnDelete       string
}

type DBIndex struct {
	Columns string
	Name    string
	Type    string
	Primary bool
	Unique  bool
}

type ForeignKey struct {
	Name           string
	ForeignSchema  string
	ForeignTable   string
	OnUpdate       string
	OnDelete       string
	Columns        []string
	ForeignColumns []string
}

type Index struct {
	Name    string
	Type    string
	Columns []string
	Primary bool
	Unique  bool
}
