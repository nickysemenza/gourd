import React from "react";
import { Box } from "rebass";

const Debug: React.FC<{ data: any }> = ({ data }) => (
  <Box
    sx={{
      borderWidth: "1",
      borderStyle: "solid",
      borderColor: "highlight",
    }}
  >
    <pre>{JSON.stringify(data, null, 2)}</pre>
  </Box>
);

export default Debug;
