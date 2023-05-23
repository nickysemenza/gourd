import { Disclosure, Menu, Transition } from "@headlessui/react";
import React, { Fragment } from "react";
import { Menu as MenuIcon, X } from "react-feather";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { isLoggedIn, parseJWT } from "../../auth/auth";
import { Logo } from "../Logo";
import { navItems } from "./navItems";

const NavBar: React.FC = () => {
  const history = useNavigate();
  const location = useLocation();

  function classNames(...classes: string[]) {
    return classes.filter(Boolean).join(" ");
  }
  return (
    <Disclosure
      as="nav"
      className="bg-white/95 dark:bg-gray-800 mb-2 w-100 border-b border-slate-900/10 backdrop-blur"
    >
      {({ open }) => (
        <>
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 ">
            <div className="flex items-center justify-between h-16">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <div
                    onClick={() => history("/")}
                    className="duration-300 transform hover:rotate-45"
                  >
                    <Logo />
                  </div>
                </div>
                <div className="hidden md:block">
                  <div className="ml-10 flex items-baseline space-x-4">
                    {navItems.map((item) => (
                      <Link
                        key={item.path}
                        to={item.path}
                        // className={`text-gray-400 dark:hover:text-white px-3 py-2  text-sm font-medium ${
                        //   location.pathname.includes(item.path) &&
                        //   "dark:bg-gray-900 text-gray-800 dark:text-white"
                        // }`}
                        className={classNames(
                          location.pathname.includes(item.path)
                            ? "bg-gray-500 text-white"
                            : "text-gray-700 hover:bg-gray-300 hover:text-white",
                          "rounded-md px-3 py-2 text-sm font-medium"
                        )}
                      >
                        {item.title}
                      </Link>
                    ))}
                  </div>
                </div>
              </div>
              <div className="hidden md:block">
                <div className="ml-4 flex items-center md:ml-6">
                  {/* <button className="bg-gray-800 p-1 rounded-full text-gray-400 hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-white">
                    <span className="sr-only">View notifications</span>
                    <Bell className="h-6 w-6" aria-hidden="true" />
                  </button> */}

                  {/* Profile dropdown */}
                  {isLoggedIn() && (
                    <Menu as="div" className="ml-3 relative">
                      {({ open }) => (
                        <>
                          <div>
                            <Menu.Button className="max-w-xs bg-gray-800 rounded-full flex items-center text-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-white">
                              <span className="sr-only">Open user menu</span>
                              <img
                                className="h-8 w-8 rounded-full"
                                src={parseJWT()?.user_info.picture}
                                alt=""
                              />
                            </Menu.Button>
                          </div>
                          <Transition
                            show={open}
                            as={Fragment}
                            enter="transition ease-out duration-100"
                            enterFrom="transform opacity-0 scale-95"
                            enterTo="transform opacity-100 scale-100"
                            leave="transition ease-in duration-75"
                            leaveFrom="transform opacity-100 scale-100"
                            leaveTo="transform opacity-0 scale-95"
                          >
                            <Menu.Items
                              static
                              className="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-white ring-1 ring-black ring-opacity-5 focus:outline-none"
                            >
                              {["foo", "bar"].map((item) => (
                                <Menu.Item key={item}>
                                  {({ active }) => (
                                    <a
                                      href="#foo"
                                      className={classNames(
                                        active ? "bg-gray-100" : "",
                                        "block px-4 py-2 text-sm text-gray-700"
                                      )}
                                    >
                                      {item}
                                    </a>
                                  )}
                                </Menu.Item>
                              ))}
                            </Menu.Items>
                          </Transition>
                        </>
                      )}
                    </Menu>
                  )}
                </div>
              </div>
              <div className="-mr-2 flex md:hidden">
                {/* Mobile menu button */}
                <Disclosure.Button className="bg-gray-800 inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-white">
                  <span className="sr-only">Open main menu</span>
                  {open ? (
                    <X className="block h-6 w-6" aria-hidden="true" />
                  ) : (
                    <MenuIcon className="block h-6 w-6" aria-hidden="true" />
                  )}
                </Disclosure.Button>
              </div>
            </div>
          </div>
        </>
      )}
    </Disclosure>
  );
};

export default NavBar;
