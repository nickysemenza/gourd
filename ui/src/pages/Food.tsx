import React from "react";
import { useGetFoodById } from "../api/openapi-hooks/api";
import Debug from "../components/Debug";
import { Code } from "../util";

const Food: React.FC = () => {
  // const [food, setFood]
  return (
    <div className="grid grid-cols-5 gap-4">
      <FoodInfo fdc_id={9999999} />
      <FoodInfo fdc_id={171047} />
      <FoodInfo fdc_id={392941} />
      <FoodInfo fdc_id={747448} />
      <FoodInfo fdc_id={789097} />
      <FoodInfo fdc_id={335560} />
    </div>
  );
};
export default Food;

const FoodInfo: React.FC<{ fdc_id: number }> = ({ fdc_id }) => {
  const { data } = useGetFoodById({ fdc_id });
  if (!data) return null;
  const { category, branded_info, portions, ...food } = data;
  return (
    <div>
      <h2 className="font-bold text-l">{food.description}</h2>
      <div>
        <Code>{food.data_type}</Code>
      </div>
      <div>
        <Code>{fdc_id}</Code>
      </div>
      <Debug data={{ category, branded_info, portions }} />

      <table>
        {food.nutrients
          .filter((n) => n.amount > 0.1)
          .map((n) => (
            <tr>
              <td>{n.amount.toFixed(2)} </td>
              <td>{n.nutrient.unit_name} </td>
              <td>{n.nutrient.name} </td>
            </tr>
          ))}
      </table>
    </div>
  );
};
