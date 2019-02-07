package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/nickysemenza/food/app/model"
	opentracing "github.com/opentracing/opentracing-go"
)

//PutImageUpload uploads images to a recipe based on its Slug
func PutImageUpload(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	db := model.GetDBFromContext(ctx)

	span, ctx := opentracing.StartSpanFromContext(ctx, "image upload handler")
	defer span.Finish()
	span.LogEvent("begin")

	var finishedImages []model.Image

	//get a ref to the parsed multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(500, err)
	}
	slug := c.PostForm("slug")
	span.SetTag("recipe-slug", slug)

	recipe := model.Recipe{}
	if err := db.Where("slug = ?", slug).First(&recipe).Error; err != nil {
		c.JSON(500, errors.New("recipe "+slug+" not found"))
		span.LogEvent("slug not found")
		return
	}

	files := form.File["file"]
	log.Printf("recieving %d images via upload for recipe %s", len(files), slug)
	span.SetTag("image-count", len(files))
	span.SetTag("recipe-id", recipe.ID)

	for i := range files {
		fSpan, ctx := opentracing.StartSpanFromContext(ctx, "image upload process")
		db := model.GetDBFromContext(ctx)
		originalFileName := files[i].Filename
		fSpan.LogEvent("begin")
		fSpan.SetTag("filename", originalFileName)
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			c.JSON(500, err)
			return
		}

		fileData, md5Hash, err := ReadAndHash(file)
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

		fSpan.LogEvent("assigned image-id")
		fSpan.SetTag("image-id", imageObj.ID)

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

		uploadToS3 := os.Getenv("S3_IMAGES") == "true"

		fSpan.SetTag("use-s3", uploadToS3)
		if uploadToS3 {
			if err := AddFileToS3(ctx, localImageFile.Name(), imagePath); err != nil {
				fSpan.LogEvent("s3 upload error")
				imageObj.IsInS3 = false
				log.Println(err)
			} else {
				fSpan.LogEventWithPayload("uploaded to s3", imagePath)
				imageObj.IsInS3 = true
				finishedImages = append(finishedImages, imageObj)
			}
		} else {
			finishedImages = append(finishedImages, imageObj)
		}
		imageObj.OriginalFileName = originalFileName
		imageObj.Path = imagePath
		db.Save(&imageObj)
		fSpan.LogEvent("finished")
		fSpan.Finish()
		//imageObj.MakeSizes()
	}
	c.JSON(200, finishedImages)
}

//GetAllImages gets all images, with their related recipes
func GetAllImages(c *gin.Context) {
	db := model.GetDBFromContext(c.MustGet("ctx").(context.Context))
	var images []model.Image
	db.Preload("Recipes").Find(&images)
	c.JSON(200, images)
}
