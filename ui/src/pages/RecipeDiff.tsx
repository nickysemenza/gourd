import React, { useMemo } from "react";
import { Helmet } from "react-helmet";
import { useLocation } from "react-router-dom";
import RecipeDiffView from "../components/RecipeDiffView";
import queryString from "query-string";

const Playground: React.FC = () => {
  const loc = useLocation();
  const ids = useMemo(() => {
    const url = queryString.parse(loc.search).recipes;
    const a = url ? (Array.isArray(url) ? url : [url]) : [];
    let u = a.filter((x) => x !== null) as string[];
    return u;
  }, [loc]);

  return (
    <div>
      {/* @ts-ignore */}
      <Helmet>
        <title>diff | gourd</title>
      </Helmet>
      <RecipeDiffView ids={ids} />
    </div>
  );
};
export default Playground;
