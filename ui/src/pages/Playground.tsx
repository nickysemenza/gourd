import React from "react";
import { Configuration, RecipesApi } from "../api/openapi-fetch";
import { useListRecipes } from "../api/openapi-hooks/api";
import Debug from "../components/Debug";

const Playground: React.FC = () => {
  const foo = useListRecipes({ base: "http://localhost:4242/api" });

  const c = new Configuration({ basePath: "http://localhost:4242/api" });
  const bar = new RecipesApi(c);
  const r2 = bar.listRecipes({});
  return (
    <div className="grid grid-cols-5 gap-4">
      <Debug data={{ foo, r: foo.data?.recipes, r2 }} />
    </div>
  );
};
export default Playground;
