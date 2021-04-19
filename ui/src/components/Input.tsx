import React from "react";
import { formatText } from "../util";

// <input> can't do onBlur?
export type TallOrBlur =
  | {
      tall: true;
      blur?: false;
    }
  | {
      blur?: boolean;
      tall?: false;
    };

export type TableInputProps = {
  edit: boolean;
  softEdit?: boolean;
  value: string | number;
  width?: number | "full";
  highlight?: boolean;
  onChange: (event: string) => void;
} & TallOrBlur;

export const TableInput: React.FC<TableInputProps> = ({
  edit,
  softEdit = false,
  width = 10,
  tall = false,
  blur = false,
  highlight = false,
  value,
  onChange,
  ...props
}) => {
  const controlledVal = value.toString();
  const [internalVal, setVal] = React.useState(controlledVal);
  React.useEffect(() => {
    setVal(controlledVal);
  }, [controlledVal]);

  const oC = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setVal(e.target.value);
    if (!blur) {
      onChange(e.target.value);
    }
  };
  const oB = (e: React.FocusEvent<HTMLInputElement>) => {
    if (!blur) {
      return;
    }
    if (internalVal !== controlledVal) {
      onChange(internalVal);
    }
  };

  const className = `border-2 border-dashedp-0 h-${tall ? 18 : 6} w-${width} ${
    highlight ? "border-blue-400" : "border-gray-200"
  } disabled:border-red-100 hover:border-black ${
    softEdit && !edit && "bg-transparent"
  } focus:bg-gray-200`;

  return edit || softEdit ? (
    tall ? (
      <textarea
        {...props}
        value={internalVal}
        onChange={oC}
        className={className}
        rows={3}
      />
    ) : (
      <input
        {...props}
        value={internalVal}
        onChange={oC}
        onBlur={oB}
        className={className}
        disabled={!edit && controlledVal === "0"}
      />
    )
  ) : (
    <p className="flex flex-wrap">{formatText(internalVal)}</p>
  );
};
