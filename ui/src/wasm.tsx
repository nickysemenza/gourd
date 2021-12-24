import { createContext, useEffect, useState } from "react";

export type wasm = typeof import("gourd_rs");

export const WasmContext = createContext<wasm | undefined>(undefined);

export const WasmContextProvider: React.FC = ({ children }) => {
  // const [cursor, setCursor] = useState({ active: false });
  const [state, setState] = useState<wasm>();
  useEffect(() => {
    const fetchWasm = async () => {
      console.time("wasm-load");
      const wasm = await import("gourd_rs");
      setState(wasm);
      console.timeEnd("wasm-load");
    };
    fetchWasm();
  }, []);

  return <WasmContext.Provider value={state}>{children}</WasmContext.Provider>;
};
