package model

import (
	"time"

	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"os"
)

//Model is what all the models are based off of
type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

//User is the `user` table
type User struct {
	Model
	Email      string `json:"email" gorm:"unique"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	FacebookID string `json:"fb_id"`
	Admin      bool   `json:"admin"`
}

//GetJWTToken returns a JWT token for a User
func (u *User) GetJWTToken(db *gorm.DB) string {

	mySigningKey := []byte("AllYourBase")

	// Create the Claims
	claims := &jwt.StandardClaims{
		//ExpiresAt: 15000,
		Issuer: "test",
		Id:     fmt.Sprint(u.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Println(err)
	}

	return ss
}

//Recipe is the main Recipe object. `recipes` in SQL
type Recipe struct {
	Model
	Slug         string       `json:"slug" gorm:"unique"`
	Title        string       `json:"title"`
	TotalMinutes uint         `json:"total_minutes"`
	Equipment    string       `json:"equipment"`
	Source       string       `json:"source"`
	Servings     uint         `json:"servings"`
	Unit         string       `json:"unit"`
	Quantity     uint         `json:"quantity"`
	Sections     []Section    `json:"sections"`
	Notes        []RecipeNote `json:"notes"`
	Images       []Image      `json:"images" gorm:"many2many:recipe_images;"`
	Categories   []Category   `json:"categories" gorm:"many2many:recipe_categories;"`
	RecipeMeal   []RecipeMeal `json:"recipe_meals"`
}

//Section is `sections`
type Section struct {
	Model
	SortOrder    uint                 `json:"sort_order"`
	Ingredients  []SectionIngredient  `json:"ingredients"`
	Instructions []SectionInstruction `json:"instructions"`
	RecipeID     uint                 `json:"recipe_id"`
	Minutes      uint                 `json:"minutes"`
}

//SectionInstruction is `section_instructions`
type SectionInstruction struct {
	Model
	Name      string `json:"name"`
	SortOrder uint   `json:"sort_order"`
	SectionID uint   `json:"section_id"`
}

//SectionIngredient is `section_ingredients`
type SectionIngredient struct {
	Model
	Item       Ingredient `json:"item"`
	SortOrder  uint       `json:"sort_order"`
	ItemID     uint       `json:"item_id"`
	Grams      float32    `json:"grams"`
	Amount     float32    `json:"amount"`
	Unit       string     `json:"amount_unit"`
	Substitute string     `json:"substitute"`
	Modifier   string     `json:"modifier"`
	Optional   bool       `json:"optional"`
	SectionID  uint       `json:"section_id"`
}

func (r *Recipe) reIndexSectionSortOrder() {
	for x := range r.Sections {
		r.Sections[x].SortOrder = uint(x)
	}
}
func (s *Section) reIndexSectionContentsSortOrder() {
	for x := range s.Ingredients {
		s.Ingredients[x].SortOrder = uint(x)
	}
	for x := range s.Instructions {
		s.Instructions[x].SortOrder = uint(x)
	}
}

//Ingredient is `ingredients`
type Ingredient struct {
	Model
	Name string `json:"name"`
}

//RecipeNote is `recipe_note`
type RecipeNote struct {
	Model
	Body     string `json:"body" sql:"type:text"`
	RecipeID uint   `json:"recipe_id"`
}

//Image is `images`
type Image struct {
	Model
	Path             string   `json:"path"`
	OriginalFileName string   `json:"original_name"`
	IsInS3           bool     `json:"in_s3"`
	Recipes          []Recipe `json:"recipes" gorm:"many2many:recipe_images;"`
	Md5Hash          string   `json:"md5"`
}

//Category is `categories`
type Category struct {
	Model
	Name    string   `json:"name"`
	Recipes []Recipe `json:"recipes" gorm:"many2many:recipe_categories;"`
}

//Meal is `meals`
type Meal struct {
	Model
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Description string       `json:"description"`
	RecipeMeal  []RecipeMeal `json:"recipe_meals"`
	Time        time.Time    `json:"time"`
}

//RecipeMeal is `recipe_meal`
type RecipeMeal struct {
	Recipe     Recipe `json:"recipe"`
	RecipeID   uint
	Meal       Meal `json:"meal"`
	MealID     uint
	Multiplier uint `json:"multiplier" gorm:"default:1"`
}

//MarshalJSON is a helper to marshal an Image.
//It returns s3 or local url depending on if the image exists in s3
//useful for local dev with large images!
func (i *Image) MarshalJSON() ([]byte, error) {
	type Alias Image
	var url string
	if i.IsInS3 {
		url = "https://" + os.Getenv("S3_BUCKET") + ".s3.amazonaws.com/" + i.Path
	} else {
		url = os.Getenv("API_PUBLIC_URL") + "/public/" + i.Path
	}
	return json.Marshal(&struct {
		FullURL string `json:"url"`
		*Alias
	}{
		FullURL: url,
		Alias:   (*Alias)(i),
	})
}

//AddToMeal adds the recipe to a Meal, with specified multiplier
func (r *Recipe) AddToMeal(db *gorm.DB, meal *Meal, multiplier uint) {
	recipemeal := RecipeMeal{}
	recipemeal.Multiplier = multiplier
	recipemeal.RecipeID = r.ID
	recipemeal.MealID = meal.ID
	db.Create(&recipemeal)
}

//FindOrCreateUsingName finds an ingredient by name, otherwise it creates it
func (ingredient *Ingredient) FindOrCreateUsingName(db *gorm.DB) {
	db.FirstOrCreate(&ingredient, Ingredient{Name: ingredient.Name})
}

//GetFresh gets a fresh DB instance of an Ingredient
func (ingredient *Ingredient) GetFresh(db *gorm.DB) {
	db.First(&ingredient, ingredient.ID)
}

//AfterCreate is a hook on Ingredient's gorm saving
func (ingredient *Ingredient) AfterCreate() (err error) {
	log.Printf("[ingredient] created: %s, %d", ingredient.Name, ingredient.ID)
	return
}

//CreateOrUpdate updates or creates a Recipe with its children.
//	recursivelyStripIDs is useful for importing, when you want to ignore the primary keys of a json obj
func (r Recipe) CreateOrUpdate(db *gorm.DB, recursivelyStripIDs bool) {
	//todo: ensure that we aren't overwriting something with same slug, by checking for presence of ID
	//update Categories and Sections relations
	db.Model(&r).Association("Categories").Replace(r.Categories)
	db.Model(&r).Association("Sections").Replace(r.Sections)
	for x := range r.Sections {
		eachSection := &r.Sections[x]
		if recursivelyStripIDs {
			eachSection.ID = 0
		}
		for y := range eachSection.Ingredients {
			eachSectionIngredient := &eachSection.Ingredients[y]
			eachItem := &eachSectionIngredient.Item
			if recursivelyStripIDs {
				eachSectionIngredient.ID = 0
				eachItem.ID = 0
			}
			if eachItem.ID == 0 {
				//	new ingredient!
				//	find by name, to see if we have existing
				eachItem.FindOrCreateUsingName(db)
				//eachItem = model.GetIngredientByName(e.DB, eachItem.Name)
				log.Printf("[ingredient] %s does not have an ID, giving it %d: ", eachItem.Name, eachItem.ID)
			} else {
				//	get fresh obj via eachIngredient.ID
				fresh := *eachItem
				fresh.GetFresh(db)
				//	if eachIngredient.Name != fresh.Name IT WAS MUTATED AAH!
				if eachItem.Name != fresh.Name {
					log.Printf("[ingredient] name of #%d was mutated! (%s->%s)", eachItem.ID, eachItem.Name, fresh.Name)
					//we want to preserve the original eachItem; create new w/ eachItem.Name

					// find by name, or create new
					newItem := Ingredient{Name: eachItem.Name}
					newItem.FindOrCreateUsingName(db)
					eachSectionIngredient.Item = newItem
				}
			}

		}
		for y := range eachSection.Instructions {
			eachSectionInstruction := &eachSection.Instructions[y]
			if recursivelyStripIDs {
				eachSectionInstruction.ID = 0
			}
		}
		db.Model(&eachSection).Association("Ingredients").Replace(eachSection.Ingredients)
		db.Model(&eachSection).Association("Instructions").Replace(eachSection.Instructions)

		eachSection.reIndexSectionContentsSortOrder()
	}
	if recursivelyStripIDs {
		r.ID = 0
	}

	r.reIndexSectionSortOrder()

	if err := db.Save(&r).Error; err != nil {
		log.Println(err)
	}
}

//GetFromSlug will populate a Recipe obj based on slug
func (r *Recipe) GetFromSlug(db *gorm.DB, slug string) error {
	return db.Where("slug = ?", slug).
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
		First(&r).Error
}

var modelsInOrder = []interface{}{
	&Recipe{},
	&Ingredient{},
	&Image{},
	&RecipeNote{},
	&Meal{},
	&RecipeMeal{},
	&Section{},
	&SectionIngredient{},
	&SectionInstruction{},
	&Category{},
	&User{},
}

//DBMigrate migrates the DB
func DBMigrate(db *gorm.DB) *gorm.DB {
	for y := range modelsInOrder {
		db.AutoMigrate(modelsInOrder[y])
	}
	db.Model(&Section{}).AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")
	db.Model(&RecipeNote{}).AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")
	db.Model(&SectionInstruction{}).AddForeignKey("section_id", "sections(id)", "RESTRICT", "RESTRICT")
	db.Model(&SectionIngredient{}).AddForeignKey("section_id", "sections(id)", "RESTRICT", "RESTRICT")
	db.Model(&SectionIngredient{}).AddForeignKey("item_id", "ingredients(id)", "RESTRICT", "RESTRICT")

	db.Table("recipe_categories").AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")
	db.Table("recipe_categories").AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")

	db.Table("recipe_images").AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")
	db.Table("recipe_images").AddForeignKey("image_id", "images(id)", "RESTRICT", "RESTRICT")

	db.Table("recipe_meals").AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")
	db.Table("recipe_meals").AddForeignKey("meal_id", "meals(id)", "RESTRICT", "RESTRICT")
	return db
}

//DBReset resets the DB.
func DBReset(db *gorm.DB) *gorm.DB {
	db.Exec("SET foreign_key_checks = 0;")
	for i := len(modelsInOrder) - 1; i >= 0; i-- {
		db.DropTable(modelsInOrder[i])
	}
	db.Exec("SET foreign_key_checks = 1;")
	return db
}
