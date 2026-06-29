package main

import (
	"log"
	"net/http"

	"github.com/hazubeep/personal-blog/internal/middleware"
	"github.com/hazubeep/personal-blog/internal/modules/blog"
)

const port = ":8080"

func main() {

	repo := blog.NewJSONRepository("./data/posts.json")
	service := blog.NewPostService(repo)
	handler := blog.NewPostHandler(service)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// Public routes
	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/article/{id}", handler.ShowPost)

	// Protected routes
	mux.HandleFunc("/admin", middleware.BasicAuth(handler.HomeAdmin))
	mux.HandleFunc("/new", middleware.BasicAuth(handler.CreatePost))
	mux.HandleFunc("/edit/{id}", middleware.BasicAuth(handler.EditPost))
	mux.HandleFunc("/delete/{id}", middleware.BasicAuth(handler.DeletePost))

	log.Printf("server running on localhost:%s", port)
	log.Fatal(http.ListenAndServe(port, mux))

}
