/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./view/*.templ", "./view/*_templ.go"],
  theme: {
    extend: {
      spacing: {
        'wi-0': '0%',
        'wi-10': '10%',
        'wi-20': '20%',
        'wi-30': '30%',
        'wi-40': '40%',
        'wi-50': '50%',
        'wi-60': '60%',
        'wi-70': '70%',
        'wi-80': '80%',
        'wi-90': '90%',
        'wi-100': '100%',
      }
    },
  },
  plugins: [],
}
