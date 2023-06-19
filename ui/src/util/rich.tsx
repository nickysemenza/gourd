import { wasm, RichItem } from "./wasmContext";

export const formatRichText = (w: wasm, text: RichItem[]) => {
  return text.map((t, x) => {
    if (t.kind === "Text") {
      return t.value;
    } else if (t.kind === "Ing") {
      return (
        <div
          className="inline text-orange-800 m-0 underline decoration-grey decoration-solid"
          key={x + "a"}
        >
          {t.value}
        </div>
      );
    } else if (t.kind === "Measure") {
      const val = t.value.pop();
      if (!val) {
        return null;
      }
      if (val.unit === "whole") {
        val.unit = "";
      }
      if (val.value === null) {
        val.value = 0;
      }
      return (
        <div
          className="inline text-green-800 m-0 underline decoration-grey decoration-solid"
          key={x}
        >
          {/* // sum_time_amounts converts the units to the best scaled */}
          {w.format_amount(val)}
        </div>
      );
    } else {
      return null;
    }
  });
};
