import React, { useState } from "react";
import { Icon, MinusCircle, PlusCircle } from "react-feather";

type ButtonProps =
  | {
      // icons from https://feathericons.com/
      IconLeft?: Icon;
      IconRight?: Icon;
      text?: string;
      disabled?: boolean;
    } & ({ submit: true } | { submit?: false; onClick: () => void });
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
      {buttons.map(({ text, IconLeft, IconRight, submit, ...props }, x) => {
        const iconMargins = compact ? 1 : 3.5;
        const iconDim = compact ? 12 : 18;
        return (
          <button
            type={submit ? "submit" : "button"}
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
