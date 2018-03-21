package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app/model"
	"github.com/nickysemenza/food/backend/app/utils"
	"io"
	"log"
	"os"
	"path"
)

//PutImageUpload uploads images to a recipe based on its Slug
func PutImageUpload(c *gin.Context) {

	db := c.MustGet("DB").(*gorm.DB)
	var finishedImages []model.Image

	//get a ref to the parsed multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(500, err)
	}
	slug := c.PostForm("slug")
	recipe := model.Recipe{}
	if err := db.Where("slug = ?", slug).First(&recipe).Error; err != nil {
		c.JSON(500, errors.New("recipe "+slug+" not found"))
		return
	}

	files := form.File["files"]
	log.Printf("recieving %d images via upload for recipe %s", len(files), slug)
	for i := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			c.JSON(500, err)
			return
		}
		originalFileName := files[i].Filename

		fileData, md5Hash, err := utils.ReadAndHash(file)
		if err != nil {
			c.JSON(500, err)
			return
		}
		//todo: dedup using md5Hash

		//persist an image obj to DB so we get an PK for s3 path
		imageObj := model.Image{}
		imageObj.Md5Hash = md5Hash
		db.Create(&imageObj)
		db.Model(&recipe).Association("Images").Append(&imageObj)

		originalImageSize := model.ImageSize{}
		originalImageSize.IsOriginal = true
		originalImageSize.ImageID = imageObj.ID
		db.Create(&originalImageSize)

		//form filesystem / s3 path
		imagePath := fmt.Sprintf("images/%d%s", imageObj.ID, path.Ext(originalFileName))

		os.MkdirAll("public/images", 0777)
		localImageFile, err := os.Create("public/" + imagePath)
		defer localImageFile.Close()
		if err != nil {
			c.JSON(500, err)
			return
		}

		log.Printf("file: %s -> %s", originalFileName, localImageFile.Name())
		//copy the uploaded file to the destination file
		if _, err := io.Copy(localImageFile, fileData); err != nil {
			c.JSON(500, err)
			return
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
		db.Save(&imageObj)
		//imageObj.MakeSizes()
	}
	c.JSON(200, finishedImages)
}

//GetAllImages gets all images, with their related recipes
func GetAllImages(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var images []model.Image
	db.Preload("Recipes").Find(&images)
	c.JSON(200, images)
}
