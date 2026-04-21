package gmaps

// Entry holds the scraped data for a single Google Maps result.
type Entry struct {
	Title    string  `json:"title"`
	Address  string  `json:"address"`
	Phone    string  `json:"phone"`
	Website  string  `json:"website"`
	Rating   float64 `json:"rating"`
	Reviews  int     `json:"reviews"`
	Category string  `json:"category"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

// IsValid reports whether the entry contains at least a title.
func (e *Entry) IsValid() bool {
	if e == nil {
		return false
	}
	return e.Title != ""
}

// Merge copies non-zero fields from other into e.
// Fields already set on e are not overwritten.
func (e *Entry) Merge(other *Entry) {
	if other == nil {
		return
	}
	if e.Title == "" {
		e.Title = other.Title
	}
	if e.Address == "" {
		e.Address = other.Address
	}
	if e.Phone == "" {
		e.Phone = other.Phone
	}
	if e.Website == "" {
		e.Website = other.Website
	}
	if e.Rating == 0 {
		e.Rating = other.Rating
	}
	if e.Reviews == 0 {
		e.Reviews = other.Reviews
	}
	if e.Category == "" {
		e.Category = other.Category
	}
	if e.Lat == 0 {
		e.Lat = other.Lat
	}
	if e.Lon == 0 {
		e.Lon = other.Lon
	}
}
