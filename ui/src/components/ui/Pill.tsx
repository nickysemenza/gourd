import { Badge } from "./Badge";

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
}) => <Badge>{kind === "letter" ? String.fromCharCode(65 + x) : x}</Badge>;
