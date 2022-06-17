import React, { useState } from "react";
import { Icon, MinusCircle, PlusCircle } from "react-feather";

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
  } text-violet-100 transition-colors duration-150 text-sm
  bg-violet-700 focus:shadow-outline hover:bg-violet-800
  disabled:bg-violet-400 inline-flex items-center`;
  return (
    <div className="inline-flex mx-1" role="group" aria-label="Button group">
      {buttons.map(({ text, IconLeft, IconRight, ...props }, x) => {
        const iconMargins = compact ? 1 : 3.5;
        const iconDim = compact ? 12 : 18;
        return (
          <button
            key={x}
            className={`${baseStyles} ${x === 0 && "rounded-l"} ${
              x === buttons.length - 1 && "rounded-r"
            }`}
            {...props}
          >
            {!!IconLeft && (
              <IconLeft
                width={iconDim}
                height={iconDim}
                className={`mr-${text ? iconMargins : 0}`}
              />
            )}
            {text && <span className="mx-1">{text}</span>}
            {!!IconRight && (
              <IconRight
                width={iconDim}
                height={iconDim}
                className={`ml-${text ? iconMargins : 0}`}
              />
            )}
          </button>
        );
      })}
    </div>
  );
};
export const Pill: React.FC<{
  children?: React.ReactNode;
}> = ({ children }) => (
  <span className="inline-flex items-center justify-center px-2 py-1 mr-1 text-xs font-bold leading-none bg-violet-200 text-violet-800 rounded-lg">
    {children}
  </span>
);
export const Pill2: React.FC<{
  color: "red" | "green";
  children?: React.ReactNode;
}> = ({ children, color }) => (
  <span
    className={`px-1 h-5 inline-flex text-xs leading-5 font-semibold rounded-full bg-${color}-100 text-${color}-800`}
  >
    {children}
  </span>
);

export const HideShowHOC: React.FC<{
  children?: React.ReactNode;
}> = ({ children }) => {
  const [show, setShow] = useState(false);
  return (
    <div>
      <HideShowButton show={show} setVal={setShow} /> {show && children}
    </div>
  );
};

export const makeHideShowButton = (
  show: boolean,
  setVal: (newVal: boolean) => void,
  text?: string
) => {
  return {
    onClick: () => {
      setVal(!show);
    },
    text: `${text ? text : show ? "hide" : "show"}`,
    IconLeft: show ? MinusCircle : PlusCircle,
  };
};
export const HideShowButton: React.FC<{
  show: boolean;
  setVal: (newVal: boolean) => void;
}> = ({ show, setVal }) => (
  <ButtonGroup buttons={[makeHideShowButton(show, setVal)]} />
);

export const PillLabel: React.FC<{ x: number; kind: "letter" | "number" }> = ({
  x,
  kind,
}) => <Pill>{kind === "letter" ? String.fromCharCode(65 + x) : x}</Pill>;
