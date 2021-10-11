module.exports = {
  purge: ['./**/*.html','./assets/*'],  
  darkMode: "class", // or 'media' or 'class'
  theme: {
    extend: {},
  },
  variants: {
    extend: {
      hidden: ['dark'],
      display: ['dark'],
      inline: ['dark']
    }
  },
  plugins: [],
}
