import { createContext, useEffect, useState } from "react";

export type wasm = typeof import("gourd_rs");
// const useWasm = () => {
//   const [state, setState] = useState<wasm>();
//   useEffect(() => {
//     const fetchWasm = async () => {
//       const startTime = new Date().getMilliseconds();
//       const wasm = await import("gourd_rs");
//       setState(wasm);
//       console.log(
//         `loaded wasm in ${new Date().getMilliseconds() - startTime}ms`
//       );
//     };
//     fetchWasm();
//   }, []);
//   return state;
// };

export const WasmContext = createContext<wasm | undefined>(undefined);

export const WasmContextProvider: React.FC = ({ children }) => {
  // const [cursor, setCursor] = useState({ active: false });
  const [state, setState] = useState<wasm>();
  useEffect(() => {
    const fetchWasm = async () => {
      const startTime = new Date().getMilliseconds();
      const wasm = await import("gourd_rs");
      setState(wasm);
      const endTime = new Date().getMilliseconds();
      const delta = endTime - startTime;
      console.info(`loaded wasm in ${delta}ms`);
    };
    fetchWasm();
  }, []);

  return <WasmContext.Provider value={state}>{children}</WasmContext.Provider>;
};
