import React, { useState } from "react";
import { gql } from "apollo-boost";

import styled from "styled-components";
import {
  color,
  SpaceProps,
  ColorProps,
  space,
  BordersProps,
  border,
} from "styled-system";
import { Box, Card, Heading, Text, Button } from "rebass";
import { Input } from "theme-ui";

const Box2 = styled.div<SpaceProps & ColorProps & BordersProps>`
  ${color}
  ${space}
${border}
`;

const Table = styled.table<SpaceProps & ColorProps & BordersProps>`
  ${color}
  ${space}
${border}
`;

const Td = styled.td<SpaceProps & ColorProps & BordersProps>`
  ${color}
  ${space}
${border}
`;

const Tr = styled.tr<SpaceProps & ColorProps & BordersProps>`
  ${color}
  ${space}
${border}
`;

const Test: React.FC = () => {
  return (
    <div>
      <Box2 color="red" mx={3} bg="navy" py={4}>
        asd
      </Box2>
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
