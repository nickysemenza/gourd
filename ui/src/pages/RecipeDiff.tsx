import React, { useMemo } from "react";
import { Helmet } from "react-helmet";
import { useLocation } from "react-router-dom";
import RecipeDiffView from "../components/RecipeDiffView";
import queryString from "query-string";

const Playground: React.FC = () => {
  const loc = useLocation();
  const ids = useMemo(() => {
    const url = queryString.parse(loc.search).recipes;
    return url ? (Array.isArray(url) ? url : [url]) : [];
  }, [loc]);

  return (
    <div>
      <Helmet>
        <title>diff | gourd</title>
      </Helmet>
      <RecipeDiffView ids={ids} />
    </div>
  );
};
export default Playground;
