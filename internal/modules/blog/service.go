package blog

import (
	"fmt"
	"time"
)

type PostService struct {
	repo *JSONRepository
}

func NewPostService(repo *JSONRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Index() ([]Post, error) {
	posts, err := s.repo.Load()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostService) Show(id int) (*Post, error) {
	posts, err := s.repo.Load()
	if err != nil {
		return nil, err
	}

	for _, article := range posts {
		if article.ID == id {
			a := article
			return &a, nil
		}
	}

	return nil, fmt.Errorf("Article with ID %d not found", id)
}

func (s *PostService) Create(post Post) error {
	posts, err := s.repo.Load()
	if err != nil {
		return err
	}

	lastID := 0
	for _, a := range posts {
		if a.ID > lastID {
			lastID = a.ID
		}
	}

	newPost := Post{
		ID:        lastID + 1,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: time.Now(),
	}

	posts = append(posts, newPost)
	return s.repo.Save(posts)
}

func (s *PostService) Edit(id int, post Post) error {
	posts, err := s.repo.Load()
	if err != nil {
		return err
	}

	for i, a := range posts {
		if a.ID == id {
			posts[i].Title = post.Title
			posts[i].Content = post.Content
			posts[i].CreatedAt = post.CreatedAt
		}
	}

	return s.repo.Save(posts)
}

func (s *PostService) Delete(id int) error {
	posts, err := s.repo.Load()
	if err != nil {
		return err
	}

	filtered := []Post{}
	for _, a := range posts {
		if a.ID != id {
			filtered = append(filtered, a)
		}
	}

	return s.repo.Save(filtered)
}
