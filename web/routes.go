package web

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"codeberg.org/polarhive/pasta/util"
	"github.com/labstack/echo/v4"
)

const maxFileSize int64 = 1 * 1024 * 1024 // 1 MB in bytes

func registerCrudRoutes(router *echo.Echo, db *util.DB) {
	router.POST("/", handleCreate(db))
	router.GET("/data/:id", handleRead(db))
	router.POST("/update/:id", handleUpdate(db))
	router.POST("/delete/:id", handleDelete(db))
}

func handleCreate(db *util.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("file")
		sitename := "http://localhost:" + os.Getenv("SERVER_PORT")
		if err != nil {
			return c.String(http.StatusBadRequest, "Failed to get file from form")
		}

		src, err := file.Open()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error opening file")
		}
		defer src.Close()

		body, err := io.ReadAll(src)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error reading file data")
		}

		if len(body) == 0 || isEmptyFile(body) {
			return c.String(http.StatusBadRequest, "Empty file")
		}

		if int64(len(body)) > maxFileSize {
			return c.String(http.StatusBadRequest, "File size > 1 MB")
		}

		id := generateRandomID()
		paste := &util.Paste{
			ID:      id,
			Content: string(body),
		}

		if err := db.Create(paste); err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "Error saving data")
		}

		return c.String(http.StatusOK, fmt.Sprintf("%s/data/%s\n", sitename, id))
	}
}

func handleRead(db *util.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		paste, err := db.GetOne(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Paste not found")
		}
		return c.Blob(http.StatusOK, "text/plain; charset=utf-8", []byte(paste.Content))
	}
}

func handleUpdate(db *util.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		file, err := c.FormFile("file")
		sitename := "http://localhost:" + os.Getenv("SERVER_PORT")
		if err != nil {
			return c.String(http.StatusBadRequest, "Failed to get file from form")
		}

		src, err := file.Open()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error opening file")
		}
		defer src.Close()

		body, err := io.ReadAll(src)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error reading file data")
		}

		if len(body) == 0 || isEmptyFile(body) {
			return c.String(http.StatusBadRequest, "Empty file")
		}

		if int64(len(body)) > maxFileSize {
			return c.String(http.StatusBadRequest, "File size > 1 MB")
		}

		paste := &util.Paste{
			ID:      id,
			Content: string(body),
		}

		if err := db.Update(paste); err != nil {
			return c.String(http.StatusInternalServerError, "Error updating data")
		}

		return c.String(http.StatusOK, fmt.Sprintf("Updated paste: %s/data/%s\n", sitename, id))
	}
}

func handleDelete(db *util.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if err := db.Delete(id); err != nil {
			return c.String(http.StatusInternalServerError, "Error deleting paste")
		}
		return c.String(http.StatusOK, "Paste deleted successfully")
	}
}

func isEmptyFile(data []byte) bool {
	for _, b := range data {
		if b != 0 {
			return false
		}
	}
	return true
}

func generateRandomID() string {
	id := make([]byte, 4)
	rand.Read(id)
	return hex.EncodeToString(id)
}
