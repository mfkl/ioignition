FROM golang:1.21.1-alpine as build_stage

# For these to work you need to do two things on the dokku host:
# 1. dokku config:set <app-name> KEY=VALUE
# 2. dokku docker-options:add <app-name> build "--build-arg ARG_NAME=${KEY}
# Where key is the one you've set in step 1.
# https://dokku.com/docs/deployment/builders/dockerfiles/?h=dockerfile#build-time-configuration-variables
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
