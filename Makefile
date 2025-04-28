build:
	npx @tailwindcss/cli -i ./src/html/css/input.css -o ./src/html/css/output.css

dev: 
	# Run Tailwind in watch mode AND Go app together
	npx @tailwindcss/cli -i ./src/html/css/input.css -o ./src/html/css/output.css --watch &
	air
