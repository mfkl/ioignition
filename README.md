<div id="header" align="left">
  <img src="https://github.com/Cijin/ioignition/assets/1990966/0c01de29-b659-4640-81cb-97235dae6bec" width="100"/>
</div> 

# ioignition
Privacy first, open, and simple analytics app

## Add secrets
Copy `env.example` into `.env` and add secrets required to run the app

## Running the app locally
* Clone the repo
* Start services: `docker compose -d up`
* Stop services: `docker compse -d down`
* If you want to see the logs skip `-d` and run docker compose
* Not neccessary but recommended:
  * Build Js Script: `npm run build-script`
  * Build tailwind: `npm run build-css`
* To start server: `go run main.go` OR `go build -o main && ./main` -- To run as binary

## Running development server
* Follow steps above for starting services
* There are watch scripts to watch for changes in the package.json
* To watch changes for `.go` files, install air `go install github.com/cosmtrek/air@latest`
* Then run `air`, the config already exists so you don't need to init air
* *NOTE*: the installation of `air` can be skipped, but remember to restart go server everytime you make a change `go run main.go`
