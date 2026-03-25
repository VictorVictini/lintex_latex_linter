# Installation
- go to project folder
- run `go build .`
- run `./lintex` (or equivalent for operating system)
- go to `http://localhost:8080/`
- requires 'lualatex' CLI

for css setup, refer to tailwindcss CLI documentation (this require npm/npx)
- `npm install tailwindcss @tailwindcss/cli`
- `@import "tailwindcss";` at the top of the affected .css files
- `npx @tailwindcss/cli -i ./styles/main.css -o ./styles/output.css`
    - swap `main.css` with the affected file within the `styles` folder
    - if there are several, changing the name of the output css file is necessary to ensure all can be included
- `<link href="/styles/output.css" rel="stylesheet">` within the .html webpages to include the generated output
    - swap `output.css` as necessary

- now requires flowbite installed for handling css