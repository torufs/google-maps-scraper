package gmaps

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONLWriter_Write(t *testing.T) {
	f, err := os.CreateTemp("", "*.jsonl")
	require.NoError(t, err)
	path := f.Name()
	f.Close()
	defer os.Remove(path)

	w, err := NewJSONLWriter(path)
	require.NoError(t, err)

	entry := &Entry{Title: "Test Place", Category: "Restaurant"}
	require.NoError(t, w.Write(entry))
	require.NoError(t, w.Close())

	data, err := os.Open(path)
	require.NoError(t, err)
	defer data.Close()

	var got Entry
	require.NoError(t, json.NewDecoder(data).Decode(&got))
	assert.Equal(t, "Test Place", got.Title)
	assert.Equal(t, "Restaurant", got.Category)
}

func TestJSONLWriter_WriteNil(t *testing.T) {
	f, err := os.CreateTemp("", "*.jsonl")
	require.NoError(t, err)
	path := f.Name()
	f.Close()
	defer os.Remove(path)

	w, err := NewJSONLWriter(path)
	require.NoError(t, err)

	// Writing nil should be a no-op and leave the file empty
	require.NoError(t, w.Write(nil))
	require.NoError(t, w.Close())

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, int64(0), info.Size())
}

func TestJSONLWriter_MultipleEntries(t *testing.T) {
	f, err := os.CreateTemp("", "*.jsonl")
	require.NoError(t, err)
	path := f.Name()
	f.Close()
	defer os.Remove(path)

	w, err := NewJSONLWriter(path)
	require.NoError(t, err)

	entries := []*Entry{
		{Title: "Place A", Category: "Cafe"},
		{Title: "Place B", Category: "Bar"},
		{Title: "Place C", Category: "Hotel"},
	}

	for _, e := range entries {
		require.NoError(t, w.Write(e))
	}
	require.NoError(t, w.Close())

	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var results []Entry
	for scanner.Scan() {
		var e Entry
		require.NoError(t, json.Unmarshal(scanner.Bytes(), &e))
		results = append(results, e)
	}
	require.NoError(t, scanner.Err())

	require.Len(t, results, 3)
	assert.Equal(t, "Place A", results[0].Title)
	assert.Equal(t, "Place B", results[1].Title)
	assert.Equal(t, "Place C", results[2].Title)
}

func TestNewJSONLWriter_EmptyPath(t *testing.T) {
	// An empty path should return an error rather than silently failing
	_, err := NewJSONLWriter("")
	assert.Error(t, err)
}
