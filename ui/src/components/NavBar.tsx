import React, { useState } from "react";
import { Box, Flex, Text } from "rebass";
import { useHistory } from "react-router-dom";

const NavBar: React.FC = () => {
  const history = useHistory();
  return (
    <Flex px={2} color="white" bg="black" alignItems="center">
      <Text p={2} fontWeight="bold" onClick={() => history.push("/")}>
        Recipes
      </Text>
      <Box mx="auto" />
      <Box
        sx={{
          display: "inline-block",
          fontWeight: "bold",
          px: 2,
          py: 1,
          color: "inherit",
        }}
      >
        <a onClick={() => history.push("/recipes")}>Recipes</a>
      </Box>
      <Box
        sx={{
          display: "inline-block",
          fontWeight: "bold",
          px: 2,
          py: 1,
          color: "inherit",
        }}
      >
        <a onClick={() => history.push("/ingredients")}>Ingredients</a>
      </Box>
    </Flex>
  );
};

export default NavBar;
