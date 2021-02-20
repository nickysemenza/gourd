import { useEffect, useState } from "react";

export type wasm = typeof import("gourd_rs");
export const useWasm = () => {
  const [state, setState] = useState<wasm>();
  useEffect(() => {
    const fetchWasm = async () => {
      const startTime = new Date().getMilliseconds();
      const wasm = await import("gourd_rs");
      setState(wasm);
      console.log(
        `loaded wasm in ${new Date().getMilliseconds() - startTime}ms`
      );
    };
    fetchWasm();
  }, []);
  return state;
};
