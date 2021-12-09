module.exports = {
  content: ["./public/**/*.html", "./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    fontFamily: {
      sans: [
        "Libre Franklin",
        "ui-sans-serif",
        "system-ui",
        "-apple-system",
        "BlinkMacSystemFont",
        '"Segoe UI"',
        "Roboto",
        '"Helvetica Neue"',
        "Arial",
        '"Noto Sans"',
        "sans-serif",
        '"Apple Color Emoji"',
        '"Segoe UI Emoji"',
        '"Segoe UI Symbol"',
        '"Noto Color Emoji"',
      ],
    },
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
  plugins: [require("@tailwindcss/forms")],
};
