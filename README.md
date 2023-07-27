# Backend API Documentaion

### Installation and Running
This is the backend REST API for Aucto's marketplace, it is currently in `v1`. The technology used is mostly Golang and PostgreSQL. The project has been containerised with the help of Docker and in order to run the project you simply have to:

- `git clone ...`
- `make docker-up`

There is a `MAKEFILE` that simplifies the build and run commands for docker to `docker-up`. The `docker-compose` file contains all the services that are run when running the API. The `Dockerfile` contains the build information of the API. The API is run on `localhost:8080` and the base path is `/api/v1`

Once run, you can run a sanity check by testing the following endpoint: `localhost:8080/api/v1/test/ping`. If setup correctly you should recieve a response in the form of `"messaage":"pong"`.

### API documentations