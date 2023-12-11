FROM golang:1.21.1-alpine as build_stage

ARG DATABASE_URL
ARG PORT
ARG JWT_SECRET

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

ENV DATABASE_URL=${DATABASE_URL}
ENV PORT=${PORT}
ENV JWT_SECRET=${JWT_SECRET}

# Build server
RUN scripts/buildprod.sh

# Install migration tool
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
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

EXPOSE ${PORT}

# Command to run when starting the container
CMD ["./server"]
