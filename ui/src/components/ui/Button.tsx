import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { VariantProps, cva } from "class-variance-authority";
import { cn } from "./lib";

const buttonVariants = cva(
  "text-violet-100 transition-colors duration-150 text-sm bg-violet-700 focus:shadow-outline hover:bg-violet-800 disabled:bg-violet-400 inline-flex items-center",
  {
    variants: {
      size: {
        default: "h-8 py-5",
        compact: "h-5 px-2",
      },
      round: {
        left: "rounded-l",
        right: "rounded-r",
        all: "rounded",
      },
    },
    defaultVariants: {
      size: "default",
      round: "all",
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
