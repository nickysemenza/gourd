import { Popover, Transition } from "@headlessui/react";
import React, { Fragment } from "react";
import { AlertCircle } from "react-feather";
import { IngredientDetail } from "../api/openapi-hooks/api";
import Debug from "./Debug";
import { FoodRow } from "./FoodSearch";
import { UnitMappingList } from "./Misc";

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
              <AlertCircle
                className={`w-6 h-6 p-1 ${
                  unit_mappings.length > 0 ? "text-blue-700" : "text-red-700"
                } `}
              />
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
                  <div className="relative grid gap-4 bg-white p-1 lg:grid-cols-1">
                    <div className="p-4 bg-gray-50">
                      <span className="text-sm font-medium text-gray-900">
                        Metadata
                      </span>
                      <Debug data={ingredient} />
                    </div>
                    <div className="p-4 bg-gray-50">
                      <span className="text-sm font-medium text-gray-900">
                        Unit Mappings
                      </span>
                      <UnitMappingList
                        unit_mappings={unit_mappings}
                        includeDot
                      />
                    </div>
                    {food && (
                      <div className="p-4 bg-gray-50">
                        <span className="text-sm font-medium text-gray-900">
                          USDA food info
                        </span>
                        <FoodRow info={food} loading={false} />
                      </div>
                    )}
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
