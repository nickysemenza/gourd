import React from "react";

const Debug: React.FC<{ data: any }> = ({ data }) => (
  <div className="rounded-md bg-gray-800 text-purple-300 text-xs">
    <pre className="scrollbar-none m-0 p-0">
      <code className="inline-block w-auto p-4 scrolling-touch">
        {JSON.stringify(data, null, 2)}
      </code>
    </pre>
  </div>
);

export default Debug;
