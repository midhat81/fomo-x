package checkpoint

import (
	"encoding/json"
	"os"
)

// Store persists the last processed signature to a local file so the
// ingestor can resume from where it left off after a restart. This is a
// simple file-based checkpoint for Day 1 — can be moved to Postgres/Redis
// later if needed.
type Store struct {
	path string
}

type state struct {
	LastSignature string `json:"last_signature"`
}

// NewStore creates a checkpoint store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the last processed signature to disk.
func (s *Store) Save(signature string) error {
	data, err := json.Marshal(state{LastSignature: signature})
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// Load reads the last processed signature from disk. Returns an empty
// string if no checkpoint exists yet.
func (s *Store) Load() (string, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	var st state
	if err := json.Unmarshal(data, &st); err != nil {
		return "", err
	}
	return st.LastSignature, nil
}