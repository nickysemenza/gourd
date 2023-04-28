import React, { useContext } from "react";
import {
  Highlight,
  HitsPerPage,
  InstantSearch,
  RefinementList,
  SearchBox,
  useHits,
} from "react-instantsearch-hooks-web";

import { instantMeiliSearch } from "@meilisearch/instant-meilisearch";
import { BrandedFoodItem, RecipeDetail } from "../api/openapi-hooks/api";
import { WasmContext } from "../wasmContext";
import { FoodRow } from "../components/FoodSearch";
import { FoodInfo } from "../api/openapi-fetch";
import { RecipeGridCell } from "../components/RecipeGrid";

const searchClient = instantMeiliSearch("http://localhost:7700", "FOO");

function BrandedHit(props: { hit: any }) {
  let hit = props.hit as BrandedFoodItem;
  const w = useContext(WasmContext);
  let f: FoodInfo | undefined = w?.bfi_to_info(hit);
  return (
    <div className="border-1">
      {f && (
        <FoodRow
          info={f}
          loading={false}
          wide
          descriptionComponent={
            <Highlight attribute="description" hit={props.hit} />
          }
          brandOwnerComponent={
            <Highlight attribute="brandOwner" hit={props.hit} />
          }
        />
      )}
    </div>
  );
}

function RecipeDetailHit(props: { hit: any }) {
  let hit = props.hit as RecipeDetail;
  return (
    <RecipeGridCell
      detail={hit}
      nameComponent={<Highlight attribute="name" hit={props.hit} />}
    />
  );
}

const BrandedHits: React.FC = () => {
  const { hits } = useHits();
  return (
    <div>
      {hits.map((hit) => (
        <BrandedHit key={(hit as unknown as BrandedFoodItem).fdcId} hit={hit} />
      ))}
    </div>
  );
};

const RecipeDetailsHits: React.FC = () => {
  const { hits } = useHits();
  return (
    <div className="grid gap-5 row-gap-5 mb-8 lg:grid-cols-6 md:grid-cols-4 sm:grid-cols-2">
      {hits.map((hit) => (
        <RecipeDetailHit key={(hit as unknown as RecipeDetail).id} hit={hit} />
      ))}
    </div>
  );
};
const searchClassNames = {
  root: "p-3 shadow-sm",
  form: "relative",
  input:
    "block w-full pl-9 pr-3 py-2 bg-white border border-slate-300 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-sky-500 rounded-md focus:ring-1",
  // submitIcon: "absolute top-0 left-0 bottom-0 w-6",
};
const refinementClassNames = {
  label: "flex space-x-1",
  count: "text-sm font-bold",
};
const Search: React.FC = () => (
  <div>
    <BrandedSearch />
    <hr />
    <RecipeDetailsSearch />
  </div>
);
const BrandedSearch: React.FC = () => (
  <InstantSearch indexName="BrandedFoods" searchClient={searchClient}>
    <div>
      <SearchBox classNames={searchClassNames} />
      <div className="flex">
        <div>
          <RefinementList
            classNames={refinementClassNames}
            attribute="brandOwner"
            showMore={true}
          />
          <RefinementList
            classNames={refinementClassNames}
            attribute="brandedFoodCategory"
          />
          <RefinementList
            classNames={refinementClassNames}
            attribute="servingSizeUnit"
          />
          <HitsPerPage
            items={[
              { label: "8 hits per page", value: 8, default: true },
              { label: "16 hits per page", value: 16 },
            ]}
          />
        </div>
        <BrandedHits />
      </div>
    </div>
  </InstantSearch>
);

const RecipeDetailsSearch: React.FC = () => (
  <InstantSearch indexName="RecipeDetails" searchClient={searchClient}>
    <div>
      <SearchBox classNames={searchClassNames} />
      <div className="flex">
        <div>
          {/* <RefinementList attribute="brandOwner" showMore={true} /> */}
          <RefinementList classNames={refinementClassNames} attribute="tags" />
          {/* <RefinementList attribute="servingSizeUnit" /> */}
        </div>
        <RecipeDetailsHits />
      </div>
    </div>
  </InstantSearch>
);

export default Search;
