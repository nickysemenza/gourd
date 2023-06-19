import React from "react";
import { useGetFoodById } from "../api/openapi-hooks/api";
import Debug from "../components/ui/Debug";
import { UnitMappingList } from "../components/misc/Misc";
import { scaledRound } from "../util/util";
import { Code } from "../components/Code";

const Food: React.FC = () => {
  // const [food, setFood]
  return (
    <div className="grid grid-cols-3 gap-1">
      {/* <FoodInfo fdc_id={9999999} /> */}
      <FoodInfo fdc_id={171047} />
      <FoodInfo fdc_id={2028185} />
      <FoodInfo fdc_id={790018} />
      {/* <FoodInfo fdc_id={789097} /> */}
      {/* <FoodInfo fdc_id={335560} /> */}
    </div>
  );
};
export default Food;

const FoodInfo: React.FC<{ fdc_id: number }> = ({ fdc_id }) => {
  const { data } = useGetFoodById({ fdc_id });
  if (!data) return null;
  const { description, dataType } = data.wrapper;
  // const
  const { unit_mappings } = data;
  return (
    <div>
      <h2 className="font-bold text-l">{description}</h2>
      <div>
        <Code>{dataType}</Code>
      </div>
      <div>
        <Code>{fdc_id}</Code>
      </div>
      <UnitMappingList unit_mappings={unit_mappings} />
      <Debug data={{ data }} />

      <table>
        {(data.foodNutrients || [])
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
