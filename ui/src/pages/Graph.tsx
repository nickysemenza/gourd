import React, { useRef } from "react";
import Highcharts, { Point, PointClickEventObject } from "highcharts";
import HighchartsReact from "highcharts-react-official";
import HighchartsNetworkGraph from "highcharts/modules/networkgraph";

import { useRecipeDependencies } from "../api/openapi-hooks/api";
import { useNavigate } from "react-router-dom";

type CustomPoint = Highcharts.Point & {
  raw_kind: "ingredient" | "recipe";
  raw_id: string;
};

const Graph: React.FC = () => {
  HighchartsNetworkGraph(Highcharts);

  const { data } = useRecipeDependencies({});
  const history = useNavigate();
  const chartComponentRef = useRef<HighchartsReact.RefObject>(null);

  if (!data || !data.items) return null;

  const links = data.items.map((d) => {
    return {
      from: d.ingredient_name,
      to: d.recipe_name,
      value: d.ingredient_kind === "ingredient" ? 1 : 40,
    };
  });
  const nodes = [
    ...data.items.map((d) => ({
      id: d.ingredient_name,
      raw_id: d.ingredient_id,
      raw_kind: d.ingredient_kind,
      group: d.ingredient_kind === "ingredient" ? 1 : 4,
      color: "orange",
    })),
    ...data.items.map((d) => ({
      id: d.recipe_name,
      raw_id: d.recipe_id,
      raw_kind: "recipe",
      group: 30,
      color: "green",
    })),
  ].filter(
    (thing, index, self) => index === self.findIndex((t) => t.id === thing.id)
  );

  const options: Highcharts.Options = {
    chart: {
      type: "networkgraph",
    },
    title: {
      text: "ingredients",
    },
    plotOptions: {
      networkgraph: {
        keys: ["from", "to"],
        layoutAlgorithm: {
          enableSimulation: true,
          // friction: -0.9,
          integration: "verlet",
          // linkLength: 80,
        },
        marker: {
          radius: 5,
          lineWidth: 1,
        },
      },
    },
    series: [
      {
        link: {
          color: "rgba(100, 100, 100, 0.8)",
          dashStyle: "dash",
        },
        point: {
          events: {
            click: function (this: Point, _event: PointClickEventObject) {
              const node = this as unknown as CustomPoint;
              console.log(node);

              history(
                node.raw_kind === "recipe"
                  ? `/recipe/${node.raw_id}`
                  : `/ingredients/${node.raw_id}`
              );
            },
          },
        },

        type: "networkgraph",
        dataLabels: {
          enabled: true,
          linkFormat: "",
        },
        id: "lang-tree",
        data: links,
        nodes: nodes,
      },
    ],
  };

  console.log({ links, nodes });
  return (
    <div>
      <div className="text-gray-900 flex">Graph</div>
      <div className="border border-indigo-600 w-full h-full">
        <HighchartsReact
          highcharts={Highcharts}
          options={options}
          ref={chartComponentRef}
          height={900}
        />
      </div>
    </div>
  );
};
export default Graph;
