package model

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

type Recipe struct {
	Model
	Slug         string    `json:"slug" gorm:"unique"`
	Title        string    `json:"title"`
	TotalMinutes uint      `json:"total_minutes"`
	Equipment    string    `json:"equipment"`
	Source       string    `json:"source"`
	Servings     uint      `json:"servings"`
	Unit         string    `json:"unit"`
	Quantity     uint      `json:"quantity"`
	Sections     []Section `json:"sections"`
}
type Section struct {
	Model
	SortOrder    uint                 `json:"sort_order"`
	Ingredients  []SectionIngredient  `json:"ingredients"`
	Instructions []SectionInstruction `json:"instructions"`
	RecipeID     uint                 `json:"recipe_id"`
	Minutes      uint                 `json:"minutes"`
}
type SectionInstruction struct {
	Model
	Name      string `json:"name"`
	SectionID uint   `json:"section_id"`
}
type SectionIngredient struct {
	Model
	Item       Ingredient `json:"item"`
	ItemID     uint       `json:"item_id"`
	Grams      float32    `json:"grams"`
	Amount     float32    `json:"amount"`
	Unit       string     `json:"amount_unit"`
	Substitute string     `json:"substitute"`
	Modifier   string     `json:"modifier"`
	Optional   bool       `json:"optional"`
	SectionID  uint       `json:"section_id"`
}
type Ingredient struct {
	Model
	Name string `json:"name"`
}

type Project struct {
	Model
	Title    string `gorm:"unique" json:"title"`
	Archived bool   `json:"archived"`
	Tasks    []Task `gorm:"ForeignKey:ProjectID" json:"tasks"`
}

func (p *Project) Archive() {
	p.Archived = true
}

func (p *Project) Restore() {
	p.Archived = false
}

type Task struct {
	Model
	Title     string     `json:"title"`
	Priority  string     `gorm:"type:ENUM('0', '1', '2', '3');default:'0'" json:"priority"`
	Deadline  *time.Time `gorm:"default:null" json:"deadline"`
	Done      bool       `json:"done"`
	ProjectID uint       `json:"project_id"`
}

func (t *Task) Complete() {
	t.Done = true
}

func (t *Task) Undo() {
	t.Done = false
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Project{}, &Task{}, &Section{}, &SectionInstruction{}, &SectionIngredient{}, &Recipe{}, &Ingredient{})
	db.Model(&Section{}).AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")
	db.Model(&SectionInstruction{}).AddForeignKey("section_id", "sections(id)", "RESTRICT", "RESTRICT")
	db.Model(&SectionIngredient{}).AddForeignKey("section_id", "sections(id)", "RESTRICT", "RESTRICT")
	db.Model(&SectionIngredient{}).AddForeignKey("item_id", "ingredients(id)", "RESTRICT", "RESTRICT")
	return db
}
