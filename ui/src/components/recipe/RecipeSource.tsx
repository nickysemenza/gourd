import React from "react";

export interface Props {
  source: any;
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
