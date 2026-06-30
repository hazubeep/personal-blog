package main

import (
	"log"
	"net/http"

	"github.com/hazubeep/personal-blog/internal/config"
	"github.com/hazubeep/personal-blog/internal/middleware"
	"github.com/hazubeep/personal-blog/internal/modules/blog"
)

func main() {
	cfg := config.Load()

	repo := blog.NewJSONRepository(cfg.StoragePath)
	service := blog.NewPostService(repo)
	handler := blog.NewPostHandler(service)

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/posts/{id}", handler.ShowPost)

	// Protected routes
	mux.HandleFunc("/admin", middleware.BasicAuth(handler.HomeAdmin))
	mux.HandleFunc("/new", middleware.BasicAuth(handler.CreatePost))
	mux.HandleFunc("/edit/{id}", middleware.BasicAuth(handler.EditPost))
	mux.HandleFunc("/delete/{id}", middleware.BasicAuth(handler.DeletePost))

	log.Printf("server running on http://localhost:%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))

}
