import React, { useEffect, useState } from "react";
import { PaginatedRecipes, RecipesApi } from "../api/openapi-fetch";
import RecipeDiff from "../components/RecipeDiff";
import { getAPIURL, getOpenapiFetchConfig } from "../config";
import { useWasm } from "../wasm";

const Playground: React.FC = () => {
  const url = getAPIURL();
  const [r2, setR2] = useState<PaginatedRecipes>();

  const instance = useWasm();

  useEffect(() => {
    if (!instance) return;
    instance.greet();
    console.log(instance.parse2("2 cups flour, sifted"));
    console.log(instance.parse3("2 cups flour, sifted"));
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
