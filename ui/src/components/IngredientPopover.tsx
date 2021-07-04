import { Popover, Transition } from "@headlessui/react";
import React, { Fragment } from "react";
import { AlertCircle } from "react-feather";
import { IngredientDetail } from "../api/openapi-hooks/api";
import Debug from "./Debug";
import { FoodRow } from "./FoodSearch";
import { UnitMappingList } from "./Misc";

// const solutions = [
//   {
//     name: "Insights",
//     description: "Measure actions your users take",
//     href: "##",
//     icon: IconOne,
//   },
//   {
//     name: "Automations",
//     description: "Create your own targeted content",
//     href: "##",
//     icon: IconTwo,
//   },
//   {
//     name: "Reports",
//     description: "Keep track of your growth",
//     href: "##",
//     icon: IconThree,
//   },
// ];

const IngredientPopover: React.FC<{ detail: IngredientDetail }> = ({
  detail,
}) => {
  const { ingredient, food, unit_mappings } = detail;
  return (
    <div className="">
      <Popover className="relative">
        {({ open }) => (
          <>
            <Popover.Button>
              <AlertCircle className="w-6 h-6 p-1 text-blue-700" />
            </Popover.Button>
            <Transition
              show={open}
              as={Fragment}
              enter="transition ease-out duration-200"
              enterFrom="opacity-0 translate-y-1"
              enterTo="opacity-100 translate-y-0"
              leave="transition ease-in duration-150"
              leaveFrom="opacity-100 translate-y-0"
              leaveTo="opacity-0 translate-y-1"
            >
              <Popover.Panel
                static
                className="absolute z-10 w-screen max-w-sm px-4 mt-3 transform -translate-x-1/2 left-1/2 sm:px-0 lg:max-w-3xl"
              >
                <div className="overflow-hidden rounded-lg shadow-lg ring-1 ring-black ring-opacity-5">
                  <div className="relative grid gap-8 bg-white p-7 lg:grid-cols-2">
                    <Debug data={ingredient} />
                    <UnitMappingList unit_mappings={unit_mappings} includeDot />
                    {food && <FoodRow food={food} />}
                    {/* <div>
                      {solutions.map((item) => (
                        <a
                          key={item.name}
                          href={item.href}
                          className="flex items-center p-2 -m-3 transition duration-150 ease-in-out rounded-lg hover:bg-gray-50 focus:outline-none focus-visible:ring focus-visible:ring-blue-500 focus-visible:ring-opacity-50"
                        >
                          <div className="flex items-center justify-center flex-shrink-0 w-10 h-10 text-white sm:h-12 sm:w-12">
                            <item.icon aria-hidden="true" />
                          </div>
                          <div className="ml-4">
                            <p className="text-sm font-medium text-gray-900">
                              {item.name}
                            </p>
                            <p className="text-sm text-gray-500">
                              {item.description}
                            </p>
                          </div>
                        </a>
                      ))}
                    </div> */}
                  </div>
                  <div className="p-4 bg-gray-50">
                    <a
                      href="##"
                      className="flow-root px-2 py-2 transition duration-150 ease-in-out rounded-md hover:bg-gray-100 focus:outline-none focus-visible:ring focus-visible:ring-blue-500 focus-visible:ring-opacity-50"
                    >
                      <span className="flex items-center">
                        <span className="text-sm font-medium text-gray-900">
                          Documentation
                        </span>
                      </span>
                      <span className="block text-sm text-gray-500">
                        Start integrating products and tools
                      </span>
                    </a>
                  </div>
                </div>
              </Popover.Panel>
            </Transition>
          </>
        )}
      </Popover>
    </div>
  );
};
export default IngredientPopover;
