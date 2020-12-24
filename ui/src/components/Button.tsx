import React from "react";
interface ButtonProps {
  onClick: () => void;
  label: string;
  kind?: "primary" | "secondary";
  disabled?: boolean;
}
export const Button: React.FC<ButtonProps> = ({
  label,
  onClick,
  kind = "primary",
  ...props
}) => {
  return (
    <button
      className="px-1 py-1 text-sm text-purple-600 disabled:border-red-200 font-semibold rounded-full border border-purple-200 hover:text-white hover:bg-purple-600 hover:border-transparent focus:outline-none focus:ring-2 focus:ring-purple-600 focus:ring-offset-2"
      onClick={onClick}
      {...props}
    >
      {label}
    </button>
  );
};
