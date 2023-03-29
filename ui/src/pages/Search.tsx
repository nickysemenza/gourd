import React, { useContext } from "react";
import {
  Highlight,
  InstantSearch,
  RefinementList,
  SearchBox,
  useHits,
} from "react-instantsearch-hooks-web";

import { instantMeiliSearch } from "@meilisearch/instant-meilisearch";
import { BrandedFoodItem } from "../api/openapi-hooks/api";
import { WasmContext } from "../wasmContext";
import { FoodRow } from "../components/FoodSearch";
import { FoodInfo } from "../api/openapi-fetch";

const searchClient = instantMeiliSearch("http://localhost:7700", "FOO");

function Hit(props: { hit: any }) {
  let hit = props.hit as BrandedFoodItem;
  const w = useContext(WasmContext);
  let f: FoodInfo | undefined = w?.bfi_to_info(hit);
  return (
    <div className="border-1">
      {f && <FoodRow info={f} loading={false} />}
      {/* <Debug data={props.hit.name} /> */}
      <Highlight attribute="brandOwner" hit={props.hit} />
      <Highlight attribute="description" hit={props.hit} />
    </div>
  );
}

const Hits: React.FC = () => {
  const { hits } = useHits();
  return (
    <div>
      {hits.map((hit) => (
        <Hit key={(hit as unknown as BrandedFoodItem).fdcId} hit={hit} />
      ))}
    </div>
  );
};

const Search: React.FC = () => (
  <InstantSearch indexName="BrandedFoods" searchClient={searchClient}>
    <div>
      <SearchBox
        classNames={{
          root: "p-3 shadow-sm",
          form: "relative",
          input:
            "block w-full pl-9 pr-3 py-2 bg-white border border-slate-300 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-sky-500 rounded-md focus:ring-1",
          // submitIcon: "absolute top-0 left-0 bottom-0 w-6",
        }}
      />
      <div className="flex">
        <div>
          <RefinementList attribute="brandOwner" showMore={true} />
          <RefinementList attribute="brandedFoodCategory" />
          <RefinementList attribute="servingSizeUnit" />
        </div>
        <Hits />
      </div>
    </div>
    {/* <Hits hitComponent={Hit as unknown as Foo} /> */}
  </InstantSearch>
);

export default Search;
