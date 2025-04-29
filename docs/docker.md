[<- Go back to "README"](../README.md)

# Docker Documentation

Woole can be built and run using Docker for easier setup and usage. The Dockerfile supports building images for both the client and the server. The [Dockerfile](https://github.com/ECRomaneli/woole/blob/master/docker/Dockerfile) is available under the `docker` folder in the root path of the project.

## Dockerfile Arguments

The Dockerfile accepts the following arguments:

- **`MODULE`**: Specifies which module to build. Possible values are:
  - `client` (default): Builds the Woole client.
  - `server`: Builds the Woole server.
- **`VERSION`**: Specifies the version of the source code to use. Possible values are:
  - Branch (default): Uses the `master` branch as default.
  - Or any specific tag or branch, such as `v1.0.0`.

## Building and Running the Images

Download the [Dockerfile](https://github.com/ECRomaneli/woole/blob/master/docker/Dockerfile) or copy the code to a local file, open the terminal, go to the folder where the [Dockerfile](https://github.com/ECRomaneli/woole/blob/master/docker/Dockerfile) is located, and follow the steps below.

### Client

To build and run the Woole client:

1. Build the client image:
   ```sh
   docker build -t woole --build-arg VERSION=v1.0.0 -f Dockerfile .
   ```

   - Here, `MODULE=client` is the default value, so it does not need to be specified.
   - `VERSION=v1.0.0` indicates that version `v1.0.0` of the repository will be used.

2. Run the client:
   ```sh
   docker run --rm -p 8000:8000 woole $ARGS
   ```

   - The client will be available on port `8000` (sniffer/dashboard).
   - Replace `$ARGS` with any additional arguments you want to pass to the client (see the [Client Options](client.md#available-options) section).

### Server

To build and run the Woole server:

1. Build the server image:
   ```sh
   docker build -t woole-server --build-arg MODULE=server --build-arg VERSION=v1.0.0 -f Dockerfile .
   ```

   - Here, `MODULE=server` specifies that the server will be built.
   - `VERSION=v1.0.0` indicates that version `v1.0.0` of the repository will be used.

2. Run the server:
   ```sh
   docker run --rm -p 9653:9653 -p 80:80 woole-server $ARGS
   ```

   - The server will be available on ports `9653` (tunnel) and `80` (HTTP).
   - Replace `$ARGS` with any additional arguments you want to pass to the server (see the [Server Options](server.md#available-options) section).

### Examples

#### Server with default configuration

```sh
docker build -t woole-server --build-arg MODULE=server -f Dockerfile .
docker run --rm -p 9653:9653 -p 80:80 woole-server
```

#### Client with a configured tunnel

```sh
docker build -t woole -f Dockerfile .
docker run --rm -p 8000:8000 woole -proxy http://localhost:8080 -tunnel woole.me
```

*If the server and client are running in the same machine, remember to put the tunnel URL to a network visible on both containers. The docker option `--network host` can also be used. However, it is not recommended for security reasons.*

For more information on available options, refer to the [Client](client.md#available-options) and [Server](server.md#available-options) sections.