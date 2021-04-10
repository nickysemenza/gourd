import React from "react";
import ForceGraph2D, { GraphData, NodeObject } from "react-force-graph-2d";
import { useRecipeDependencies } from "../api/openapi-hooks/api";
import { useHistory } from "react-router-dom";

type CustomNode = NodeObject & {
  raw_kind: "ingredient" | "recipe";
  raw_id: string;
};
const Graph: React.FC = () => {
  const { data } = useRecipeDependencies({});
  let history = useHistory();

  if (!data || !data.items) return null;
  const data3: GraphData = {
    links: data.items.map((d) => {
      return {
        source: d.ingredient_name,
        target: d.recipe_name,
        value: d.ingredient_kind === "ingredient" ? 1 : 4,
      };
    }),
    nodes: [
      ...data.items.map((d) => ({
        id: d.ingredient_name,
        raw_id: d.ingredient_id,
        raw_kind: "ingredient",
        group: 1,
      })),
      ...data.items.map((d) => ({
        id: d.recipe_name,
        raw_id: d.recipe_id,
        raw_kind: "recipe",
        group: 3,
      })),
    ].filter(
      (thing, index, self) => index === self.findIndex((t) => t.id === thing.id)
    ),
  };
  console.log(data3);
  return (
    <div>
      <div className="text-gray-900 flex">Graph</div>
      <div
        className="border border-indigo-600"
        style={{ width: "1000px", height: "800px" }}
      >
        <ForceGraph2D
          width={1000}
          height={800}
          // dagMode="radialout"
          onNodeClick={(n, _e) => {
            const node = n as CustomNode;
            history.push(
              node.raw_kind === "recipe"
                ? `/recipe/${node.raw_id}`
                : `/ingredients/${node.raw_id}`
            );
            console.log(node);
          }}
          graphData={data3}
          nodeAutoColorBy="group"
          // @ts-ignore
          nodeCanvasObject={(node, ctx, globalScale) => {
            const label = node.id || "";
            const fontSize = 12 / globalScale;
            ctx.font = `${fontSize}px Sans-Serif`;
            // @ts-ignore
            const textWidth = ctx.measureText(label).width;
            const bckgDimensions = [textWidth, fontSize].map(
              (n) => n + fontSize * 0.2
            ); // some padding

            ctx.fillStyle = "rgba(255, 255, 255, 0.8)";
            // @ts-ignore
            ctx.fillRect(
              // @ts-ignore
              node.x - bckgDimensions[0] / 2,
              // @ts-ignore
              node.y - bckgDimensions[1] / 2,
              // @ts-ignore
              ...bckgDimensions
              // @ts-ignore
            );

            ctx.textAlign = "center";
            ctx.textBaseline = "middle";
            // @ts-ignore
            ctx.fillStyle = node.color;
            // @ts-ignore
            ctx.fillText(label, node.x, node.y);

            // @ts-ignore
            node.__bckgDimensions = bckgDimensions; // to re-use in nodePointerAreaPaint
          }}
          nodePointerAreaPaint={(node, color, ctx) => {
            ctx.fillStyle = color;
            // @ts-ignore
            const bckgDimensions = node.__bckgDimensions;
            bckgDimensions &&
              ctx.fillRect(
                // @ts-ignore
                node.x - bckgDimensions[0] / 2,
                // @ts-ignore
                node.y - bckgDimensions[1] / 2,
                // @ts-ignore
                ...bckgDimensions
              );
          }}
        />
      </div>
    </div>
  );
};
export default Graph;
