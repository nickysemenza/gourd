package handler

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"github.com/nickysemenza/food/backend/app/model"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func GetAllRecipes(e *Env, w http.ResponseWriter, r *http.Request) error {
	var recipes []model.Recipe
	e.DB.Preload("Images").Preload("Categories").Find(&recipes)
	respondSuccess(w, recipes)
	return nil
}
func ErrorTest(e *Env, w http.ResponseWriter, r *http.Request) error {
	return StatusError{Code: 201, Err: errors.New("sad..")}
}
func GetRecipe(e *Env, w http.ResponseWriter, r *http.Request) error {
	recipe := model.Recipe{}
	slug := mux.Vars(r)["slug"]
	if err := e.DB.Where("slug = ?", slug).Preload("Sections.Instructions").Preload("Sections.Ingredients.Item").Preload("Notes").Preload("Images").Preload("Categories").First(&recipe).Error; err != nil {
		return StatusError{Code: 404, Err: errors.New("recipe " + slug + " not found")}
	}

	respondSuccess(w, recipe)
	return nil
}

func PutRecipe(e *Env, w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	var updatedRecipe model.Recipe
	err := decoder.Decode(&updatedRecipe)
	if err != nil {
		log.Println(err)
	}

	updatedRecipe.CreateOrUpdate(e.DB, false)

	slug := updatedRecipe.Slug
	recipe := model.Recipe{}
	if err := e.DB.Where("slug = ?", slug).Preload("Sections.Instructions").Preload("Sections.Ingredients.Item").Preload("Notes").Preload("Images").Preload("Categories").First(&recipe).Error; err != nil {
		return StatusError{Code: 404, Err: errors.New("recipe " + slug + " not found")}
	}

	respondSuccess(w, recipe)
	return nil
}

func CreateRecipe(e *Env, w http.ResponseWriter, r *http.Request) error {
	//decode the data from JSON encoded request body
	decoder := json.NewDecoder(r.Body)
	var parsed struct {
		Slug  string `json:"slug"`
		Title string `json:"title"`
	}
	err := decoder.Decode(&parsed)
	if err != nil {
		log.Println(err)
	}

	//see if one exists
	recipe := model.Recipe{}
	if !e.DB.Where("slug = ?", parsed.Slug).First(&recipe).RecordNotFound() {
		respondError(w, 500, "slug exists already")
		return nil
	}
	recipe.Slug = parsed.Slug
	recipe.Title = parsed.Title
	e.DB.Save(&recipe)
	respondSuccess(w, "added!")
	return nil
}

func AddNote(e *Env, w http.ResponseWriter, r *http.Request) error {
	//find the recipe we are adding a note to
	recipe := model.Recipe{}
	slug := mux.Vars(r)["slug"]
	if err := e.DB.Where("slug = ?", slug).First(&recipe).Error; err != nil {
		return StatusError{Code: 404, Err: errors.New("recipe " + slug + " not found")}
	}

	//decode the note from JSON encoded request body
	decoder := json.NewDecoder(r.Body)
	var parsed struct {
		Note string `json:"note"`
	}
	err := decoder.Decode(&parsed)
	if err != nil {
		log.Println(err)
	}
	//add a new RecipeNote Model, save it
	note := model.RecipeNote{
		Body:     parsed.Note,
		RecipeID: recipe.ID,
	}
	e.DB.Save(&note)

	respondSuccess(w, note)
	return nil
}

func readAndHash(a io.Reader) (io.Reader, string, error) {
	var b bytes.Buffer

	hash := md5.New()
	_, err := io.Copy(&b, io.TeeReader(a, hash))

	if err != nil {
		return nil, "", err
	}

	return &b, hex.EncodeToString(hash.Sum(nil)), nil
}
func ImageUploadTest(e *Env, w http.ResponseWriter, r *http.Request) error {

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

		fileData, md5Hash, err := readAndHash(file)
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

			if err := addFileToS3(localImageFile.Name(), imagePath); err != nil {
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

func GetAllImages(e *Env, w http.ResponseWriter, r *http.Request) error {
	var images []model.Image
	e.DB.Preload("Recipes").Find(&images)
	respondSuccess(w, images)
	return nil
}

func GetAllMeals(e *Env, w http.ResponseWriter, r *http.Request) error {
	var meals []model.Meal
	e.DB.Preload("RecipeMeal.Recipe").Find(&meals)
	respondSuccess(w, meals)
	return nil
}
func GetAllCategories(e *Env, w http.ResponseWriter, r *http.Request) error {
	var categories []model.Category
	e.DB.Find(&categories)
	respondSuccess(w, categories)
	return nil
}

func getAWSSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		//LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
}

func addFileToS3(fileDir string, s3Path string) error {

	s, err := getAWSSession()
	if err != nil {
		log.Fatal(err)
	}

	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	log.Printf("saving file to S3 at %s", s3Path)
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("S3_BUCKET")),
		Key:           aws.String(s3Path),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(size),
		ACL:           aws.String("public-read"),
		ContentType:   aws.String(http.DetectContentType(buffer)),
	})
	return err
}
