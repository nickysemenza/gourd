import React from 'react';
import { Form } from 'semantic-ui-react';

const RecipeEditorBasicInfo = ({ recipe, editTopLevelItem }) => (
  <Form>
    <Form.Group>
      <Form.Field width={16}>
        <label>Title</label>
        <input
          type="text"
          value={recipe.title}
          onChange={e => editTopLevelItem('title', e)}
        />
      </Form.Field>
      <Form.Field width={16}>
        <label>Source</label>
        <input
          type="text"
          value={recipe.source}
          onChange={e => editTopLevelItem('source', e)}
        />
      </Form.Field>
      <Form.Field width={16}>
        <label>Quantity</label>
        <input
          type="number"
          value={recipe.quantity}
          onChange={e => editTopLevelItem('quantity', e)}
        />
      </Form.Field>
      <Form.Field width={16}>
        <label>Quantity Unit</label>
        <input
          type="text"
          value={recipe.unit}
          onChange={e => editTopLevelItem('unit', e)}
        />
      </Form.Field>
      <Form.Field width={16}>
        <label>Servings</label>
        <input
          type="number"
          value={recipe.servings}
          onChange={e => editTopLevelItem('servings', e)}
        />
      </Form.Field>
      <Form.Field width={16}>
        <label>Total Minutes</label>
        <input
          type="number"
          value={recipe.total_minutes}
          onChange={e => editTopLevelItem.bind('total_minutes', e)}
        />
      </Form.Field>
    </Form.Group>
  </Form>
);
export default RecipeEditorBasicInfo;
