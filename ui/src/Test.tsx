import React from "react";

import { Box, Card, Heading, Text, Button } from "rebass";
import { useColorMode } from "theme-ui";

const Test: React.FC = () => {
  const [colorMode, setColorMode] = useColorMode();
  return (
    <div>
      <button
        onClick={(e) => {
          setColorMode(colorMode === "default" ? "dark" : "default");
        }}
      >
        Toggle {colorMode === "default" ? "Dark" : "Light"}
      </button>
      <Box width={256}>
        <Card
          sx={{
            p: 1,
            borderRadius: 2,
            boxShadow: "0 0 16px rgba(0, 0, 0, .25)",
          }}
        >
          <Box px={2}>
            <Heading as="h3">Card Demo</Heading>
            <Text fontSize={0}>You can edit this code</Text>
          </Box>
        </Card>
      </Box>

      <Button variant="primary" mr={2}>
        Beep
      </Button>
      <Button variant="secondary">Boop</Button>
    </div>
  );
};

export default Test;
