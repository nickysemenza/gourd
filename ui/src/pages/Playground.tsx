import React, { useContext, useEffect, useState } from "react";
import { Helmet } from "react-helmet";
import { PaginatedRecipeWrappers, RecipesApi } from "../api/openapi-fetch";
import {
  UnitConversionRequest,
  UnitConversionRequestTargetEnum,
} from "../api/openapi-fetch/models/UnitConversionRequest";
import { getAPIURL, getOpenapiFetchConfig } from "../config";
import { WasmContext } from "../wasm";

const ParseTest: React.FC = () => {
  const w = useContext(WasmContext);

  useEffect(() => {
    if (!w) return;
    console.log({ parse: w.parse("2 cups (240g) flour, sifted") });
    console.log({ parse2: w.parse2("2 cups (240g) flour, sifted") });
    console.log({ parse3: w.parse3("2 cups (240g) flour, sifted") });
    console.log({ parse4: w.parse4("2 cups (240g) flour, sifted") });
    // ingredients.forEach((i) => {
    let foo: UnitConversionRequest = {
      target: UnitConversionRequestTargetEnum.MONEY,
      unit_mappings: [
        {
          a: {
            unit: "cents",
            value: 240,
          },
          b: {
            unit: "gram",
            value: 1,
          },
          source: "fdc",
        },
      ],
      input: [{ unit: "grams", value: 100 }],
    };
    let foo2: UnitConversionRequest = {
      target: UnitConversionRequestTargetEnum.WEIGHT,
      unit_mappings: [
        {
          a: {
            unit: "cup",
            value: 1,
          },
          b: {
            unit: "gram",
            value: 125,
          },
          source: "fdc",
        },
      ],
      input: [{ unit: "cups", value: 1.5 }],
    };
    console.time("dolla");
    try {
      console.log("dolla", w.dolla(foo));
      console.log("dolla2", w.dolla(foo2));
    } catch (e) {
      console.error({ e });
    }
    console.timeEnd("dolla");
    // greet();
    // parse("2 cups flour");
  }, [w]);
  return null;
};

const Playground: React.FC = () => {
  const url = getAPIURL();
  const [r2, setR2] = useState<PaginatedRecipeWrappers>();

  useEffect(() => {
    const fetchData = async () => {
      const bar = new RecipesApi(getOpenapiFetchConfig());
      const result = await bar.listRecipes({});
      setR2(result);
    };
    fetchData();
  }, [url]);

  if (!r2 || !r2.recipes) return null;
  return (
    <div className="grid grid-cols-2 gap-4">
      {/* @ts-ignore */}
      <Helmet>
        <title>playground | gourd</title>
      </Helmet>
      <h1>playground</h1>
      <ParseTest />
    </div>
  );
};
export default Playground;
