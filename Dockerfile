FROM alpine:3.14

# Set the Current Working Directory inside the container
WORKDIR /app

ARG DATABASE_URL
ARG JWT_SECRET

# Copy binary & assets folder
COPY main .
COPY public public

RUN echo "DATABASE_URL=${DATABASE_URL}" >> .env
RUN echo "JWT_SECRET=${JWT_SECRET}" >> .env

ENV HOST=0.0.0.0

EXPOSE 8080

CMD ["./main"]
