package models

type Pager struct {
	Offset        int `json:"page"`
	PageSize      int `json:"page_size"`
	RemainingDocs int `json:"remaining_docs"`
	TotalDocs     int `json:"total_docs"`
}

// Creates the pager object
func GetPager(totalDocs, offset, limit int) Pager {
	take := offset * limit

	var remainingDocs int

	if totalDocs-take <= 0 {
		remainingDocs = 0
	} else {
		remainingDocs = totalDocs - take
	}

	return Pager{offset, limit, remainingDocs, totalDocs}
}
