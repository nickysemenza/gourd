import React from "react";
import { useHistory } from "react-router-dom";

const NavBar: React.FC = () => {
  const history = useHistory();
  const navItems: {
    title: string;
    path: string;
  }[] = [
    { title: "recipes", path: "recipes" },
    { title: "ingredients", path: "ingredients" },
    { title: "create", path: "create" },
    { title: "food (usda)", path: "food" },
    { title: "playground", path: "playground" },
    { title: "photos", path: "photos" },
    { title: "meals", path: "meals" },
  ];
  return (
    <nav className="flex items-center justify-between flex-wrap bg-teal-500 px-6 py-3 mb-2">
      <div className="flex items-center flex-shrink-0 text-white mr-1">
        <svg
          onClick={() => history.push("/")}
          width="158"
          height="40"
          viewBox="0 0 158 59"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M42.8111 20.6392C44.1847 17.1087 51.8288 17.1261 53.0829 20.6565C56.5273 19.1413 61.8939 24.4923 60.3709 27.8804C60.5787 27.9618 60.7741 28.0646 60.9572 28.1865C60.9852 28.202 61.0014 28.2125 61 28.2154C63.875 30.1955 63.6447 36.8982 60.3461 38.0746C61.8663 41.5441 56.4611 47.6078 53.0829 46C51.7094 49.5302 44.0659 49.5298 42.8111 46C39.3466 47.4437 34.0137 41.3934 35.61 38.0261C32.11 36.659 32.0294 29.1466 35.4829 27.8144C34.0048 24.3438 39.4422 19.0358 42.8111 20.6392Z"
            fill="#F4CA59"
          />
          <path
            d="M35.4829 27.8144C35.4787 27.8049 35.487 27.824 35.4829 27.8144ZM35.4829 27.8144C34.0048 24.3438 39.4422 19.0358 42.8111 20.6392C44.1847 17.1087 51.8288 17.1261 53.0829 20.6565C56.5273 19.1413 61.8939 24.4923 60.3709 27.8804M35.4829 27.8144C32.0294 29.1466 32.11 36.659 35.61 38.0261C34.0137 41.3934 39.3466 47.4437 42.8111 46C44.0659 49.5298 51.7094 49.5302 53.0829 46C56.4611 47.6078 61.8663 41.5441 60.3461 38.0746C63.8688 36.8182 63.8921 29.259 60.3709 27.8804M35.4829 27.8144C35.506 27.8055 35.4598 27.8233 35.4829 27.8144ZM60.3687 27.8795L60.3709 27.8804M60.3709 27.8804C60.3629 27.8984 61.0085 28.1976 61 28.2154"
            stroke="black"
            stroke-linejoin="round"
          />
          <circle cx="48" cy="33" r="9.5" fill="white" stroke="black" />
          <path
            d="M29.9 21.5V44.4C29.9 49.0333 28.75 52.45 26.45 54.65C24.15 56.8833 20.7167 58 16.15 58C13.65 58 11.2667 57.65 9 56.95C6.73333 56.2833 4.9 55.3167 3.5 54.05L5.8 50.35C7.03333 51.4167 8.55 52.25 10.35 52.85C12.1833 53.4833 14.0667 53.8 16 53.8C19.1 53.8 21.3833 53.0667 22.85 51.6C24.35 50.1333 25.1 47.9 25.1 44.9V42.8C23.9667 44.1667 22.5667 45.2 20.9 45.9C19.2667 46.5667 17.4833 46.9 15.55 46.9C13.0167 46.9 10.7167 46.3667 8.65 45.3C6.61667 44.2 5.01667 42.6833 3.85 40.75C2.68333 38.7833 2.1 36.55 2.1 34.05C2.1 31.55 2.68333 29.3333 3.85 27.4C5.01667 25.4333 6.61667 23.9167 8.65 22.85C10.7167 21.7833 13.0167 21.25 15.55 21.25C17.55 21.25 19.4 21.6167 21.1 22.35C22.8333 23.0833 24.25 24.1667 25.35 25.6V21.5H29.9ZM16.1 42.7C17.8333 42.7 19.3833 42.3333 20.75 41.6C22.15 40.8667 23.2333 39.85 24 38.55C24.8 37.2167 25.2 35.7167 25.2 34.05C25.2 31.4833 24.35 29.4167 22.65 27.85C20.95 26.25 18.7667 25.45 16.1 25.45C13.4 25.45 11.2 26.25 9.5 27.85C7.8 29.4167 6.95 31.4833 6.95 34.05C6.95 35.7167 7.33333 37.2167 8.1 38.55C8.9 39.85 9.98333 40.8667 11.35 41.6C12.75 42.3333 14.3333 42.7 16.1 42.7ZM90.627 21.5V48H86.077V44C85.1103 45.3667 83.827 46.4333 82.227 47.2C80.6603 47.9333 78.9436 48.3 77.077 48.3C73.5436 48.3 70.7603 47.3333 68.727 45.4C66.6936 43.4333 65.677 40.55 65.677 36.75V21.5H70.477V36.2C70.477 38.7667 71.0936 40.7167 72.327 42.05C73.5603 43.35 75.327 44 77.627 44C80.1603 44 82.1603 43.2333 83.627 41.7C85.0936 40.1667 85.827 38 85.827 35.2V21.5H90.627ZM104.365 25.95C105.198 24.4167 106.432 23.25 108.065 22.45C109.698 21.65 111.682 21.25 114.015 21.25V25.9C113.748 25.8667 113.382 25.85 112.915 25.85C110.315 25.85 108.265 26.6333 106.765 28.2C105.298 29.7333 104.565 31.9333 104.565 34.8V48H99.7648V21.5H104.365V25.95ZM144.687 10.9V48H140.087V43.8C139.02 45.2667 137.67 46.3833 136.037 47.15C134.403 47.9167 132.603 48.3 130.637 48.3C128.07 48.3 125.77 47.7333 123.737 46.6C121.703 45.4667 120.103 43.8833 118.937 41.85C117.803 39.7833 117.237 37.4167 117.237 34.75C117.237 32.0833 117.803 29.7333 118.937 27.7C120.103 25.6667 121.703 24.0833 123.737 22.95C125.77 21.8167 128.07 21.25 130.637 21.25C132.537 21.25 134.287 21.6167 135.887 22.35C137.487 23.05 138.82 24.1 139.887 25.5V10.9H144.687ZM131.037 44.1C132.703 44.1 134.22 43.7167 135.587 42.95C136.953 42.15 138.02 41.05 138.787 39.65C139.553 38.2167 139.937 36.5833 139.937 34.75C139.937 32.9167 139.553 31.3 138.787 29.9C138.02 28.4667 136.953 27.3667 135.587 26.6C134.22 25.8333 132.703 25.45 131.037 25.45C129.337 25.45 127.803 25.8333 126.437 26.6C125.103 27.3667 124.037 28.4667 123.237 29.9C122.47 31.3 122.087 32.9167 122.087 34.75C122.087 36.5833 122.47 38.2167 123.237 39.65C124.037 41.05 125.103 42.15 126.437 42.95C127.803 43.7167 129.337 44.1 131.037 44.1Z"
            fill="black"
          />
        </svg>
      </div>
      <div className="w-full block flex-grow lg:flex lg:items-center lg:w-auto">
        <div className="text-sm lg:flex-grow">
          {navItems.map((item) => (
            <a
              href="#link"
              onClick={() => history.push("/" + item.path)}
              className="block mt-4 lg:inline-block lg:mt-0 text-teal-200 hover:text-white mr-4 font-medium"
            >
              {item.title}
            </a>
          ))}
        </div>
      </div>
    </nav>
  );
};

export default NavBar;
