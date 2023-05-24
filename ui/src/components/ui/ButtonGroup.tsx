import React, { useState } from "react";
import { Icon, MinusCircle, PlusCircle } from "react-feather";
import { Button } from "./Button";

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
  return (
    <div className="inline-flex mx-1" role="group" aria-label="Button group">
      {buttons.map(({ text, IconLeft, IconRight, submit, ...props }, x) => {
        const iconMargins = compact ? 1 : 3.5;
        const iconDim = compact ? 12 : 18;
        return (
          <Button
            type={submit ? "submit" : "button"}
            key={x}
            round={
              buttons.length === 1
                ? "all"
                : x === 0
                ? "left"
                : x === buttons.length - 1
                ? "right"
                : "none"
            }
            size={compact ? "compact" : "default"}
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
          </Button>
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
