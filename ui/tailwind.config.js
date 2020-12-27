module.exports = {
  purge: [],
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
