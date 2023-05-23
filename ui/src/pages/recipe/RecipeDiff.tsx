import React, { useMemo } from "react";
import { Helmet } from "react-helmet";
import { useLocation } from "react-router-dom";
import RecipeDiffView from "../../components/recipe/RecipeDiffView";
import queryString from "query-string";
import { EntitySummary, IngredientKind } from "../../api/openapi-fetch";

const Playground: React.FC = () => {
  const loc = useLocation();
  const ids = useMemo(() => {
    const url = queryString.parse(loc.search).recipes;
    const a = url ? (Array.isArray(url) ? url : [url]) : [];
    let u = a.filter((x) => x !== null) as string[];
    return u;
  }, [loc]);

  let input: EntitySummary[] = ids.map((id) => {
    return { id, multiplier: 1, name: "", kind: IngredientKind.RECIPE };
  });
  return (
    <div>
      {/* @ts-ignore */}
      <Helmet>
        <title>diff | gourd</title>
      </Helmet>
      <RecipeDiffView entitiesToDiff={input} />
    </div>
  );
};
export default Playground;
