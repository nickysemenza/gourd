import React from 'react';
import { Button, Form, Input, Label, Segment } from 'semantic-ui-react';

const RecipeEditorIngredientItem = ({
  sectionNum,
  ingredientNum,
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
        onClick={() => this.addIngredient(sectionNum, ingredientNum)}
        content="New"
      />
      <Button
        icon="arrow down"
        onClick={() => this.addIngredient(sectionNum, ingredientNum + 1)}
        content="New"
      />
      <Button
        icon="trash"
        onClick={() => this.deleteIngredient(sectionNum, ingredientNum)}
      />
    </Button.Group>
    <Form>
      <Form.Group>
        <Form.Field width={8}>
          <input
            type="text"
            placeholder="Name"
            value={ingredient.item.name}
            onChange={e => editIngredient(sectionNum, ingredientNum, 'item', e)}
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
              editIngredient(sectionNum, ingredientNum, 'grams', e)
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
              editIngredient(sectionNum, ingredientNum, 'amount', e)
            }
          />
        </Form.Field>
        <Form.Field width={8}>
          <input
            type="text"
            placeholder="Amount Unit"
            value={ingredient.amount_unit}
            onChange={e =>
              editIngredient(sectionNum, ingredientNum, 'amount_unit', e)
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
              editIngredient(sectionNum, ingredientNum, 'modifier', e)
            }
          />
        </Form.Field>
        <Form.Field width={8}>
          <input
            type="text"
            placeholder="Substitute"
            value={ingredient.substitute}
            onChange={e =>
              editIngredient(sectionNum, ingredientNum, 'substitute', e)
            }
          />
        </Form.Field>
      </Form.Group>
    </Form>
  </Segment>
);
export default RecipeEditorIngredientItem;
