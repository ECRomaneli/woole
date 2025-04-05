<p align='center'>
    <a href="https://github.com/ECRomaneli/woole" style='text-decoration:none'>
        <img src="https://i.postimg.cc/XJs1WfVc/logo.png" alt='Woole'>
    </a>
</p>
<p align='center'>
    The Wormhole (or just Woole) is an Open-Source reverse-proxy, sniffing, and tunneling tool developed in Go
</p>
<p align='center'>
    <a href="https://github.com/ECRomaneli/woole/tags" style='text-decoration:none'>
        <img src="https://img.shields.io/github/v/tag/ecromaneli/woole?label=version&sort=semver&style=for-the-badge" alt="Version">
    </a>
    &nbsp;
    <a href="https://github.com/ECRomaneli/woole/commits/master" style='text-decoration:none'>
        <img src="https://img.shields.io/github/last-commit/ecromaneli/woole?style=for-the-badge" alt="Last Commit">
    </a>
    &nbsp;
    <a href="https://github.com/ECRomaneli/woole/blob/master/LICENSE" style='text-decoration:none'>
        <img src="https://img.shields.io/github/license/ecromaneli/woole?style=for-the-badge" alt="License">
    </a>
    &nbsp;
    <a href="https://github.com/ECRomaneli/woole/issues" style='text-decoration:none'>
        <img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=for-the-badge" alt="Contributions Welcome">
    </a>
</p>

## Summary

- [How it Works](#how-it-works)
- [Client](#client)
- [Server](#server)
- [Build](#build)
- [Docker](#docker)
- [Releases](#releases)
- [Disclaimer](#disclaimer)

## How it Works

Woole provides two modules: the server and the client. The server sets up an HTTP Tunnel, sends requests to the client, and waits for responses. The client retrieves requests using the configured tunnel, performs reverse-proxy operations, stores the information, and sends responses back to the server. Additionally, the client provides a sniffing tool.

&nbsp;

<p align='center'>
    <a href="https://github.com/ECRomaneli/woole" style='text-decoration:none'>
        <img src="https://i.postimg.cc/VkkFygg7/diagram.png" alt='Diagram'>
    </a>
</p>

## Releases

Pre-built binaries for the client and server are available in the [Releases](https://github.com/ECRomaneli/woole/releases) section. Follow the steps below to download and use them:

### Windows

1. Visit the [Releases Page](https://github.com/ECRomaneli/woole/releases) and download the appropriate binary for your platform (e.g., `woole-windows-x64.zip`).
2. Extract the downloaded file using a compression tool (or just double-click the zip file).
3. Open a Command Prompt or PowerShell in the folder where the binaries were extracted.
4. Run the binaries:
    ```sh
    ./woole.exe --help
    ./woole-server.exe --help
    ```
5. **(Optional)** Add the folder containing the binaries to your PATH environment variable for easier access.

### MacOS/Linux

1. Visit the [Releases Page](https://github.com/ECRomaneli/woole/releases) and download the appropriate binary for your platform (e.g., `woole-linux-x64.zip`, `woole-macos-arm64.zip` or `woole-macos-x64.zip`).
2. Navigate to the download folder and extract the downloaded file:
   ```sh
   unzip woole-<version>-<platform>.zip
   ```
3. Make the binaries executable:
    ```sh
    chmod +x woole woole-server
    ```
4. **(Optional)** Move the binaries to a directory in your `PATH`:
    ```sh
    sudo mv woole woole-server /usr/local/bin/
    ```
5. Run the binaries:
    ```sh
    ./woole --help
    ./woole-server --help

    # If the binaries are in your PATH variable
    woole --help
    woole-server --help
    ```

*Some systems may need administrator permissions to bind well known ports (e.g. 80, 443, etc.)*

## Client

### Basic Usage

#### Using a standalone server (provided by the client)

```sh
./woole -http :80 -proxy https://github.com/ECRomaneli/woole

===============
HTTP URL:  http://localhost:80
Proxying:  https://github.com/ECRomaneli/woole
 Sniffer:  http://localhost:8000
===============
```
*More about standalone config under the [Standalone Mode section](#standalone-mode).*

#### Using the tunnel to connect with an external server

```sh
./woole -proxy https://github.com/ECRomaneli/woole -tunnel grpc://woole.me:8001

===============
 HTTP URL: http://x5ck9p8e.woole.me
HTTPS URL: https://x5ck9p8e.woole.me (Only with TLS configured)
 Proxying: https://github.com/ECRomaneli/woole
  Sniffer: http://localhost:8000
===============
```

*The HTTPS URL requires a certified Server. Otherwise, only the HTTP URL will be displayed.*

### Available Options

| Option              | Description                                                                 |
|---------------------|-----------------------------------------------------------------------------|
| `-client`           | Unique identifier of the client                                            |
| `-http`             | Port to start the standalone server (disables tunnel)                      |
| `-proxy`            | URL of the target server to be proxied (default `:80`)                     |
| `-tunnel`           | URL of the tunnel (default `:8001`)                                        |
| `-custom-host`      | Custom host to be used when proxying                                        |
| `-sniffer`          | Port on which the sniffer is available (default `:8000`)                |
| `-records`          | Max records to store. Use `0` for unlimited (default `1000`)              |
| `-log-level`        | Level of detail for the logs to be displayed (default `INFO`)              |
| `-tls-skip-verify`  | Disables the validation of the integrity of the Server's certificate       |
| `-tls-ca`           | Path to the TLS CA file (only for self-signed certificates)                |

### URL Patterns

All options that requires URLs or ports **MUST** follow one of the patterns below:

- `protocol`://`host`:`port`;
- `protocol`://`host` (default port `80`);
- `host`:`port`, (default protocol `HTTP`);
- `host`, (default protocol `HTTP` and port `80`);
- :`port`, where the colon is important (default protocol `HTTP` and host `localhost`).

### Proxy

Woole is capable to proxy local and online HTTP and HTTPS webservers. Custom names defined using a DNS or `hosts` file are also supported.

Define the URL to proxy using the option `-proxy`. The URL must follows one of [these patterns](#url-patterns):

```sh
./woole -proxy :<port>
```

#### Custom Host

The option `custom-host` can be used along with the proxy to customize the host provided during the HTTP requests.
The default value is the proxy URL.

```sh
./woole -proxy localhost:8080 -custom-host mywebsite.com
```


### Client ID

The `-client` is optional. However, when creating the URL, the provided client ID will be prioritized.
Read more about the URL and how the client is used by the server in the [Server Pattern Section](#pattern).

### Standalone Mode

Standalone mode initiates a self-served HTTP server via the Woole client. In this mode, the client bypasses the need for a tunnel connection, and can be used as a local sniffing tool and reverse proxy. 

To enable it, provide the port using the option `-http`. If a custom name is provided along with the port, the URL will only be accessible through it.

#### Example

```sh
./woole -http :80 -proxy <interval-or-external-url>

# Or, if a custom domain name was configured

./woole -http custom_name:80 -proxy <internal-or-external-url>
```

When using the Standalone mode, the tunnel related options are going to be ignored.

### Tunnel

To use the tunneling tool, a server must be configured and provided using the `tunnel` option.

A single server allows many client connections at the same time. Once configured, copy the Tunnel URL provided by the server and use it as `tunnel`. Consult the [Server Section](#server) for more details on how to create and configure a server. The URL follows [this pattern](#url-patterns) but the `grpc` protocol is used instead of the `http`.

```sh
./woole [...] -tunnel grpc://woole.me:8001
```

If the objective is to only use the sniffing tool and the reverse proxy, without the tunnel, consider using the [Standalone Mode](#standalone-mode).

### Troubleshooting

#### Expired or unsafe cerficate

For servers with expired or unsafe certificates, if trusted by the user, use the option `-tls-skip-verify` to disable the validation of the integrity. Otherwise, the connection with the tunnel will not be possible. This is only required by HTTPS servers.

#### [EXPERIMENTAL] Self-signed Certificates

If the server and the client shares a self-signed certificate, use the `-tls-ca` option to provide the CA file path.

#### Requests not showing

Some browsers and websites utilize efficient caching mechanisms to minimize unnecessary requests. To temporarily disable this caching, disable the cache in your browser's DevTools (Network tab > Check "Disable Cache") and remove any cache headers from website requests. Note that this behavior is not a bug.

### Sniffing Tool

<p align='center'>
    <a href="https://github.com/ECRomaneli/woole" style='text-decoration:none'>
        <img src="https://i.postimg.cc/cL27kFvc/sniffing-tool.png" alt='Sniffing Tool'>
    </a>
</p>

The sniffing tool is accessible through the port configured using the `sniffer` option (default port is available in the [options list](#available-options)). To change the port use:

```sh
./woole [...] -sniffer :9094
```

#### Features:
- Light/Dark Theme;
- Fuzzy Search (status, host, url, name, headers, request body, cookies);
- Media preview (audio, video [chunks are not supported], and images);
- Replay requests with or without changes (with editor);
- Generate request cURLs;
- ACE Editor as viewer for the request and response body;
- Beautify XML, HTML, JSON, javascript, and CSS bodies.

#### Fuzzy Search:

The search uses the pattern `root.parent.child: value` recursively where one or more levels can be used starting from the root parent or not. The value is optional. The root parent is not required, the search can start by any level.

For instance, to search for a specific header called `Content-Type`, the following options are valid:

```
Content-Type
header.Content-Type
response.header.Content-Type
```

and to search for `XML` bodies:

```
response.header.Content-Type: xml
```

Note that the value does not need to match the entire field.

#### Hierarchical Structure

```js
request
├── proto: string (Protocol)
├── method: string (HTTP Verbs)
├── url: string
├── path: string
├── header
│   ├── name_1: string
│   └── name_n: string
└── body: text
response
├── proto: string (Protocol)
├── status: string (e.g. Not Found)
├── code: int (e.g. 404)
├── header
│   ├── name_1: string
│   └── name_n: string
├── body: text
├── elapsed: int
└── serverElapsed: int
```

## Server

With Woole, you can create **YOUR** own server. But before setup your Woole Server, be sure that your server port is open and the firewall (if configurable) has the HTTP, HTTPS and Tunnel port configured. Consult the server provider documentation to know how to configure that. Domains and hostings are not provided by Woole Server.

The https://woole.me website was created to provide a free-to-use Woole Server. Just use the tunnel URL `grpc://woole.me:8001`. The virtual machine is not so powerful so use with moderation. Note that the website disponibility can change without further advise.

### Basic Usage

```sh
    ./woole-server 

    ===============
    HTTP listening on http://{client}
    Tunnel listening on http://{client}:8001
    ===============
```

### Available Options

| Option                  | Description                                                                 |
|-------------------------|-----------------------------------------------------------------------------|
| `-pattern`              | Set the server hostname pattern. Example: `{client}.mysite.com`            |
| `-http`                 | Port on which the server listens for HTTP requests (default `80`)          |
| `-https`                | Port on which the server listens for HTTPS requests (default `443`)        |
| `-tunnel`               | Port on which the gRPC tunnel listens (default `8001`)                     |
| `-key`                  | Key used to hash the bearer                                                |
| `-tls-cert`             | Path to the TLS certificate or fullchain file                              |
| `-tls-key`              | Path to the TLS private key file                                           |
| `-log-level`            | Level of detail for the logs to be displayed (default `INFO`)              |
| `-tunnel-reconnect-timeout` | Timeout to reconnect the stream when the connection is lost (default `10000`ms) |
| `-tunnel-request-size`  | Tunnel maximum request size in bytes (default `math.MaxInt32`)             |
| `-tunnel-response-size` | Tunnel maximum response size in bytes (default `math.MaxInt32`)            |
| `-tunnel-response-timeout` | Timeout to receive a client response (default `20000`ms) 

### Pattern

The `pattern` is used to define the host format and where the [Client ID](#client-id) will be displayed in the URL. Example, `{client}.woole.me` will generate URLs such as:
- client-name-here.woole.me;
- test.woole.me
- l2rhwi87aira.woole.me;

#### Custom URL Rules

The [Client ID](#client-id) will be used for the first attached tunnel and the subsequents will be appended with a 5 digits hash. The [Client ID](#client-id) will become available again once the tunnel dettach.

If no name is provided, an 8 digits hash will be returned instead.

#### Example

Using the server pattern https://{client}.woole.me and the [Client ID](#client-id) `test` will return the following URL:
- https://test.woole.me, if the name test is not in use right now OR
- https://test-3ld8f.woole.me, with a 5 digits hash.

Otherwise, if the name is not provided, an 8 digits hash will be used instead:
- https://2hv9e4lf.woole.me

### Using HTTPS

The HTTPS URL is only available for certified servers. Provide the certification path and the key path using `-tls-cert` and `-tls-key` respectively. The HTTPS port can be changed using `-https`.

## Build

Manually:

```sh
    git clone --depth 1 https://github.com/ecromaneli/woole.git

    # to build the client
    cd woole/cmd/client
    go build -o woole
    chmod +x woole

    # to build the server
    cd woole/cmd/server
    go build -o woole-server
    chmod +x woole-server

```

Now, just run the executable using the options above. You can also use `-help` to see the available options.

## Docker

Woole can be run using Docker for easier setup and usage. The Dockerfile supports building images for both the client and the server. Follow the instructions below to build and run the images.

### Dockerfile Arguments

The Dockerfile accepts the following arguments:

- **`MODULE`**: Specifies which module to build. Possible values are:
  - `client` (default): Builds the Woole client.
  - `server`: Builds the Woole server.
- **`VERSION`**: Specifies the version of the source code to use. Possible values are:
  - Branch (default): Uses the `master` branch as default.
  - Or any specific tag or branch, such as `v1.0.0`.

### Building and Running the Images

#### Server

To build and run the Woole server:

1. Build the server image:
   ```sh
   docker build -t woole-server --build-arg MODULE=server --build-arg VERSION=v1.0.0 -f Dockerfile .
   ```

   - Here, `MODULE=server` specifies that the server will be built.
   - `VERSION=v1.0.0` indicates that version `v1.0.0` of the repository will be used.

2. Run the server:
   ```sh
   docker run --rm -p 8001:8001 -p 80:80 woole-server $ARGS
   ```

   - The server will be available on ports `8001` (tunnel) and `80` (HTTP).
   - Replace `$ARGS` with any additional arguments you want to pass to the server (see the [Server Options](#server) section).

#### Client

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
   - Replace `$ARGS` with any additional arguments you want to pass to the client (see the [Client Options](#client) section).

### Examples

#### Server with default configuration

```sh
 docker build -t woole-server --build-arg MODULE=server -f Dockerfile .
 docker run --rm -p 8001:8001 -p 80:80 woole-server
```

#### Client with a configured tunnel

```sh
 docker build -t woole -f Dockerfile .
 docker run --rm -p 8000:8000 woole -proxy http://localhost:8080 -tunnel grpc://woole.me:8001
```

If the server and client are running in the same machine, remember to put the tunnel URL to a network visible on both containers.

Example:

```sh
# Run the server and export the tunnel and HTTP port
docker run --rm -p 8001:8001 -p 80:80 woole-server

# Access the tunnel through the server exported tunnel port using the host IP address (localhost will not work because the container is isolated)
docker run --rm -p 8000:8000 woole -proxy http://localhost:8080 -tunnel grpc://<host-ip-address>:8001
```

*The docker option `--network host` can also be used. However, it is not recommended for security reasons.*

For more information on available options, refer to the [Client](#client) and [Server](#server) sections.

## Author

- Created by [Emerson Capuchi Romaneli](https://github.com/ECRomaneli) (@ECRomaneli).

## Disclaimer

The Woole project, the woole.me website and all its contributors are not responsible for and do not encourage the use of this tool for any illegal activity. You as the user are solely responsible for its use. Report cybercrimes.

## License

[MIT License](https://github.com/ECRomaneli/woole/blob/master/LICENSE)
