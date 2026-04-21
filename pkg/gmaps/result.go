package gmaps

// Entry represents a single Google Maps business listing result.
type Entry struct {
	Title           string   `json:"title"`
	Category        string   `json:"category"`
	Address         string   `json:"address"`
	OpenHours       string   `json:"open_hours"`
	Website         string   `json:"website"`
	Phone           string   `json:"phone"`
	PlusCode        string   `json:"plus_code"`
	ReviewCount     int      `json:"review_count"`
	Rating          float64  `json:"rating"`
	Latitude        float64  `json:"latitude"`
	Longitude       float64  `json:"longitude"`
	CID             string   `json:"cid"`
	CompleteAddress Address  `json:"complete_address"`
	Images          []string `json:"images,omitempty"`
}

// Address holds structured address components.
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`
	Country    string `json:"country"`
}

// IsValid returns true if the entry has at minimum a title.
func (e *Entry) IsValid() bool {
	return e.Title != ""
}

// Merge copies non-zero fields from other into e.
func (e *Entry) Merge(other *Entry) {
	if other == nil {
		return
	}
	if e.Website == "" {
		e.Website = other.Website
	}
	if e.Phone == "" {
		e.Phone = other.Phone
	}
	if e.OpenHours == "" {
		e.OpenHours = other.OpenHours
	}
}
