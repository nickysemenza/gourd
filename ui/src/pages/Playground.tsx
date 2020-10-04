import React, { useEffect, useState } from "react";
import { useGet } from "restful-react";
import { PaginatedIngredients, DefaultApi } from "../api/openapi-fetch";
import { useListIngredients } from "../api/openapi-hooks/api";
import Debug from "../components/Debug";
import { getAPIURL, getOpenapiFetchConfig } from "../config";

const Playground: React.FC = () => {
  const url = getAPIURL();
  const foo = useListIngredients({});
  const bar = useGet({ path: url + "/ingredients?limit=5&offset=10" });
  const [r2, setR2] = useState<PaginatedIngredients>();

  useEffect(() => {
    const fetchData = async () => {
      const bar = new DefaultApi(getOpenapiFetchConfig());
      const result = await bar.listIngredients({});
      setR2(result);
    };

    fetchData();
  }, [url]);

  return (
    <div className="grid grid-cols-2 gap-4">
      <Debug data={{ foo, r: foo.data, bar }} />
      <Debug data={{ r2 }} />
    </div>
  );
};
export default Playground;
