import * as React from "react";
import { VariantProps, cva } from "class-variance-authority";
import { cn } from "./lib";

const badgeVariants = cva(
  "inline-flex items-center border rounded-full px-2 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
  {
    variants: {
      variant: {
        default:
          "bg-slate-300 hover:bg-slate-300/80 border-transparent text-slate-700",
        violet:
          "bg-violet-300 hover:bg-violet-300/80 border-transparent text-violet-700",
        red: "bg-red-300 hover:bg-red-300/80 border-transparent text-red-700",
        outline: "text-foreground",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  );
}

export { Badge, badgeVariants };
