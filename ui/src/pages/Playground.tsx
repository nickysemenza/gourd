import React, { useContext, useEffect, useState } from "react";
import { PaginatedRecipes, RecipesApi } from "../api/openapi-fetch";
import RecipeDiff from "../components/RecipeDiff";
import { getAPIURL, getOpenapiFetchConfig } from "../config";
import { WasmContext } from "../wasm";

const Playground: React.FC = () => {
  const url = getAPIURL();
  const [r2, setR2] = useState<PaginatedRecipes>();

  const instance = useContext(WasmContext);

  useEffect(() => {
    if (!instance) return;
    console.log("parse", instance.parse("2 cups (240g) flour, sifted"));
    console.log("parse3", instance.parse3("2 cups (240g) flour, sifted"));
    console.log("parse4", instance.parse4("2 cups (240g) flour, sifted"));
    // greet();
    // parse("2 cups flour");
  }, [instance]);

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
      <RecipeDiff details={r2.recipes[0].versions} />
    </div>
  );
};
export default Playground;
