package paginations

type Request struct {
	Page    int32
	Counter uint64
	Limit   int32
	Fields  []string
	Values  []string
}
