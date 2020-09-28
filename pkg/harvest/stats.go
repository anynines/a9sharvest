package harvest

type Stats struct {
	GroupedByTag map[string]float64
}

func NewStats(grouped_by_tags map[string]float64) *Stats {
	return &Stats{
		GroupedByTag: grouped_by_tags,
	}
}
