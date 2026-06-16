package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"log"

	"forum/back-end/internal/models"
	"forum/back-end/internal/render"
)

type CreationPageData struct {
	Categories   []models.Category
	ErrorMessage string
	Content      string
}

func HandlerCreation(db *sql.DB, v *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/post/create" {
			http.NotFound(w, r)
			return
		}

		user, err := models.GetCurrentUser(db, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if r.Method == http.MethodGet {
			renderCreationPage(w, db, v, "", "")
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 25*1024*1024)

		err = r.ParseMultipartForm(25 * 1024 * 1024)
		if err != nil {
			renderCreationPage( w, db, v, "Le formulaire est trop lourd.", "", )
			return
		}

		content := strings.TrimSpace(r.FormValue("content"))
		if content == "" {
			renderCreationPage( w, db, v, "Le contenu du post ne peut pas être vide.", content, )
			return
		}

		title := strings.TrimSpace(r.FormValue("title"))
		if title == "" {
			renderCreationPage( w, db, v, "Le titre du post ne peut pas être vide.", title, )
			return
		}

		categoryIDs, err := parseCategoryIDs(r)
		if err != nil {
			renderCreationPage( w, db, v, "Catégorie invalide.", content, )
			return
		}

		imageLink, err := savePostImage(r)
		if err != nil {
			renderCreationPage( w, db, v, err.Error(), content, )
			return
		}

		err = models.CreatePost( db, user.ID, title, content, imageLink, categoryIDs, )
		if err != nil {
			log.Printf("Erreur CreatePost : %v", err)
			renderCreationPage( w, db, v, "Erreur lors de la création du post : "+err.Error(), content, )
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}


func renderCreationPage( w http.ResponseWriter, db *sql.DB, v *render.Render, errorMessage string, content string, ) {
	categories, err := models.GetCategories(db)
	if err != nil {
		http.Error( w, "Erreur chargement catégories", http.StatusInternalServerError, )
		return
	}

	data := CreationPageData{
		Categories:   categories,
		ErrorMessage: errorMessage,
		Content:      content,
	}

	v.Render(w, "creationPage.html", data)
}

func parseCategoryIDs(r *http.Request) ([]int, error) {
	values := r.Form["categories"]

	categoryIDs := []int{}
	alreadyAdded := map[int]bool{}

	for _, value := range values {
		categoryID, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}

		if !alreadyAdded[categoryID] {
			categoryIDs = append(categoryIDs, categoryID)
			alreadyAdded[categoryID] = true
		}
	}

	return categoryIDs, nil
}

func savePostImage(r *http.Request) (string, error) {
	file, header, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'upload de l'image")
	}
	defer file.Close()

	if header.Size > 20*1024*1024 {
		return "", fmt.Errorf("l'image ne doit pas dépasser 20 Mo")
	}

	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("image invalide")
	}

	contentType := http.DetectContentType(buffer[:n])

	var extension string

	switch contentType {
	case "image/jpeg":
		extension = ".jpg"
	case "image/png":
		extension = ".png"
	case "image/gif":
		extension = ".gif"
	default:
		return "", fmt.Errorf("format d'image non autorisé")
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("erreur lecture image")
	}

	fileName, err := randomFileName(extension)
	if err != nil {
		return "", fmt.Errorf("erreur génération nom image")
	}

	uploadDir := "front-end/static/uploads"

	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		return "", fmt.Errorf("erreur création dossier upload")
	}

	filePath := filepath.Join(uploadDir, fileName)

	destination, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("erreur sauvegarde image")
	}
	defer destination.Close()

	_, err = io.Copy(destination, file)
	if err != nil {
		return "", fmt.Errorf("erreur copie image")
	}

	return "/static/uploads/" + fileName, nil
}

func randomFileName(extension string) (string, error) {
	bytes := make([]byte, 16)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes) + extension, nil
}