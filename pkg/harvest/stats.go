package harvest

type Stats struct {
	GroupedByTag    map[string]float64
	percentageByTag map[string]float64
	totalHours      float64
}

func NewStats(grouped_by_tags map[string]float64) *Stats {
	s := Stats{
		GroupedByTag: grouped_by_tags,
	}
	s.prepare()
	return &s
}

func (s *Stats) prepare() {
	s.totalHours = s.calculateTotalHours()

	p := make(map[string]float64)
	for k, v := range s.GroupedByTag {
		p[k] = v * 100.0 / s.totalHours
	}
	s.percentageByTag = p
}

// Returns percentage for total hours.
//
// The hard part is to calculate percentage values so that the summed up
// percentages are not higher than 100%.
// See for example https://revs.runtime-revolution.com/getting-100-with-rounded-percentages-273ffa70252b
// Let's go with a simple solution that can calculate sum(percentages) > 100%.
func (s *Stats) PercentageForTag(tag string) float64 {
	v, ok := s.percentageByTag[tag]
	if ok {
		return v
	}
	return 0.0
}

func (s *Stats) TotalHours() float64 {
	return s.totalHours
}

func (s *Stats) calculateTotalHours() float64 {
	sum := 0.0
	for _, v := range s.GroupedByTag {
		sum += v
	}
	return sum
}
