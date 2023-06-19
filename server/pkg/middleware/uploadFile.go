package middleware

import (
	"io"
	"io/ioutil"
	dto "mytask/dto/result"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

func UploadFile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//Mengambil file dari form dengan nama "image" dari request context menggunakan c.FormFile("image")
		file, err := c.FormFile("image")
		if file != nil {
			if err != nil {
				return c.JSON(http.StatusBadRequest, "failed get images")
			}

			// handle extensi filename
			ext := filepath.Ext(file.Filename)
			allowedExts := []string{".png", ".jpg", ".jpeg"} // Ekstensi yang diperbolehkan

			// check valid extensi
			validExt := false
			for _, allowedExt := range allowedExts {
				if strings.ToLower(ext) == allowedExt {
					validExt = true
					break
				}
			}

			if !validExt {
				return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: "Invalid file extension"})
			}
			//open file
			src, err := file.Open()
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			defer src.Close()
			//Membuat file sementara dengan nama unik di direktori "uploads" menggunakan ioutil.TempFile(). Nama file akan dimulai dengan "image-" dan diikuti dengan karakter acak.
			tempFile, err := ioutil.TempFile("uploads", "image-*"+ext)

			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			defer tempFile.Close()

			if _, err = io.Copy(tempFile, src); err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}

			data := tempFile.Name()

			// update := strings.Split(data, "\\")[1]
			// fmt.Println(update)

			c.Set("dataFile", data)
			return next(c)
		}
		c.Set("dataFile", "")
		return next(c)
	}
}
