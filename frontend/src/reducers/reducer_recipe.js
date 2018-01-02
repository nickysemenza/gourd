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
  MOVE_INSTRUCTION,
  DELETE_INGREDIENT,
  ADD_INGREDIENT,
  EDIT_INGREDIENT,
  RECEIVE_IMAGES,
  RECEIVE_MEAL_LIST,
  RECEIVE_CATEGORIES,
  REMOVE_CATEGORY_FROM_RECIPE,
  ADD_CATEGORY_TO_RECIPE
} from '../actions/recipe';

import update from 'immutability-helper';
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

const BLANK_INSTRUCTION = { name: '' };
const BLANK_SECTION = {
  ingredients: [BLANK_INGREDIENT],
  instructions: [BLANK_INSTRUCTION]
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
      return update(state, {
        recipe_detail: {
          [action.slug]: { sections: { $splice: [[action.index, 1]] } }
        }
      });
    case ADD_SECTION:
      return update(state, {
        recipe_detail: {
          [action.slug]: {
            sections: { $splice: [[action.index, 0, BLANK_SECTION]] }
          }
        }
      });
    case DELETE_INSTRUCTION:
      return update(state, {
        recipe_detail: {
          [action.slug]: {
            sections: {
              [action.sectionIndex]: {
                instructions: { $splice: [[action.instructionIndex, 1]] }
              }
            }
          }
        }
      });
    case ADD_INSTRUCTION:
      return update(state, {
        recipe_detail: {
          [action.slug]: {
            sections: {
              [action.sectionIndex]: {
                instructions: {
                  $splice: [[action.instructionIndex, 0, BLANK_INSTRUCTION]]
                }
              }
            }
          }
        }
      });
    case EDIT_INSTRUCTION:
      return update(state, {
        recipe_detail: {
          [action.slug]: {
            sections: {
              [action.sectionIndex]: {
                instructions: {
                  [action.instructionIndex]: { name: { $set: action.value } }
                }
              }
            }
          }
        }
      });
    case MOVE_INSTRUCTION:
      let movingInstruction =
        state.recipe_detail[action.slug].sections[action.sectionIndex]
          .instructions[action.instructionIndex];
      return update(state, {
        recipe_detail: {
          [action.slug]: {
            sections: {
              [action.sectionIndex]: {
                instructions: {
                  $splice: [
                    [action.instructionIndex, 1],
                    [action.targetIndex, 0, movingInstruction]
                  ]
                }
              }
            }
          }
        }
      });
    case DELETE_INGREDIENT:
      return update(state, {
        recipe_detail: {
          [action.slug]: {
            sections: {
              [action.sectionIndex]: {
                ingredients: { $splice: [[action.ingredientIndex, 1]] }
              }
            }
          }
        }
      });
    case ADD_INGREDIENT:
      r = state.recipe_detail[action.slug];
      return update(state, {
        recipe_detail: {
          [action.slug]: {
            sections: {
              [action.sectionIndex]: {
                ingredients: {
                  $splice: [[action.ingredientIndex, 0, BLANK_INGREDIENT]]
                }
              }
            }
          }
        }
      });
    case EDIT_INGREDIENT:
      //TODO: cleanup
      r = state.recipe_detail[action.slug];
      let thisSectionIngredient = {
        ...r.sections[action.sectionIndex].ingredients[action.ingredientIndex],
        [action.field]:
          action.field === 'item'
            ? {
                ...r.sections[action.sectionIndex].ingredients[
                  action.ingredientIndex
                ].item,
                name: action.value
              }
            : action.value
      };
      return update(state, {
        recipe_detail: {
          [action.slug]: {
            sections: {
              [action.sectionIndex]: {
                ingredients: {
                  [action.ingredientIndex]: { $set: thisSectionIngredient }
                }
              }
            }
          }
        }
      });
    case REMOVE_CATEGORY_FROM_RECIPE:
      return update(state, {
        recipe_detail: {
          [action.slug]: {
            categories: {
              $set: state.recipe_detail[action.slug].categories.filter(
                x => x.id !== action.categoryId
              )
            }
          }
        }
      });
    case ADD_CATEGORY_TO_RECIPE:
      return update(state, {
        recipe_detail: {
          [action.slug]: {
            categories: {
              $push: state.category_list.filter(x => x.id === action.categoryId)
            }
          }
        }
      });
    default:
      return state;
  }
}
