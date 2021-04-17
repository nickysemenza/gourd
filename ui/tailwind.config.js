module.exports = {
  mode: "jit",
  purge: ["./public/**/*.html", "./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {},
  },
  variants: {
    backgroundColor: ["responsive", "odd", "hover", "focus", "disabled"],
    borderStyle: ["responsive", "hover", "focus"],
    borderColor: [
      "responsive",
      "hover",
      "focus",
      "active",
      "group-hover",
      "disabled",
    ],
  },
  plugins: [],
};
