import React from 'react'
export const GramMeasurement = ({grams, scale=1}) => {
    const str = (grams === undefined)
        ? "\u00a0"
        : parseFloat((grams * scale).toFixed(2));

    return scale===1
        ? <div>{str}</div>
        : <b>{str}</b>;
};
export const VolumeMeasurement = ({measurement, scale=1}) => {
    const str = (measurement === undefined)
        ? "\u00a0"
        : `${parseFloat((measurement.amount * scale).toFixed(2))} ${measurement.unit}`;

    return scale===1
        ? <div>{str}</div>
        : <b>{str}</b>;
};
