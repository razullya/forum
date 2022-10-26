package delivery

import (
	"forum/models"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) userProfilePage(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/profile/")

	userI := h.userIdentity(w, r)
	// if userI == (models.User{}) {
	// 	h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	// 	return
	// }
	if userI.Username == username {

		postsUser, err := h.Services.Post.GetPostsByUsername(username)
		if err != nil {
			h.errorPage(w, http.StatusBadRequest, err.Error())
		}
		info := models.Info{
			ThatUser: userI,
			Posts:    postsUser,
		}
		if err := h.Tmpl.ExecuteTemplate(w, "user.html", info); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

	} else {

		user, err := h.Services.User.GetUserByUsername(username)
		if err != nil {
			log.Println(err)
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}
		postsUser, err := h.Services.Post.GetPostsByUsername(username)
		if err != nil {
			h.errorPage(w, http.StatusBadRequest, err.Error())
		}
		info := models.Info{
			User:     user,
			Posts:    postsUser,
			ThatUser: userI,
		}
		if err := h.Tmpl.ExecuteTemplate(w, "user.html", info); err != nil {
			log.Println(err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

	}
}
