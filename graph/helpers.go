package graph

import (
	"github.com/nickysemenza/food/graph/model"
	"github.com/nickysemenza/food/manager"
)

func fromRecipe(res *manager.Recipe) *model.Recipe {
	mr := model.Recipe{
		UUID:         res.UUID,
		Name:         res.Name,
		TotalMinutes: int(res.TotalMinutes),
		Unit:         res.Unit}

	// for _, x := range res.Sections {
	// 	s := &model.Section{Minutes: int(x.Minutes)}
	// 	for _, i := range x.Instructions {
	// 		s.Instructions = append(s.Instructions, &model.Instruction{
	// 			Instruction: i.Instruction,
	// 		})
	// 	}
	// 	for _, i := range x.Ingredients {
	// 		s.Ingredients = append(s.Ingredients, &model.SectionIngredient{
	// 			Grams: i.Grams,
	// 			Name:  i.Name,
	// 			// Amount:    i.Amount.Float64,
	// 			// Unit:      i.Unit.String,
	// 			// Adjective: i.Adjective.String,
	// 			// Optional:  i.Optional.Bool,
	// 		})
	// 	}
	// 	mr.Sections = append(mr.Sections, s)
	// }
	return &mr
}
