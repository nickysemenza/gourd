const Pill3: React.FC<{
  children?: React.ReactNode;
}> = ({ children }) => (
  <span className="inline-flex items-center justify-center px-2 py-1 mr-1 text-xs font-bold leading-none bg-violet-200 text-violet-800 rounded-lg">
    {children}
  </span>
);
export const Pill: React.FC<{
  color: "red" | "green";
  children?: React.ReactNode;
}> = ({ children, color }) => (
  <span
    className={`px-1 h-5 inline-flex text-xs leading-5 font-semibold rounded-full bg-${color}-100 text-${color}-800`}
  >
    {children}
  </span>
);

export const PillLabel: React.FC<{ x: number; kind: "letter" | "number" }> = ({
  x,
  kind,
}) => <Pill3>{kind === "letter" ? String.fromCharCode(65 + x) : x}</Pill3>;
