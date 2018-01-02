import React from 'react';
import { Button, Form, Input, Label, Segment } from 'semantic-ui-react';

const RecipeEditorIngredientItem = ({
  sectionIndex,
  ingredientIndex,
  ingredient,
  editIngredient,
  addIngredient,
  deleteIngredient
}) => (
  <Segment>
    <Label as="a" color="purple" ribbon>
      {ingredient.item.name}
    </Label>
    <Button.Group>
      <Button
        icon="arrow up"
        onClick={() => addIngredient(sectionIndex, ingredientIndex)}
        content="New"
      />
      <Button
        icon="arrow down"
        onClick={() => addIngredient(sectionIndex, ingredientIndex + 1)}
        content="New"
      />
      <Button
        icon="trash"
        onClick={() => deleteIngredient(sectionIndex, ingredientIndex)}
      />
    </Button.Group>
    <Form>
      <Form.Group>
        <Form.Field width={8}>
          <input
            type="text"
            placeholder="Name"
            value={ingredient.item.name}
            onChange={e =>
              editIngredient(sectionIndex, ingredientIndex, 'item', e)
            }
          />
        </Form.Field>
        <Form.Field width={8}>
          <Input
            label={{ basic: true, content: 'g' }}
            placeholder="Grams"
            labelPosition="right"
            type="number"
            value={ingredient.grams}
            onChange={e =>
              editIngredient(sectionIndex, ingredientIndex, 'grams', e)
            }
          />
        </Form.Field>
      </Form.Group>
      <Form.Group>
        <Form.Field width={8}>
          <input
            type="number"
            placeholder="Amount"
            value={ingredient.amount}
            onChange={e =>
              editIngredient(sectionIndex, ingredientIndex, 'amount', e)
            }
          />
        </Form.Field>
        <Form.Field width={8}>
          <input
            type="text"
            placeholder="Amount Unit"
            value={ingredient.amount_unit}
            onChange={e =>
              editIngredient(sectionIndex, ingredientIndex, 'amount_unit', e)
            }
          />
        </Form.Field>
      </Form.Group>
      <Form.Group>
        <Form.Field width={8}>
          <input
            type="text"
            placeholder="Modifier"
            value={ingredient.modifier}
            onChange={e =>
              editIngredient(sectionIndex, ingredientIndex, 'modifier', e)
            }
          />
        </Form.Field>
        <Form.Field width={8}>
          <input
            type="text"
            placeholder="Substitute"
            value={ingredient.substitute}
            onChange={e =>
              editIngredient(sectionIndex, ingredientIndex, 'substitute', e)
            }
          />
        </Form.Field>
      </Form.Group>
    </Form>
  </Segment>
);
export default RecipeEditorIngredientItem;
