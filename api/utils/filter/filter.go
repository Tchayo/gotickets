package filter

// Filter : paginator filter helper
type Filter struct {
	Page    int
	Limit   int
	OrderBy []string
	Search  string
}
