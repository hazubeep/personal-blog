package blog

import (
	"fmt"
	"strings"
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

	for _, post := range posts {
		if post.ID == id {
			a := post
			return &a, nil
		}
	}

	return nil, fmt.Errorf("Post with ID %d not found", id)
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
		Content:   strings.TrimSpace(post.Content),
		CreatedAt: post.CreatedAt,
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
