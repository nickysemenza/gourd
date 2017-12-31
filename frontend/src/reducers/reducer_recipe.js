import {
  REQUEST_RECIPES,
  RECEIVE_RECIPES,
  REQUEST_RECIPE_DETAIL,
  RECEIVE_RECIPE_DETAIL,
  EDIT_TOP_LEVEL_ITEM,
  DELETE_SECTION,
  ADD_SECTION,
  DELETE_INSTRUCTION,
  ADD_INSTRUCTION,
  EDIT_INSTRUCTION,
  DELETE_INGREDIENT,
  ADD_INGREDIENT,
  EDIT_INGREDIENT,
  RECEIVE_IMAGES,
  RECEIVE_MEAL_LIST,
  RECEIVE_CATEGORIES,
  REMOVE_CATEGORY_FROM_RECIPE,
  ADD_CATEGORY_TO_RECIPE
} from '../actions/recipe';

const INITIAL_STATE = {
  recipe_list: [],
  recipe_list_loading: false,

  recipe_detail: {},
  recipe_detail_loading: false,

  image_list: [],

  meal_list: [],
  category_list: []
};

const BLANK_INGREDIENT = {
  item: {
    name: 'name'
  },
  grams: 0,
  amount_unit: 'cup',
  amount: 1,
  substitute: '',
  modifier: '',
  optional: false
};

export default function(state = INITIAL_STATE, action) {
  let r;
  switch (action.type) {
    case REQUEST_RECIPES:
      return { ...state, recipe_list_loading: true };
    case RECEIVE_RECIPES:
      //todo: error checking
      return {
        ...state,
        recipe_list_loading: false,
        recipe_list: action.json
      };
    case REQUEST_RECIPE_DETAIL:
      return { ...state, recipe_detail_loading: true };
    case RECEIVE_RECIPE_DETAIL:
      //todo: error checking
      return {
        ...state,
        recipe_detail_loading: false,
        recipe_detail: {
          ...state.recipe_detail,
          [action.recipeSlug]: {
            ...action.json,
            categories:
              action.json.categories === null ? [] : action.json.categories
          }
        }
      };

    case RECEIVE_IMAGES:
      return {
        ...state,
        image_list: action.json
      };
    case RECEIVE_MEAL_LIST:
      return {
        ...state,
        meal_list: action.json
      };
    case RECEIVE_CATEGORIES:
      return {
        ...state,
        category_list: action.json
      };
    //recipe editing stuff
    case EDIT_TOP_LEVEL_ITEM:
      return {
        ...state,
        recipe_detail: {
          ...state.recipe_detail,
          [action.slug]: {
            ...state.recipe_detail[action.slug],
            [action.fieldName]: action.value
          }
        }
      };
    case DELETE_SECTION:
      return {
        ...state,
        recipe_detail: {
          ...state.recipe_detail,
          [action.slug]: {
            ...state.recipe_detail[action.slug],
            sections: [
              ...state.recipe_detail[action.slug].sections.slice(
                0,
                action.index
              ),
              ...state.recipe_detail[action.slug].sections.slice(
                action.index + 1
              )
            ]
          }
        }
      };
    case ADD_SECTION:
      return {
        ...state,
        recipe_detail: {
          ...state.recipe_detail,
          [action.slug]: {
            ...state.recipe_detail[action.slug],
            sections: [
              ...state.recipe_detail[action.slug].sections.slice(
                0,
                action.index
              ),
              { ingredients: [BLANK_INGREDIENT], instructions: [{ name: '' }] },
              ...state.recipe_detail[action.slug].sections.slice(action.index)
            ]
          }
        }
      };
    case DELETE_INSTRUCTION:
      r = state.recipe_detail[action.slug];
      return {
        ...state,
        recipe_detail: {
          ...state.recipe_detail,
          [action.slug]: {
            ...r,
            sections: [
              ...r.sections.slice(0, action.sectionNum),
              {
                ...r.sections[action.sectionNum],
                instructions: [
                  ...r.sections[action.sectionNum].instructions.slice(
                    0,
                    action.instructionNum
                  ),
                  ...r.sections[action.sectionNum].instructions.slice(
                    action.instructionNum + 1
                  )
                ]
              },
              ...r.sections.slice(action.sectionNum + 1)
            ]
          }
        }
      };
    case ADD_INSTRUCTION:
      r = state.recipe_detail[action.slug];
      return {
        ...state,
        recipe_detail: {
          ...state.recipe_detail,
          [action.slug]: {
            ...r,
            sections: [
              ...r.sections.slice(0, action.sectionNum),
              {
                ...r.sections[action.sectionNum],
                instructions: [
                  ...r.sections[action.sectionNum].instructions.slice(
                    0,
                    action.instructionNum
                  ),
                  { name: '' },
                  ...r.sections[action.sectionNum].instructions.slice(
                    action.instructionNum
                  )
                ]
              },
              ...r.sections.slice(action.sectionNum + 1)
            ]
          }
        }
      };
    case EDIT_INSTRUCTION:
      r = state.recipe_detail[action.slug];
      return {
        ...state,
        recipe_detail: {
          ...state.recipe_detail,
          [action.slug]: {
            ...r,
            sections: [
              ...r.sections.slice(0, action.sectionNum),
              {
                ...r.sections[action.sectionNum],
                instructions: [
                  ...r.sections[action.sectionNum].instructions.slice(
                    0,
                    action.instructionNum
                  ),
                  {
                    ...r.sections[action.sectionNum].instructions[
                      action.instructionNum
                    ],
                    name: action.value
                  },
                  ...r.sections[action.sectionNum].instructions.slice(
                    action.instructionNum + 1
                  )
                ]
              },
              ...r.sections.slice(action.sectionNum + 1)
            ]
          }
        }
      };
    case DELETE_INGREDIENT:
      r = state.recipe_detail[action.slug];
      return {
        ...state,
        recipe_detail: {
          ...state.recipe_detail,
          [action.slug]: {
            ...r,
            sections: [
              ...r.sections.slice(0, action.sectionNum),
              {
                ...r.sections[action.sectionNum],
                ingredients: [
                  ...r.sections[action.sectionNum].ingredients.slice(
                    0,
                    action.ingredientNum
                  ),
                  ...r.sections[action.sectionNum].ingredients.slice(
                    action.ingredientNum + 1
                  )
                ]
              },
              ...r.sections.slice(action.sectionNum + 1)
            ]
          }
        }
      };
    case ADD_INGREDIENT:
      r = state.recipe_detail[action.slug];
      return {
        ...state,
        recipe_detail: {
          ...state.recipe_detail,
          [action.slug]: {
            ...r,
            sections: [
              ...r.sections.slice(0, action.sectionNum),
              {
                ...r.sections[action.sectionNum],
                ingredients: [
                  ...r.sections[action.sectionNum].ingredients.slice(
                    0,
                    action.ingredientNum
                  ),
                  BLANK_INGREDIENT,
                  ...r.sections[action.sectionNum].ingredients.slice(
                    action.ingredientNum
                  )
                ]
              },
              ...r.sections.slice(action.sectionNum + 1)
            ]
          }
        }
      };
    case EDIT_INGREDIENT:
      r = state.recipe_detail[action.slug];
      return {
        ...state,
        recipe_detail: {
          ...state.recipe_detail,
          [action.slug]: {
            ...r,
            sections: [
              ...r.sections.slice(0, action.sectionNum),
              {
                ...r.sections[action.sectionNum],
                ingredients: [
                  ...r.sections[action.sectionNum].ingredients.slice(
                    0,
                    action.ingredientNum
                  ),
                  {
                    ...r.sections[action.sectionNum].ingredients[
                      action.ingredientNum
                    ],
                    [action.field]:
                      action.field === 'item'
                        ? {
                            ...r.sections[action.sectionNum].ingredients[
                              action.ingredientNum
                            ].item,
                            name: action.value
                          }
                        : action.value
                  },
                  ...r.sections[action.sectionNum].ingredients.slice(
                    action.ingredientNum + 1
                  )
                ]
              },
              ...r.sections.slice(action.sectionNum + 1)
            ]
          }
        }
      };
    case REMOVE_CATEGORY_FROM_RECIPE:
      r = state.recipe_detail[action.slug];
      return {
        ...state,
        recipe_detail: {
          ...state.recipe_detail,
          [action.slug]: {
            ...r,
            categories: r.categories.filter(x => x.id !== action.categoryId)
          }
        }
      };
    case ADD_CATEGORY_TO_RECIPE:
      r = state.recipe_detail[action.slug];
      let categoryToAdd = state.category_list.filter(
        x => x.id === action.categoryId
      );
      return {
        ...state,
        recipe_detail: {
          ...state.recipe_detail,
          [action.slug]: {
            ...r,
            categories: [...r.categories, ...categoryToAdd]
          }
        }
      };
    default:
      return state;
  }
}
