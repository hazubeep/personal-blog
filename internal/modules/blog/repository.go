package blog

import (
	"encoding/json"
	"os"
	"sync"
)

type JSONRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewJSONRepository(filePath string) *JSONRepository {
	return &JSONRepository{filePath: filePath}
}

func (r *JSONRepository) Load() ([]Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fileData, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Post{}, nil
		}
		return nil, err
	}

	var posts []Post
	if len(fileData) == 0 {
		return posts, nil
	}

	err = json.Unmarshal(fileData, &posts)
	return posts, err
}

func (r *JSONRepository) Save(posts []Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	fileData, err := json.MarshalIndent(posts, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, fileData, 0644)
}
