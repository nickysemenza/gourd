import React from 'react';
import {
  GramMeasurement,
  VolumeMeasurement
} from './RecipeIngredientMeasurement';
const RecipeTable = ({ recipe, scale }) => {
  let totalWeight = 0;
  let totalMinutes = 0;
  let instructionNum = 1;
  let tableRows = recipe.sections.map((recipeSection, num) => {
    let ingredientList = [],
      weightList = [],
      volumeList = [];

    let minutes = recipeSection.minutes ? recipeSection.minutes : 0;
    totalMinutes += minutes;
    let ingredients = recipeSection.ingredients
      ? recipeSection.ingredients
      : [];
    ingredients.forEach((i, n) => {
      let ingredientName = `${i.item.name}${
        i.modifier ? `, ${i.modifier}` : ''
      }`;
      ingredientList.push(
        <div className="ingredientCellItem" key={n}>
          {ingredientName}
        </div>
      );
      if (i.grams) totalWeight += i.grams;
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

    let instructionList = recipeSection.instructions.map((i, n) => (
      <div key={n}>
        <b>{`${instructionNum++}`}. </b>
        {i.name}
      </div>
    ));

    let test = [];
    let x = 0;
    for (x = 0; x < ingredientList.length; x++) {
      test.push(
        <tr key={`il-${x}`}>
          <td style={{ width: '45%' }}>{ingredientList[x]}</td>
          <td>{weightList[x]}</td>
          <td>{volumeList[x]}</td>
        </tr>
      );
    }

    return (
      <tr key={`r-${num}`}>
        <td style={{ verticalAlign: 'middle' }}>
          <b>{String.fromCharCode(num + 65)}.</b>
          <br />
          {minutes ? `${minutes} min` : null}
        </td>
        <td colSpan="3" style={{ verticalAlign: 'middle' }}>
          <table style={{ width: '100%' }} className="table borderless">
            <tbody>{test}</tbody>
          </table>
        </td>
        <td>{instructionList}</td>
      </tr>
    );
  });

  return (
    <table className="table table-sm borderheavy recipeTable">
      <thead style={{ backgroundColor: '#9a9da8' }}>
        <tr>
          <th width={'50px'}>&nbsp;</th>
          <th>ingredients</th>
          <th>weight</th>
          <th>volume</th>
          <th>steps</th>
        </tr>
      </thead>
      <tbody>{tableRows}</tbody>
    </table>
  );
};

export default RecipeTable;
