import { hasGrams, hasVolume } from './RecipeIngredientMeasurement';
it('hasGrams behaves', () => {
  expect(hasGrams({ grams: 4 })).toBeTruthy();
  expect(hasGrams({ amount: 4 })).toBeFalsy();
});
it('hasVolume behaves', () => {
  expect(hasVolume({ grams: 4 })).toBeFalsy();
  expect(hasVolume({ amount: 4 })).toBeTruthy();
});
