Current repo structure

```
backend-crm/
│
├── cmd/
│   └── main.go
│
│
├── internal/
│   ├── config/
│   │   └── <config logic>
│   ├── database/
│   │   └── mongodb/
│   │       └── <business logic>
│   ├── handler/
│   │   └── <handler logic>
|   |   └── <handler test logic>
│   ├── server/
│       ├── routes.go
│       └── server.go
│
└── pkg/
    └── auth/
        └── <auth logic>
    └── core/
        └── <core business logic>
    └── enum/
        └── <constants logic>

```

make run
make test

Setting up dev environment

```
Step 1: Install Docker

Step 2: Install mongodb community https://www.mongodb.com/docs/manual/tutorial/install-mongodb-community-with-docker/
docker pull mongodb/mongodb-community-server:latest
docker-compose -f docker-compose.dev-env.yml up -d

Step 3: Download MongoDB Compass GUI
https://www.mongodb.com/try/download/compass
Make sure you fill in the username and password under Advanced Connection Options

Step 4: Start the shipments server
make run

Step 5: Adding tests
make test
```

```
Prod env
docker-compose -f docker-compose.prod-env.yml up --build
docker-compose -f docker-compose.prod-env.yml up --no-deps --build app

docker build -t columbus-crm-backend .
aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 767397881306.dkr.ecr.ap-southeast-1.amazonaws.com
docker tag columbus-crm-backend:latest 767397881306.dkr.ecr.ap-southeast-1.amazonaws.com/columbus-crm-backend:latest
docker push 767397881306.dkr.ecr.ap-southeast-1.amazonaws.com/columbus-crm-backend:latest

Yes, the command docker-compose up --no-deps --build app will run the Dockerfile associated with the app service. Here's a breakdown of what happens when you execute this command:

--no-deps: Ensures that only the specified service (app in this case) is started without starting any of its dependencies (like the mongo service).
--build: Forces a rebuild of the Docker image for the app service before starting it. This rebuilds the image based on the instructions in the Dockerfile located in the build context specified in your docker-compose.yml.
app: Specifies that the app service should be started.
```

//Enable change streams for the collection "foo" in database "bar"
db.adminCommand({modifyChangeStreams: 1,
database: "ShipmentsStore",
collection: "FeedEmailCollection",
enable: true});

//Disable change streams on collection "foo" in database "bar"

db.adminCommand({modifyChangeStreams: 1,
database: "ShipmentsStore",
collection: "FeedEmailCollection",
enable: false});
