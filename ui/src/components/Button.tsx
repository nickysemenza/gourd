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
  const baseStyles = `${compact ? "h-6" : "h-8"}  ${
    compact ? "px-3" : "px-5"
  } text-indigo-100 transition-colors duration-150
  bg-indigo-700 focus:shadow-outline hover:bg-indigo-800
  disabled:bg-indigo-500 inline-flex items-center`;
  return (
    <div className="inline-flex" role="group" aria-label="Button group">
      {buttons.map(({ text, IconLeft, IconRight, ...props }, x) => {
        const iconMargins = compact ? 1 : 3;
        const iconDim = compact ? 3 : 4;
        return (
          <button
            className={`${baseStyles} ${x === 0 && "rounded-l-lg"} ${
              x === buttons.length - 1 && "rounded-r-lg"
            }`}
            {...props}
          >
            {!!IconLeft && (
              <IconLeft
                className={`w-${iconDim} h-${iconDim} mr-${iconMargins}`}
              />
            )}
            {text && <span>{text}</span>}
            {!!IconRight && (
              <IconRight
                className={`w-${iconDim} h-${iconDim} ml-${iconMargins}`}
              />
            )}
          </button>
        );
      })}
    </div>
  );
};
