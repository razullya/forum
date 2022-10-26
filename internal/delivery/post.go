package delivery

import (
	"fmt"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/" {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	post, err := h.Services.Post.GetPostById(id)
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}
	comments, err := h.Services.GetCommentsByIdPost(id)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	for i := 0; i < len(comments); i++ {
		comments[i].Likes, comments[i].Dislikes, err = h.Services.Reaction.GetCounts(comments[i].Id, "comment")
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	user := h.userIdentity(w, r)
	if err := r.ParseForm(); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	post.Likes, post.Dislikes, err = h.Services.Reaction.GetCounts(post.Id, "post")
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.Services.Post.UpdateCountsReactionsPost(post.Likes, post.Dislikes, post.Id); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	switch r.Method {
	case http.MethodGet:

		info := models.Info{
			Post:     post,
			Comments: comments,
			ThatUser: user,
		}

		if err := h.Tmpl.ExecuteTemplate(w, "post.html", info); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		user := h.userIdentity(w, r)
		comment, ok := r.Form["comment"]
		if comment[0] == "" {
			h.errorPage(w, http.StatusBadRequest, "comment field not found")
			return
		}

		if !ok {
			h.errorPage(w, http.StatusBadRequest, "comment field not found")
			return
		}

		if err := h.Services.Comment.CreateComment(id, user, comment[0]); err != nil {
			h.errorPage(w, http.StatusBadRequest, "comment field not found")
			return
		}

		http.Redirect(w, r, r.URL.Path+fmt.Sprintf("/?id=%d", id), http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

}
func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/post/create" {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	user := h.userIdentity(w, r)
	if err := r.ParseForm(); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	info := models.Info{
		ThatUser: user,
	}
	switch r.Method {
	case http.MethodGet:
		if err := h.Tmpl.ExecuteTemplate(w, "createPost.html", info); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:

		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		title, ok := r.Form["title"]

		if !ok || title[0] == "" {
			h.errorPage(w, http.StatusBadRequest, "title field not found")
			return
		}
		description, ok := r.Form["description"]

		if !ok || description[0] == "" {
			h.errorPage(w, http.StatusBadRequest, "description field not found")
			return
		}
		category, ok := r.Form["category"]

		if !ok || category[0] == "" {
			h.errorPage(w, http.StatusBadRequest, "category field not found")
			return
		}

		post := models.Post{
			Title:       title[0],
			Description: description[0],
			Category:    category,
		}

		user := h.userIdentity(w, r)

		if user == (models.User{}) {
			h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		if err := h.Services.Post.CreatePost(post, user); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			//todo
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}
func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/delete" {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	user := h.userIdentity(w, r)
	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	postId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/delete/"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}

	if err := h.Services.DeletePost(postId, user); err != nil {
		h.errorPage(w, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) changePost(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/post/change" {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	switch r.Method {
	case http.MethodGet:
		if err := h.Tmpl.ExecuteTemplate(w, "createPost.html", nil); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:

		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		title, ok := r.Form["title"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, "title field not found")
			return
		}
		description, ok := r.Form["description"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, "description field not found")
			return
		}
		post := models.Post{
			Title:       title[0],
			Description: description[0],
		}

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		user := h.userIdentity(w, r)
		if user == (models.User{}) {
			h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		if err := h.Services.Post.UpdatePost(id, post, user); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			//todo
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *Handler) filterPostCategories(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/categories/" {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := h.userIdentity(w, r)
	if err := r.ParseForm(); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	categories, err := h.Services.Post.GetAllCategories()
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
	posts, err := h.Services.Post.GetAllPosts("")
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	info := models.Info{
		Categories: categories,
		Posts:      posts,
		ThatUser:   user,
	}
	if err := h.Tmpl.ExecuteTemplate(w, "categories.html", info); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

}
