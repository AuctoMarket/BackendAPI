# Backend API Documentaion

### Installation and Running
This is the backend REST API for Aucto's marketplace, it is currently in `v1`. The technology used is mostly Golang and PostgreSQL. The project has been containerised with the help of Docker and in order to run the project you simply have to:

- `git clone ...`
- `make docker-up`

There is a `MAKEFILE` that simplifies the build and run commands for docker to `docker-up`. The `docker-compose` file contains all the services that are run when running the API. The `Dockerfile` contains the build information of the API. The API is run on `localhost:8080` and the base path is `/api/v1`

Once run, you can run a sanity check by testing the following endpoint:

 `localhost:8080/api/v1/tests/ping` 
 
 If setup correctly you should recieve a response in the form of:
 
  `"messaage":"pong"`.

### Schema Documentation

Aucto backend runs a Postgres Database Layer with the following ER Diagram: 

![image info](/docs/Aucto%20DB%20ER%20Diagram.png)

Considerations made are:
- Reducing data dependance using Table Normalisation techniques.

### API Documentation

This project used swagger to document the various api endpoints and the swagger docs can be found at `https://uaw1x43etb.execute-api.ap-southeast-1.amazonaws.com/api/v1/docs/index.html#/`. These API represent the API available in latest stable build.

### Reporting Bugs
If a bug is found in the API, create an issue and tag it as a bug. Make sure to add instructions on how to recreate the bug as well as the expected output and the actual output.

### Availability for dev testing
The API has been hosted as a lambda function with a API gateway that allows users to use a http method to invoke the API. The base URL to do so is `https://uaw1x43etb.execute-api.ap-southeast-1.amazonaws.com/api/v1/docs`