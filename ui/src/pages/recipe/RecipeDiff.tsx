import React, { useMemo } from "react";
import { Helmet } from "react-helmet";
import { useLocation } from "react-router-dom";
import RecipeDiffView from "../../components/recipe/RecipeDiffView";
import queryString from "query-string";
import { EntitySummary } from "../../api/react-query/gourdApiSchemas";

const Playground: React.FC = () => {
  const loc = useLocation();
  const ids = useMemo(() => {
    const url = queryString.parse(loc.search).recipes;
    const a = url ? (Array.isArray(url) ? url : [url]) : [];
    const u = a.filter((x) => x !== null) as string[];
    return u;
  }, [loc]);

  const input: EntitySummary[] = ids.map((id) => {
    return { id, multiplier: 1, name: "", kind: "recipe" };
  });
  return (
    <div>
      <Helmet>
        <title>diff | gourd</title>
      </Helmet>
      <RecipeDiffView entitiesToDiff={input} />
    </div>
  );
};
export default Playground;
