import React from "react";
import { Icon } from "react-feather";

interface ButtonProps {
  // icons from https://feathericons.com/
  IconLeft?: Icon;
  IconRight?: Icon;
  text?: string;
  onClick: () => void;
  disabled?: boolean;
}
interface ButtonGroupProps {
  buttons: ButtonProps[];
  compact?: boolean;
}
// https://tailwind-starter-kit.now.sh/docs/buttons#
export const ButtonGroup: React.FC<ButtonGroupProps> = ({
  buttons,
  compact = false,
}) => {
  const baseStyles = `${compact ? "h-5" : "h-8"}  ${
    compact ? "px-2" : "px-5"
  } text-indigo-100 transition-colors duration-150 text-sm
  bg-indigo-700 focus:shadow-outline hover:bg-indigo-800
  disabled:bg-indigo-500 inline-flex items-center`;
  return (
    <div className="inline-flex" role="group" aria-label="Button group">
      {buttons.map(({ text, IconLeft, IconRight, ...props }, x) => {
        const iconMargins = compact ? 1 : 3;
        const iconDim = compact ? 3 : 4;
        return (
          <button
            key={x}
            className={`${baseStyles} ${x === 0 && "rounded-l-lg"} ${
              x === buttons.length - 1 && "rounded-r-lg"
            }`}
            {...props}
          >
            {!!IconLeft && (
              <IconLeft
                className={`w-${iconDim} h-${iconDim} mr-${
                  text ? iconMargins : 0
                }`}
              />
            )}
            {text && <span>{text}</span>}
            {!!IconRight && (
              <IconRight
                className={`w-${iconDim} h-${iconDim} ml-${
                  text ? iconMargins : 0
                }`}
              />
            )}
          </button>
        );
      })}
    </div>
  );
};
export const Pill: React.FC = ({ children }) => (
  <span className="inline-flex items-center justify-center px-2 py-1 mr-1 text-xs font-bold leading-none bg-blue-200 text-blue-800 rounded-full">
    {children}
  </span>
);
export const Pill2: React.FC<{ color: "red" | "green" }> = ({
  children,
  color,
}) => (
  <span
    className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-${color}-100 text-${color}-800`}
  >
    {children}
  </span>
);

export const PillLabel: React.FC<{ x: number; kind: "letter" | "number" }> = ({
  x,
  kind,
}) => <Pill>{kind === "letter" ? String.fromCharCode(65 + x) : x}</Pill>;
