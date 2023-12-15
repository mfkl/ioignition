FROM alpine:3.14

# Set the Current Working Directory inside the container
WORKDIR /app

ARG DATABASE_URL
ARG JWT_SECRET

# Copy the pre-built Go binary into our container.
COPY main .

# Copy the public directory into the container
COPY public public

ENV HOST=0.0.0.0
ENV DATABASE_URL=${DATABASE_URL}
ENV JWT_SECRET=${JWT_SECRET}

EXPOSE 8080

CMD ["./main"]
