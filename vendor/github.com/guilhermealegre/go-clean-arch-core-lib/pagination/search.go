package pagination

import "fmt"

// Search
type Search struct {
	// Term
	Term string `json:"term"`
}

// IsEmpty verify is search is empty
func (s *Search) IsEmpty() bool {
	return len(s.Term) == 0
}

// ToString
func (s *Search) ToString(concat string) string {
	if s.IsEmpty() {
		return ""
	}
	return fmt.Sprintf("%ssearch=%s", concat, s.Term)
}
