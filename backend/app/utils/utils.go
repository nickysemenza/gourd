package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nickysemenza/food/backend/app/config"
	"github.com/nickysemenza/food/backend/app/model"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Utils struct {
	*config.Env
}

//Import imports a folder of recipes in json format
func (u Utils) Import(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			raw, err := ioutil.ReadFile(path + f.Name())
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			} else {
				var eachRecipe model.Recipe
				json.Unmarshal(raw, &eachRecipe)
				log.Printf("importing: %s (%s)\n", eachRecipe.Title, eachRecipe.Slug)
				eachRecipe.CreateOrUpdate(u.Env.DB, true)
			}
		}
	}
	log.Printf("Exported %d recipes from %s", len(files), path)
}

//Export exports a folder of recipes in json format
func (u Utils) Export(path string) {
	recipes := []model.Recipe{}
	u.Env.DB.Preload("Sections.Instructions").Preload("Sections.Ingredients.Item").Find(&recipes)
	for _, r := range recipes {
		jsonData, _ := json.Marshal(r)
		err := ioutil.WriteFile(path+r.Slug+".json", jsonData, 0644)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("exporting: %s (%s)\n", r.Title, r.Slug)
		}
	}
	log.Printf("Exported %d recipes to %s", len(recipes), path)
}

func getAWSSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		//LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
}

//AddFileToS3 puts a local file into s3 at a given path.
//Files are public with WRT their ACL.
func AddFileToS3(fileDir string, s3Path string) error {

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

func ReadAndHash(r io.Reader) (io.Reader, string, error) {
	var b bytes.Buffer

	hash := md5.New()
	_, err := io.Copy(&b, io.TeeReader(r, hash))

	if err != nil {
		return nil, "", err
	}

	return &b, hex.EncodeToString(hash.Sum(nil)), nil
}

func GetImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}
