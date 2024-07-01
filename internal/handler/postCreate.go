package handler

import (
	"fmt"
	"forum/internal/models"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (h *Handler) postCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		h.ErrorPage(w, http.StatusNotFound)
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
	categories, err := h.Service.PostSer.GetCategories()
	if err != nil {
		fmt.Println(err.Error())
		h.ErrorPage(w, http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodPost:
		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.Form["categories"]
		img, header, err := r.FormFile("image")
		var nameImage string
		if err != nil {
			if err == http.ErrMissingFile {
				nameImage = ""
			} else {
				fmt.Println(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusBadRequest)
				return
			}
		} else {
			defer img.Close()
			if !isImage(header) {
				http.Error(w, "Uploaded file is not an image", http.StatusUnsupportedMediaType)
				return
			}

			path := filepath.Join("ui", "static", "img", header.Filename)
			f, err := os.Create(path)
			if err != nil {
				fmt.Println("Error saving file:", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			_, err = io.Copy(f, img)
			if err != nil {
				fmt.Println("Error copying file:", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			fileInfo, err := os.Stat(path)
			if err != nil {
				fmt.Println("Error retrieving file info:", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if fileInfo.Size() == 0 {
				fmt.Println("File is empty after saving")
				http.Error(w, "File is empty after saving", http.StatusInternalServerError)
				return
			}
			nameImage = header.Filename

		}

		post := models.Post{
			IdAuth:     user.ID,
			Author:     user.Name,
			Title:      title,
			Content:    content,
			Category:   categories,
			CreateDate: time.Now().Format("January 2, 2006 15:04:05"),
			Img:        nameImage,
		}
		if err := h.Service.CreatePost(post); err != nil {
			fmt.Println(err.Error())
			h.ErrorPage(w, http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	case http.MethodGet:
		if err := h.Tmp.ExecuteTemplate(w, "postCreate.html", categories); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		h.ErrorPage(w, http.StatusMethodNotAllowed)
		return
	}
}

func isImage(header *multipart.FileHeader) bool {
	mimeType := header.Header.Get("Content-Type")
	return strings.HasPrefix(mimeType, "image/")
}
