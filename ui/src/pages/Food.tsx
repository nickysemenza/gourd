import React from "react";
import { useGetFoodById } from "../api/openapi-hooks/api";
import Debug from "../components/Debug";
import { UnitMappingList } from "../components/Misc";
import { Code, scaledRound } from "../util";

const Food: React.FC = () => {
  // const [food, setFood]
  return (
    <div className="grid grid-cols-5 gap-4">
      {/* <FoodInfo fdc_id={9999999} /> */}
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
  const {
    category,
    branded_info,
    portions,
    description,
    data_type,
    nutrients,
  } = data.wrapper;
  const { unit_mappings } = data;
  return (
    <div>
      <h2 className="font-bold text-l">{description}</h2>
      <div>
        <Code>{data_type}</Code>
      </div>
      <div>
        <Code>{fdc_id}</Code>
      </div>
      <UnitMappingList unit_mappings={unit_mappings} />
      <Debug data={{ category, branded_info, portions }} />

      <table>
        {nutrients
          .filter((n) => n.amount && n.amount > 0.1)
          .map((n, x) => (
            <tr key={x}>
              <td>{n.amount && scaledRound(n.amount)} </td>
              <td>{n.nutrient && n.nutrient.unitName} </td>
              <td>{n.nutrient && n.nutrient.name} </td>
            </tr>
          ))}
      </table>
    </div>
  );
};
