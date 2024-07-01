package handler

import (
	"fmt"
	"forum/internal/models"
	"net/http"
)

func (h *Handler) myPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/myPosts" {
		h.ErrorPage(w, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		h.ErrorPage(w, http.StatusMethodNotAllowed)
		return

	}
	userValue := r.Context().Value("user")
	if userValue == nil {
		h.ErrorPage(w, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(models.User)
	if !ok {
		h.ErrorPage(w, http.StatusUnauthorized)
		return
	}
	if !user.IsAuth {
		h.ErrorPage(w, http.StatusUnauthorized)
		return

	}
	myPost, err := h.Service.PostSer.GetMyPosts(user.ID)
	if err != nil {
		fmt.Println(err.Error())
		h.ErrorPage(w, http.StatusInternalServerError)
		return
	}
	categories, err := h.Service.PostSer.Category()
	if err != nil {
		fmt.Println(err.Error())
		h.ErrorPage(w, http.StatusInternalServerError)
		return
	}
	info := struct {
		AllPosts []models.Post
		User     models.User
		Category []models.Category
	}{
		AllPosts: myPost,
		User:     user,
		Category: categories,
	}

	if err := h.Tmp.ExecuteTemplate(w, "homePage.html", info); err != nil {
		fmt.Println(err.Error())
		h.ErrorPage(w, http.StatusInternalServerError)
		return
	}

}
