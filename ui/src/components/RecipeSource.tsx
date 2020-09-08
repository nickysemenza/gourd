import React from "react";
import { Source } from "../generated/graphql";

export interface Props {
  source: Pick<Source, "meta" | "name">;
}
const RecipeSource: React.FC<Props> = ({ source }) => (
  <div>
    From{" "}
    {source.meta.startsWith("http") ? (
      <a href={source.meta} target="_blank" rel="noopener noreferrer">
        {source.name}
      </a>
    ) : (
      <>
        {source.name} / {source.meta})
      </>
    )}
  </div>
);
export default RecipeSource;
