FROM golang:1.21.1-alpine as build_stage

# Install curl, bash, node, npm
RUN apk add --no-cache bash curl nodejs npm

# Set the working directory inside the container
WORKDIR /temp

# Copy the source code into the container
COPY . .

# Download all dependencies
RUN go mod download

# Build tailwind css
COPY package.json ./
COPY tailwind.config.js ./
RUN npm install
RUN npm run build-css

# build js script
RUN npm run build-script

# install templ
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN npm run build-templ
COPY view/*_templ.go ./view/

# Copy everything inside public
COPY public ./

# Build server
RUN scripts/buildprod.sh

# Run the migrations script
RUN scripts/migrateup.sh

# Final stage: Run the compiled binary
FROM alpine:latest

WORKDIR /app/

# Copy the pre-built binary file from the previous stage
COPY --from=build_stage /temp/.env .
COPY --from=build_stage /temp/view .
COPY --from=build_stage /temp/public .
COPY --from=build_stage /temp/server .

EXPOSE 8080

# Command to run when starting the container
CMD ["./server"]
