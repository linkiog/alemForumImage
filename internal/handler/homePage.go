package handler

import (
	"fmt"
	"forum/internal/models"
	"net/http"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
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
	posts, err := h.Service.PostSer.GetAllPost()
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
	var info struct {
		AllPosts []models.Post
		User     models.User
		Category []models.Category
	}

	if r.URL.Query().Has("category") {
		category := r.URL.Query().Get("category")
		if Exist(categories, category) {
			post, err := h.Service.PostSer.GetPostsByCategory(category)
			if err != nil {
				fmt.Println(err.Error())
				h.ErrorPage(w, http.StatusInternalServerError)
				return
			}
			info.AllPosts = post
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		info.AllPosts = posts
	}
	info.Category = categories
	info.User = user

	if err := h.Tmp.ExecuteTemplate(w, "homePage.html", info); err != nil {
		fmt.Println(err.Error())
		h.ErrorPage(w, http.StatusInternalServerError)
		return
	}
}

func Exist(categories []models.Category, category string) bool {
	for i := range categories {
		if categories[i].Name == category {
			return true

		}
	}
	return false

}
