import { useGetFoodQuery } from "../generated/graphql";
import React from "react";
import { Box } from "rebass";
import Debug from "../components/Debug";

const Food: React.FC = () => {
  // const [food, setFood]
  const { data, loading, error } = useGetFoodQuery({
    variables: {
      fdc_id: 747448,
    },
  });
  return (
    <Box>
      <Debug data={{ data, loading, error }} />
    </Box>
  );
};
export default Food;
