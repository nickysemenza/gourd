export type Override = {
  sectionID: number;
  ingredientID: number;
  value: number;
  attr: IngredientAttr;
};

export type RecipeTweaks = {
  override?: Override;
  multiplier: number;
  edit: boolean;
};

export type IngredientAttr = "grams" | "amount";
export type IngredientKind = "recipe" | "ingredient";
