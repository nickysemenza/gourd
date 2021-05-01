import React from "react";
import Login from "./components/Login";

const Test: React.FC = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Gourd Recipe Database
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            <Login />
          </p>
        </div>
      </div>
    </div>
  );
};

export default Test;
