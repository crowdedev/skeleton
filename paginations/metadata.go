package paginations

type Metadata struct {
	Record   int
	Page     int
	Previous int
	Next     int
	Limit    int
	Total    int
}
