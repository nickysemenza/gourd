package manager

import (
	"github.com/nickysemenza/food/db"
	"gopkg.in/guregu/null.v3/zero"
)

// Recipe is the presentation for a Recipe
type Recipe struct {
	UUID         string    `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	Name         string    `json:"name,omitempty" yaml:"name,omitempty"`
	TotalMinutes int64     `json:"total_minutes,omitempty" yaml:"total_minutes,omitempty"`
	Equipment    string    `json:"equipment,omitempty" yaml:"equipment,omitempty"`
	Source       string    `json:"source,omitempty" yaml:"source,omitempty"`
	Quantity     int64     `json:"quantity,omitempty" yaml:"quantity,omitempty"`
	Servings     int64     `json:"servings,omitempty" yaml:"servings,omitempty"`
	Unit         string    `json:"unit,omitempty" yaml:"unit,omitempty"`
	Sections     []Section `json:"sections,omitempty" yaml:"sections,omitempty"`
}

// Section is the presentation for a Sections
type Section struct {
	Minutes      int64               `json:"minutes,omitempty" yaml:"minutes,omitempty"`
	Ingredients  []SectionIngredient `json:"ingredients,omitempty" yaml:"ingredients,omitempty"`
	Instructions []Instruction       `json:"instructions,omitempty" yaml:"instructions,omitempty"`
}

// SectionIngredient is the presentation for a Ingredient
type SectionIngredient struct {
	Name      string  `json:"name,omitempty" yaml:"name,omitempty"`
	Grams     float64 `json:"grams,omitempty" yaml:"grams,omitempty"`
	Amount    float64 `json:"amount,omitempty" yaml:"amount,omitempty"`
	Unit      string  `json:"unit,omitempty" yaml:"unit,omitempty"`
	Adjective string  `json:"adjective,omitempty" yaml:"adjective,omitempty"`
	Optional  bool    `json:"optional,omitempty" yaml:"optional,omitempty"`
}

// Instruction is the presentation for a Instruction
type Instruction struct {
	Instruction string `json:"instruction,omitempty" yaml:"instruction,omitempty"`
}

func (r *Recipe) toDB() *db.Recipe {
	if r == nil {
		return nil
	}
	dbr := db.Recipe{
		UUID:         r.UUID,
		Name:         r.Name,
		TotalMinutes: zero.IntFrom(r.TotalMinutes),
		Equipment:    zero.StringFrom(r.Equipment),
		Source:       zero.StringFrom(r.Source),
		Servings:     zero.IntFrom(r.Servings),
		Quantity:     zero.IntFrom(r.Quantity),
		Unit:         zero.StringFrom(r.Unit),
	}
	for _, s := range r.Sections {
		dbs := db.Section{
			Minutes: zero.IntFrom(s.Minutes),
		}
		for _, i := range s.Instructions {
			dbs.Instructions = append(dbs.Instructions, db.SectionInstruction{
				Instruction: i.Instruction,
			})
		}
		for _, i := range s.Ingredients {
			dbs.Ingredients = append(dbs.Ingredients, db.SectionIngredient{
				// Name:      i.Name,
				Grams:     zero.FloatFrom(i.Grams),
				Amount:    zero.FloatFrom(i.Amount),
				Unit:      zero.StringFrom(i.Unit),
				Adjective: zero.StringFrom(i.Adjective),
				Optional:  zero.BoolFrom(i.Optional),
			})
		}

		dbr.Sections = append(dbr.Sections, dbs)
	}
	return &dbr
}

func FromRecipe(dbr *db.Recipe) *Recipe {
	if dbr == nil {
		return nil
	}
	r := Recipe{
		UUID:         dbr.UUID,
		Name:         dbr.Name,
		TotalMinutes: dbr.TotalMinutes.Int64,
		Equipment:    dbr.Equipment.String,
		Source:       dbr.Source.String,
		Servings:     dbr.Servings.Int64,
		Quantity:     dbr.Quantity.Int64,
		Unit:         dbr.Unit.String,
	}
	for _, dbs := range dbr.Sections {
		s := Section{
			Minutes: dbs.Minutes.Int64,
		}
		for _, i := range dbs.Instructions {
			s.Instructions = append(s.Instructions, Instruction{
				Instruction: i.Instruction,
			})
		}
		for _, i := range dbs.Ingredients {
			s.Ingredients = append(s.Ingredients, SectionIngredient{
				Grams: i.Grams.Float64,
				// Name:      i.Name,
				Amount:    i.Amount.Float64,
				Unit:      i.Unit.String,
				Adjective: i.Adjective.String,
				Optional:  i.Optional.Bool,
			})
		}

		r.Sections = append(r.Sections, s)
	}
	return &r
}
