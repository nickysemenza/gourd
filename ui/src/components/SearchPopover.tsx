import { Transition } from "@headlessui/react";
import * as DialogPrimitive from "@radix-ui/react-dialog";
import { Cross1Icon } from "@radix-ui/react-icons";
import { clsx } from "clsx";
import React, { Fragment, useState } from "react";
import { useHotkeys } from "react-hotkeys-hook";
import { useSearch } from "../api/react-query/gourdApiComponents";
import { RecipeLink } from "./misc/Misc";
import { useNavigate } from "react-router-dom";

const Dialog = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [name, setName] = useState("");
  const navigate = useNavigate();

  useHotkeys("meta+k", () => setIsOpen(true));
  const { data } = useSearch(
    { queryParams: { name } },
    { enabled: isOpen && name !== "" }
  );

  const ingredients = data?.ingredients?.map((ingredient, x) => (
    <a
      key={x}
      className="w-full text-left px-3.5 py-2.5 rounded-md hover:bg-gray-100 dark:hover:bg-gray-800 focus:ring-1 focus:ring-gray-300 focus:outline-none flex items-center space-x-2.5 justify-between bg-white/50 dark:bg-gray-800 cursor-pointer"
      onClick={() => {
        navigate(`/ingredients/${ingredient.ingredient.id}`);
        setIsOpen(false);
      }}
    >
      <div className="flex w-3/4 items-center space-x-2.5">
        {ingredient.ingredient.name}
      </div>
      <div className="">
        {ingredient.recipes
          .filter((x) => x.meta.is_latest_version)
          .map((x) => (
            <RecipeLink recipe={x} key={x.id} />
          ))}
      </div>
    </a>
  ));

  const recipes = data?.recipes?.map((recipe, x) => (
    <div
      key={x}
      className="w-full text-left px-3.5 py-2.5 rounded-md hover:bg-gray-100 dark:hover:bg-gray-800 focus:ring-1 focus:ring-gray-300 focus:outline-none flex items-center space-x-2.5 justify-between bg-white/50 dark:bg-gray-800 cursor-pointer"
    >
      <RecipeLink recipe={recipe.detail} />
    </div>
  ));

  return (
    <DialogPrimitive.Root open={isOpen} onOpenChange={setIsOpen}>
      {/* <DialogPrimitive.Trigger asChild>
        <Button>Click</Button>
      </DialogPrimitive.Trigger> */}
      <DialogPrimitive.Portal forceMount>
        <Transition.Root show={isOpen}>
          <Transition.Child
            as={Fragment}
            enter="ease-out duration-300"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="ease-in duration-200"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <DialogPrimitive.Overlay
              forceMount
              className="fixed inset-0 z-20 bg-black/50"
            />
          </Transition.Child>
          <Transition.Child
            as={Fragment}
            enter="ease-out duration-300"
            enterFrom="opacity-0 scale-95"
            enterTo="opacity-100 scale-100"
            leave="ease-in duration-200"
            leaveFrom="opacity-100 scale-100"
            leaveTo="opacity-0 scale-95"
          >
            <DialogPrimitive.Content
              forceMount
              className={clsx(
                "fixed z-50",
                // "w-[95vw] max-w-md rounded-lg md:w-full",
                "w-full max-w-3xl",
                "top-[50%] left-[50%] -translate-x-[50%] -translate-y-[50%]",
                "bg-white dark:bg-gray-800",
                "focus:outline-none focus-visible:ring focus-visible:ring-purple-500 focus-visible:ring-opacity-75",
                "divide-y dark:divide-gray-800"
              )}
            >
              <div className="flex items-center space-x-1.5 pl-3">
                <input
                  type="text"
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Search"
                  className={clsx(
                    "py-4 px-0 border-none w-full focus:outline-none focus:border-none focus:ring-0 bg-transparent placeholder-gray-500 dark:text-white"
                  )}
                />
              </div>
              <div className="flex-1 overflow-y-auto focus:outline-none p-2 space-y-4">
                <h4 className="px-3.5 text-gray-500 text-sm font-medium">
                  Ingredients
                </h4>
                <ul tabIndex={-1}>{ingredients}</ul>
                <h4 className="px-3.5 text-gray-500 text-sm font-medium">
                  Recipes
                </h4>
                <ul>{recipes}</ul>
              </div>
              <DialogPrimitive.Close
                className={clsx(
                  "absolute top-3.5 right-3.5 inline-flex items-center justify-center rounded-full p-1",
                  "focus:outline-none focus-visible:ring focus-visible:ring-purple-500 focus-visible:ring-opacity-75"
                )}
              >
                <Cross1Icon className="h-4 w-4 text-gray-500 hover:text-gray-700 dark:text-gray-500 dark:hover:text-gray-400" />
              </DialogPrimitive.Close>
            </DialogPrimitive.Content>
          </Transition.Child>
        </Transition.Root>
      </DialogPrimitive.Portal>
    </DialogPrimitive.Root>
  );
};

export { Dialog };
