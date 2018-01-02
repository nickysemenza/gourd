package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app/config"
	"github.com/nickysemenza/food/backend/app/model"
	"github.com/nickysemenza/food/backend/app/utils"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func GetAllRecipes(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	var recipes []model.Recipe
	e.DB.Preload("Images").Preload("Categories").Find(&recipes)
	respondSuccess(w, recipes)
	return nil
}
func ErrorTest(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	return StatusError{Code: 201, Err: errors.New("sad..")}
}
func GetRecipe(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	recipe := model.Recipe{}
	slug := mux.Vars(r)["slug"]

	err := e.DB.Where("slug = ?", slug).
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("sections.sort_order ASC")
		}).
		Preload("Sections.Instructions", func(db *gorm.DB) *gorm.DB {
			return db.Order("section_instructions.sort_order ASC")
		}).
		Preload("Sections.Ingredients", func(db *gorm.DB) *gorm.DB {
			return db.Order("section_ingredients.sort_order ASC")
		}).
		Preload("Sections.Ingredients.Item").
		Preload("Notes").
		Preload("Images").
		Preload("Categories").
		First(&recipe).Error
	if err != nil {
		return StatusError{Code: 404, Err: errors.New("recipe " + slug + " not found")}
	}

	respondSuccess(w, recipe)
	return nil
}

func PutRecipe(e *config.Env, w http.ResponseWriter, r *http.Request) error {
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

func CreateRecipe(e *config.Env, w http.ResponseWriter, r *http.Request) error {
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

func AddNote(e *config.Env, w http.ResponseWriter, r *http.Request) error {
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
func GetAllImages(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	var images []model.Image
	e.DB.Preload("Recipes").Find(&images)
	respondSuccess(w, images)
	return nil
}

func GetAllMeals(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	var meals []model.Meal
	e.DB.Preload("RecipeMeal.Recipe").Find(&meals)
	respondSuccess(w, meals)
	return nil
}
func GetAllCategories(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	var categories []model.Category
	e.DB.Find(&categories)
	respondSuccess(w, categories)
	return nil
}
