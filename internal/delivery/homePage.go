package delivery

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
	"net/http"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
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

	filter, _ := r.Form["filter"]
	if len(filter) == 0 {
		filter = append(filter, "")
	}
	if filter[0] != "by_categories" && filter[0] != "by_likes" && filter[0] != "by_time" && filter[0] != "" {
		h.errorPage(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	fmt.Println(filter[0])
	posts, err := h.Services.Post.GetAllPosts(filter[0])
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	for i := 0; i < len(posts); i++ {
		posts[i].Likes, posts[i].Dislikes, err = h.Services.Reaction.GetCounts(posts[i].Id, "post")
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := h.Services.Post.UpdateCountsReactionsPost(posts[i].Likes, posts[i].Dislikes, posts[i].Id); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	info := models.Info{
		User:  user,
		Posts: posts,
	}
	if err := h.Tmpl.ExecuteTemplate(w, "homepage.html", info); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
}
