import React from "react";
import Login from "./components/Login";

const Test: React.FC = () => {
  return (
    <div>
      <Login />

      {/* https://mertjf.github.io/tailblocks/ */}
      <section className="text-gray-700 body-font">
        <div className="container px-5 py-24 mx-auto">
          <div className="flex flex-col text-center w-full mb-20">
            <h2 className="text-xs text-teal-500 tracking-widest font-medium title-font mb-1">
              Go Universal Recipe Database
            </h2>
            <h1 className="sm:text-3xl text-2xl font-medium title-font mb-4 text-gray-900">
              gourd
            </h1>
            <p className="lg:w-2/3 mx-auto leading-relaxed text-base">
              Nicky's Recipe Database
            </p>
          </div>
        </div>
      </section>
    </div>
  );
};

export default Test;
