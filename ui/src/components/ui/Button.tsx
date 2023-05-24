import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { VariantProps, cva } from "class-variance-authority";
import { cn } from "./lib";

const buttonVariants = cva(
  "transition-colors duration-150 text-sm focus:shadow-outline inline-flex items-center",
  {
    variants: {
      variant: {
        primary:
          "text-violet-100 bg-violet-700 hover:bg-violet-800 disabled:bg-violet-400",
      },
      size: {
        default: "h-8 py-2 px-4",
        compact: "h-5 py-1 px-2",
      },
      round: {
        left: "rounded-l",
        right: "rounded-r",
        all: "rounded-md",
        none: "",
      },
    },
    defaultVariants: {
      size: "default",
      round: "all",
      variant: "primary",
    },
  }
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, round, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button";
    return (
      <Comp
        className={cn(buttonVariants({ round, size, className }))}
        ref={ref}
        {...props}
      />
    );
  }
);
Button.displayName = "Button";

export { Button, buttonVariants };
