import { useGetFoodQuery } from "../generated/graphql";
import React from "react";
import { Box, Heading, Text } from "rebass";
import Debug from "../components/Debug";

const Food: React.FC = () => {
  // const [food, setFood]
  return (
    <Box
      sx={{
        display: "grid",
        gridGap: 3, // theme.space[3]
        gridTemplateColumns: "repeat(auto-fit, minmax(128px, 1fr))",
      }}
    >
      <FoodInfo fdc_id={171047} />
      <FoodInfo fdc_id={392941} />
      <FoodInfo fdc_id={747448} />
      <FoodInfo fdc_id={789097} />
      <FoodInfo fdc_id={335560} />
    </Box>
  );
};
export default Food;

const FoodInfo: React.FC<{ fdc_id: number }> = ({ fdc_id }) => {
  const { data, loading, error } = useGetFoodQuery({
    variables: {
      fdc_id,
    },
  });
  if (!data) return null;
  const { food } = data;
  // const { description = null, nutrients = null, category = null } = food;
  return (
    <Box>
      <Heading>{food?.description}</Heading>
      <Text>{food?.data_type}</Text>
      <Debug data={food?.category} />

      <table>
        {data?.food?.nutrients
          .filter((n) => n.amount > 0.1)
          .map((n) => (
            <tr>
              <td>{n.amount} </td>
              <td>{n.nutrient.unit_name} </td>
              <td>{n.nutrient.name} </td>
            </tr>
          ))}
      </table>
    </Box>
  );
};
