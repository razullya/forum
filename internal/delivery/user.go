package delivery

import (
	"fmt"
	"forum/models"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) userProfilePage(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/profile/")

	userI := h.userIdentity(w, r)
	if userI == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if userI.Username == username {
		fmt.Println(username)

		postsUser, err := h.Services.Post.GetPostsByUsername(username)
		if err != nil {
			h.errorPage(w, http.StatusBadRequest, err.Error())
		}
		info := models.Info{
			User:  userI,
			Posts: postsUser,
		}
		if err := h.Tmpl.ExecuteTemplate(w, "user.html", info); err != nil {
			log.Println(err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

	} else {
		fmt.Println(username)
		user, err := h.Services.User.GetUserByUsername(username)
		if err != nil {
			log.Println(err)
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}

		info := models.Info{
			User: user,
		}
		if err := h.Tmpl.ExecuteTemplate(w, "user.html", info); err != nil {
			log.Println(err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

	}
}
