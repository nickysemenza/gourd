export const Code: React.FC<{
  children?: React.ReactNode;
}> = ({ children }) => (
  <code className="rounded-sm px-2 py-0.5 bg-gray-100 text-red-500 h-6 flex text-sm font-bold w-fit">
    {children}
  </code>
);
