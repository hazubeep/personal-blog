package blog

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	service *PostService
}

func NewPostHandler(service *PostService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	articles, err := h.service.Index()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, articles)
}

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID is not valid", http.StatusBadRequest)
		return
	}

	article, err := h.service.Show(id)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/article.html"))
	tmpl.Execute(w, article)
}

// Handler for private route
func (h *Handler) HomeAdmin(w http.ResponseWriter, r *http.Request) {
	articles, err := h.service.Index()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin.html"))
	tmpl.Execute(w, articles)
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles("templates/add-form.html"))
		tmpl.Execute(w, nil)
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

func (h *Handler) EditPost(w http.ResponseWriter, r *http.Request) {
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

		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.ParseFiles("templates/form.html"))
		tmpl.Execute(w, post)

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

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID is not valid", http.StatusBadRequest)
		return
	}
	h.service.Delete(id)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
