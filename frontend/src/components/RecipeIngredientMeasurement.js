import React from 'react';
export const GramMeasurement = ({ grams, scale = 1 }) => {
  const str =
    grams === 0
      ? '\u00a0\u00a0\u00a0\u00a0\u00a0\u00a0\u00a0\u00a0\u00a0'
      : scaleAndTruncate(grams, scale);

  return scale === 1 ? <div>{str}</div> : <b>{str}</b>;
};
export const VolumeMeasurement = ({ i, scale = 1 }) => {
  const str =
    i.amount === 0
      ? '\u00a0'
      : `${scaleAndTruncate(i.amount, scale)} ${i.amount_unit}`;

  return scale === 1 ? <div>{str}</div> : <b>{str}</b>;
};

const scaleAndTruncate = (val, scale) => parseFloat((val * scale).toFixed(2));

//todo: markup `modifier` display with italics/small/etc
export const IngredientNameDisplay = ({ i }) =>
  `${i.item.name}${i.modifier ? `, ${i.modifier}` : ''}`;

export const getScaledMeasurementString = (ingredient, scale = 1) => {
  let { amount, amount_unit, grams, modifier } = ingredient;
  return `${scaleAndTruncate(
    amount,
    scale
  )} ${amount_unit} or ${scaleAndTruncate(grams, scale)} grams (${modifier})`;
};

export const hasGrams = i => i.grams !== undefined && i.grams !== 0;
export const hasVolume = i => i.amount !== undefined && i.amount !== 0;
