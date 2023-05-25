import { cn } from "./lib";

const Loading: React.FC<{
  children?: React.ReactNode;
  loading: boolean;
}> = ({ children, loading }) => {
  return (
    <div
      className={cn(
        loading ? "h-4 bg-gray-400 rounded animate-pulse my-2 w-full" : ""
      )}
    >
      {children}
    </div>
  );
};

export default Loading;
