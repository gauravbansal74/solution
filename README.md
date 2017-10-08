# Solution

**Prerequisites**
```
 - GoLang (go1.8.3 darwin/amd64)
 - Redis  (docker image redis:alpine OR Redis server v=3.2.9 bits=64)
 - Git 2.10.1
 - Docker API version 1.32
 - Google API Key
 - Enable Google MAP API using developer console
 ```


 **Build and Run**
  - Build
  ```
  ./make.sh
  ```
  - Dev Build and Run locally (Single Server)
  ```
  ./dev.sh
  OR if you want to use docker-compose.yml
  docker-compose up
  ```
   - Build, Create Docker Image and Push to Docker Hub
  ```
  Please update tag or release before run this script
  ./build.sh
  ```
  - to use latest docker image tag
  ```
  please update docker-compose.yml gauravbansal74/solution:{Tag}
  if we want to run multiple server(horizontal-scaling). Use shared redis server and update run.env file.
  ```

  **Questions? Please email at <gauravbansal74@gmail.com>**