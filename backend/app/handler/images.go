package handler

import (
	"errors"
	"fmt"
	"github.com/nickysemenza/food/backend/app/config"
	"github.com/nickysemenza/food/backend/app/model"
	"github.com/nickysemenza/food/backend/app/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

//PutImageUpload uploads images to a recipe based on its Slug
func PutImageUpload(e *config.Env, w http.ResponseWriter, r *http.Request) error {

	var finishedImages []model.Image
	err := r.ParseMultipartForm(100000)
	if err != nil {
		return StatusError{Code: 500, Err: err}
	}

	//get a ref to the parsed multipart form
	m := r.MultipartForm
	slug := m.Value["slug"][0]
	recipe := model.Recipe{}
	if err := e.DB.Where("slug = ?", slug).First(&recipe).Error; err != nil {
		return StatusError{Code: 404, Err: errors.New("recipe " + slug + " not found")}
	}

	files := m.File["file"]
	log.Printf("recieving %d images via upload for recipe %s", len(files), slug)
	for i := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			return StatusError{Code: 500, Err: err}
		}
		originalFileName := files[i].Filename

		fileData, md5Hash, err := utils.ReadAndHash(file)
		if err != nil {
			return StatusError{Code: 500, Err: err}
		}
		//todo: dedup using md5Hash

		//persist an image obj to DB so we get an PK for s3 path
		imageObj := model.Image{}
		imageObj.Md5Hash = md5Hash
		e.DB.Create(&imageObj)
		e.DB.Model(&recipe).Association("Images").Append(&imageObj)

		//form filesystem / s3 path
		imagePath := fmt.Sprintf("images/%d%s", imageObj.ID, path.Ext(originalFileName))

		os.MkdirAll("public/images", 0777)
		localImageFile, err := os.Create("public/" + imagePath)
		log.Printf("file: %s -> %s", originalFileName, localImageFile.Name())

		defer localImageFile.Close()
		if err != nil {
			return StatusError{Code: 500, Err: err}
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(localImageFile, fileData); err != nil {
			return StatusError{Code: 500, Err: err}
		}

		if os.Getenv("S3_IMAGES") == "true" {

			if err := utils.AddFileToS3(localImageFile.Name(), imagePath); err != nil {
				imageObj.IsInS3 = false
				log.Println(err)
			} else {
				imageObj.IsInS3 = true
				finishedImages = append(finishedImages, imageObj)
			}
		} else {
			finishedImages = append(finishedImages, imageObj)
		}
		imageObj.OriginalFileName = originalFileName
		imageObj.Path = imagePath
		e.DB.Save(&imageObj)
	}
	respondSuccess(w, finishedImages)
	return nil
}

//GetAllImages gets all images, with their related recipes
func GetAllImages(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	var images []model.Image
	e.DB.Preload("Recipes").Find(&images)
	respondSuccess(w, images)
	return nil
}
