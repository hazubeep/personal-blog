package blog

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PostHandler struct {
	service   *PostService
	tmplCache map[string]*template.Template
}

func NewPostHandler(service *PostService) *PostHandler {
	cache := make(map[string]*template.Template)

	pages := []string{
		"home.html",
		"admin.html",
		"post.html",
		"form.html",
		"add-form.html",
	}

	for _, page := range pages {
		ts, err := template.ParseFiles("templates/" + page)
		if err != nil {
			log.Fatalf("Error parsing template %s: %v", page, err)
		}

		cache[page] = ts
	}

	return &PostHandler{
		service:   service,
		tmplCache: cache,
	}
}

func (h *PostHandler) render(w http.ResponseWriter, page string, data any) {
	tmpl, ok := h.tmplCache[page]
	if !ok {
		http.Error(w, "Template "+page+" not found in cache", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render page: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *PostHandler) Home(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.Index()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.render(w, "home.html", posts)
}

func (h *PostHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID is not valid", http.StatusBadRequest)
		return
	}

	posts, err := h.service.Show(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	h.render(w, "post.html", posts)
}

// Handler for private route
func (h *PostHandler) HomeAdmin(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.Index()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.render(w, "admin.html", posts)
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		h.render(w, "add-form.html", nil)
	case http.MethodPost:
		title := r.FormValue("title")
		createdAt := r.FormValue("createdAt")
		content := r.FormValue("content")

		layout := "2006-01-02"

		parsedTime, err := time.Parse(layout, createdAt)
		if err != nil {
			http.Error(w, "Failed formating date", http.StatusBadRequest)
			return
		}

		err = h.service.Create(Post{Title: title, CreatedAt: parsedTime, Content: content})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func (h *PostHandler) EditPost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID is not valid", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		post, err := h.service.Show(id)
		if err != nil {
			http.Error(w, "ID not found", http.StatusNotFound)
			return
		}

		h.render(w, "form.html", post)

	case http.MethodPost:
		title := r.FormValue("title")
		createdAt := r.FormValue("createdAt")
		content := r.FormValue("content")

		layout := "2006-01-02"

		parsedTime, err := time.Parse(layout, createdAt)
		if err != nil {
			http.Error(w, "Failed formating date", http.StatusBadRequest)
			return
		}

		err = h.service.Edit(id, Post{Title: title, CreatedAt: parsedTime, Content: content})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID is not valid", http.StatusBadRequest)
		return
	}
	h.service.Delete(id)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
