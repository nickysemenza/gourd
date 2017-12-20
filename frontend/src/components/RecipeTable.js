import React from 'react';
import {
  GramMeasurement,
  IngredientNameDisplay,
  VolumeMeasurement
} from './RecipeIngredientMeasurement';
const RecipeTable = ({ recipe, scale }) => {
  let instructionNum = 1;
  let tableRows = recipe.sections.map((recipeSection, num) => {
    let ingredientList = [],
      weightList = [],
      volumeList = [];
    const { ingredients } = recipeSection;
    if (!ingredients) return [];

    ingredients.forEach((i, n) => {
      ingredientList.push(
        <div className="ingredientCellItem" key={n}>
          <IngredientNameDisplay i={i} />
        </div>
      );
      weightList.push(
        <div className="ingredientCellItem" key={n}>
          <GramMeasurement grams={i.grams} scale={scale} />
        </div>
      );
      volumeList.push(
        <div className="ingredientCellItem" key={n}>
          <VolumeMeasurement i={i} scale={scale} />
        </div>
      );
    });
    return (
      <tr key={`r-${num}`}>
        <td style={{ verticalAlign: 'middle' }}>
          <b>{String.fromCharCode(num + 65)}.</b>
          <br />
          {recipeSection.minutes ? `${recipeSection.minutes} min` : null}
        </td>
        <td>
          {ingredientList.map((item, key) => {
            return <span key={key}>{item}</span>;
          })}
        </td>
        <td>
          {weightList.map((item, key) => {
            return <span key={key}>{item}</span>;
          })}
        </td>
        <td>
          {volumeList.map((item, key) => {
            return <span key={key}>{item}</span>;
          })}
        </td>
        <td>
          {recipeSection.instructions.map((i, n) => (
            <div key={n}>
              <b>{`${instructionNum++}`}. </b> {i.name}
            </div>
          ))}
        </td>
      </tr>
    );
  });

  return (
    <table className="borderheavy recipeTable">
      <thead style={{ backgroundColor: '#9a9da8' }}>
        <tr style={{ textAlign: 'left' }}>
          <th className="tableHeadTH">&nbsp;</th>
          <th className="tableHeadTH">ingredients</th>
          <th className="tableHeadTH">weight</th>
          <th className="tableHeadTH">volume</th>
          <th className="tableHeadTH">steps</th>
        </tr>
      </thead>
      <tbody>{tableRows}</tbody>
    </table>
  );
};

export default RecipeTable;
