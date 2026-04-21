package gmaps_test

import (
	"testing"

	"github.com/gosom/google-maps-scraper/pkg/gmaps"
)

func TestEntry_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		entry gmaps.Entry
		want  bool
	}{
		{"valid entry", gmaps.Entry{Title: "Acme Corp"}, true},
		{"empty entry", gmaps.Entry{}, false},
		{"no title but has phone", gmaps.Entry{Phone: "+1-800-000"}, false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.entry.IsValid(); got != tc.want {
				t.Errorf("IsValid() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestEntry_Merge(t *testing.T) {
	t.Parallel()

	base := &gmaps.Entry{Title: "Base", Phone: "+1-000"}
	other := &gmaps.Entry{Website: "https://example.com", Phone: "+1-999", OpenHours: "9-5"}

	base.Merge(other)

	if base.Website != "https://example.com" {
		t.Errorf("expected website to be merged, got %q", base.Website)
	}
	if base.Phone != "+1-000" {
		t.Errorf("expected phone to remain unchanged, got %q", base.Phone)
	}
	if base.OpenHours != "9-5" {
		t.Errorf("expected open hours to be merged, got %q", base.OpenHours)
	}
}

func TestEntry_Merge_NilOther(t *testing.T) {
	t.Parallel()

	entry := &gmaps.Entry{Title: "Test"}
	// should not panic
	entry.Merge(nil)
}
