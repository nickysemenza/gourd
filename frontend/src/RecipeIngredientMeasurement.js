import React from "react";
export const GramMeasurement = ({ grams, scale = 1 }) => {
  const str = grams === undefined
    ? "\u00a0\u00a0\u00a0\u00a0\u00a0\u00a0\u00a0\u00a0\u00a0"
    : parseFloat((grams * scale).toFixed(2));

  return scale === 1 ? <div>{str}</div> : <b>{str}</b>;
};
export const VolumeMeasurement = ({ i, scale = 1 }) => {
  const str = i.amount === 0
    ? "\u00a0"
    : `${parseFloat((i.amount * scale).toFixed(2))} ${i.amount_unit}`;

  return scale === 1 ? <div>{str}</div> : <b>{str}</b>;
};
