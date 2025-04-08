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
- [Releases](#releases)
- [Client](#client)
  - [Basic Usage](#basic-usage)
  - [Available Options](#available-options)
  - [Proxy](#proxy)
  - [Client ID](#client-id)
  - [Standalone Mode](#standalone-mode)
  - [Tunnel](#tunnel)
  - [Troubleshooting](#troubleshooting)
  - [Sniffing Tool](#sniffing-tool)
    - [Features](#features)
    - [Fuzzy Search](#fuzzy-search)
    - [Hierarchical Structure](#hierarchical-structure)
- [Server](#server)
  - [Basic Usage](#basic-usage-1)
  - [Available Options](#available-options-1)
  - [Hostname Pattern](#hostname-pattern)
  - [Using HTTPS](#using-https)
- [Build](#build)
- [Docker](#docker)
- [Custom Types](#custom-types)
    - [URL Patterns](#url-patterns)
    - [Duration Format](#duration-format)
    - [Size Format](#size-format)
- [Author](#author)
- [Disclaimer](#disclaimer)
- [License](#license)

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
HTTP URL:  http://localhost
Proxying:  https://github.com/ECRomaneli/woole
 Sniffer:  http://localhost:8000
===============
```
*More about standalone config under the [Standalone Mode section](#standalone-mode).*

#### Using the tunnel to connect with an external server

```sh
./woole -proxy https://github.com/ECRomaneli/woole -tunnel woole.me

===============
 HTTP URL: http://x5ck9p8e.woole.me
HTTPS URL: https://x5ck9p8e.woole.me
 Proxying: https://github.com/ECRomaneli/woole
  Sniffer: http://localhost:8000
===============
```

*The HTTPS URL requires a certified Server. Otherwise, only the HTTP URL will be displayed.*

### Available Options

| Option                  | Description                                                                |
|-------------------------|----------------------------------------------------------------------------|
| `-client`               | Unique identifier of the client                                            |
| `-http`                 | Port to start the standalone server (disables tunnel)                      |
| `-proxy`                | URL of the target server to be proxied (default `80`)                      |
| `-tunnel`               | URL of the tunnel (default `9653`)                                         |
| `-custom-host`          | Custom host to be used when proxying                                       |
| `-sniffer`              | Port on which the sniffer is available (default `8000`)                    |
| `-records`              | Max records to store. Use `0` for unlimited (default `1000`)               |
| `-log-level`            | Level of detail for the logs to be displayed (default `INFO`)              |
| `-tls-skip-verify`      | Disables the validation of the integrity of the Server's certificate       |
| `-tls-ca`               | Path to the TLS CA file (only for self-signed certificates)                |
| `-reconnect-attempts`   | Maximum number of reconnection attempts. Use `0` for infinite (default `5`)|
| `-reconnect-interval`   | Time between reconnection attempts. [Duration format](#duration-format) (default `5s`) |

### Proxy

Woole is capable to proxy local and online HTTP and HTTPS webservers. Custom names defined using a DNS or `hosts` file are also supported.

Define the URL to proxy using the option `-proxy`. The URL must follows one of [these patterns](#url-patterns):

```sh
./woole -proxy <port>
```

#### Custom Host

The option `custom-host` can be used along with the proxy to customize the host provided during the HTTP requests.
The default value is the proxy URL.

```sh
# Proxying "http://localhost:8080" but sending mywebsite.com as header
./woole -proxy 8080 -custom-host mywebsite.com
```


### Client ID

The `-client` is optional. However, when creating the URL, the provided client ID will be prioritized.
Read more about the URL and how the client is used by the server in the [Server Hostname Pattern](#hostname-pattern) section.

### Standalone Mode

Standalone mode initiates a self-served HTTP server via the Woole client. In this mode, the client bypasses the need for a tunnel connection, and can be used as a local sniffing tool and reverse proxy. 

To enable it, provide the port using the option `-http`. If a custom name is provided along with the port, the URL will only be accessible through it.

#### Example

```sh
./woole -http 80 -proxy <internal-or-external-url>

# Or, if a custom domain name was configured

./woole -http custom_name:80 -proxy <internal-or-external-url>
```

When using the Standalone mode, the tunnel related options are going to be ignored.

### Tunnel

To use the tunneling tool, a server must be configured and provided using the `tunnel` option.

A single server allows many client connections at the same time. Once configured, copy the Tunnel URL provided by the server and use it as `tunnel`. Consult the [Server Section](#server) for more details on how to create and configure a server. The URL follows [this pattern](#url-patterns) but the default protocol is `grpc` and the default port is `9653`.

```sh
./woole [...] -tunnel woole.me

# OR, WITH A CUSTOMIZED PORT
./woole [...] -tunnel woole.me:<tunnel-port>
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
./woole [...] -sniffer 9094
```

#### Features
- Light/Dark Theme;
- Fuzzy Search (status, host, url, name, headers, request body, cookies);
- Media preview (audio, video [chunks are not supported], and images);
- Replay requests with or without changes (with editor);
- Generate request cURLs;
- ACE Editor as viewer for the request and response body;
- Beautify XML, HTML, JSON, javascript, and CSS bodies.

#### Fuzzy Search

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

```
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

## Woole.me Server

The https://woole.me website was created to offer a free-to-use Woole Server. Simply connect using the tunnel URL `woole.me`.

Please note that the virtual machine has limited resources, so we kindly ask that you use it in moderation. The server will always run the latest released version of Woole.

Keep in mind that the website’s availability may change without prior notice.

## Local or Hosted Server

Woole allows the creation of custom servers. Before setting up a Woole Server, ensure that the necessary ports (HTTP, HTTPS, and Tunnel) are open and properly configured in the firewall, if applicable. Refer to the server provider’s documentation for specific configuration instructions.

Please note that domains and hosting services are not included with Woole Server.

### Basic Usage

```sh
    ./woole-server 

    ===============
      HTTP listening on http://{client}.custom.pattern
     HTTPS listening on https://{client}.custom.pattern
    Tunnel listening on grpc://10.0.0.7:9653
    ===============
```

*If the server resolves the address to a loopback IP, the resolved IP will be displayed. Otherwise, a default placeholder in the format `grpc://<hostname-or-ip>:9653` will be used. The `hostname`, if any, can be used instead of the IP.*

*To provide an HTTPS server, the server must be certified. Consult the [Using HTTPS](#using-https) section for more details.*

### Available Options

| Option                      | Description                                                                 |
|-----------------------------|-----------------------------------------------------------------------------|
| `-pattern`                  | Set the server hostname pattern. Example: `{client}.mysite.com`            |
| `-http`                     | Port on which the server listens for HTTP requests (default `80`)          |
| `-https`                    | Port on which the server listens for HTTPS requests (default `443`)        |
| `-tunnel`                   | Port on which the gRPC tunnel listens (default `9653`)                     |
| `-key`                      | Key used to hash the bearer                                                |
| `-tls-cert`                 | Path to the TLS certificate or fullchain file                              |
| `-tls-key`                  | Path to the TLS private key file                                           |
| `-log-level`                | Level of detail for the logs to be displayed (default `INFO`)              |
| `-tunnel-reconnect-timeout` | Timeout to reconnect the stream when the connection is lost. [Duration format](#duration-format) (default `10s`) |
| `-tunnel-request-size`      | Tunnel maximum request size. [Size format](#size-format) (default `2GB`, limited by gRPC)  |
| `-tunnel-response-size`     | Tunnel maximum response size. [Size format](#size-format) (default `2GB`, limited by gRPC) |
| `-tunnel-response-timeout`  | Timeout to receive a client response. [Duration format](#duration-format) (default `10s`)      |
| `-tunnel-connection-timeout`| Timeout for client connections. [Duration format](#duration-format) (default `unset`)          |

### Hostname Pattern

The `pattern` is used to define the host format and where the [Client ID](#client-id) will be displayed in the URL. Example, `{client}.pattern.here` will generate URLs such as:
- https://clientid.pattern.here;
- https://test.pattern.here
- https://l2rhwi87aira.pattern.here;

*If using a host, configure it to allow the `*.pattern.here` DNS records.*

#### Custom URL Rules

The [Client ID](#client-id) will be used for the first attached tunnel and the subsequents will be appended with a 5 digits hash. The [Client ID](#client-id) will become available again once the tunnel dettach.

If no name is provided, an 8 digits hash will be returned instead.

#### Example

Using the server pattern "https://{client}.pattern.here" and the [Client ID](#client-id) `test` will return the following URL:
- https://test.pattern.here, if the name test is not in use right now OR
- https://test-3ld8f.pattern.here, with a 5 digits hash.

Otherwise, if the name is not provided, an 8 digits hash will be used instead:
- https://2hv9e4lf.pattern.here

### Using HTTPS

The HTTPS URL is only available for certified servers. Provide the certification path and the key path using `-tls-cert` and `-tls-key` respectively. The HTTPS port can be changed using the `-https` option.

#### Example

```sh
    ./woole-server \
        -tls-cert "/etc/tls/domain/fullchain.pem" \
        -tls-key "/etc/tls/domain/privkey.pem"
```

## Build

Manually:

```sh
    git clone --depth 1 https://github.com/ecromaneli/woole.git

    # to build the client
    go build -o ./bin/woole ./cmd/client
    chmod +x ./bin/woole

    # to build the server
    cd woole/cmd/server
    go build -o ./bin/woole-server ./cmd/server
    chmod +x ./bin/woole-server
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
   docker run --rm -p 9653:9653 -p 80:80 woole-server $ARGS
   ```

   - The server will be available on ports `9653` (tunnel) and `80` (HTTP).
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
 docker run --rm -p 9653:9653 -p 80:80 woole-server
```

#### Client with a configured tunnel

```sh
 docker build -t woole -f Dockerfile .
 docker run --rm -p 8000:8000 woole -proxy http://localhost:8080 -tunnel woole.me
```

If the server and client are running in the same machine, remember to put the tunnel URL to a network visible on both containers.

Example:

```sh
# Run the server and export the tunnel and HTTP port
docker run --rm -p 9653:9653 -p 80:80 woole-server

# Access the tunnel through the server exported tunnel port using the host IP address (localhost will not work because the container is isolated)
docker run --rm -p 8000:8000 woole -proxy http://localhost:8080 -tunnel grpc://<host-ip-address>:9653
```

*The docker option `--network host` can also be used. However, it is not recommended for security reasons.*

For more information on available options, refer to the [Client](#client) and [Server](#server) sections.

## Special Types

### URL Patterns

All options that requires URLs **MUST** follow one of the patterns below:

- `protocol`://`host`:`port`;
- `protocol`://`host` (default port `80`);
- `host`:`port`, (default protocol `HTTP`);
- `host`, (default protocol `HTTP` and port `80`);
- `port`, only digits (default protocol `HTTP` and host `localhost`).

### Duration Format

The **Duration Format** allows you to specify time intervals using a combination of numeric values and time unit qualifiers. This format is used in options like `-reconnect-interval` and other timeout-related configurations.

#### Supported Qualifiers
- `d` - Days
- `h` - Hours
- `m` or `min` - Minutes
- `s` - Seconds
- `ms` - Milliseconds
- `ns` - Nanoseconds

#### Examples
| Input String               | Description                             | Equivalent Duration         |
|----------------------------|-----------------------------------------|-----------------------------|
| `1d`                       | 1 day                                   | 24 hours                    |
| `2h 30m`                   | 2 hours and 30 minutes                  | 2 hours, 30 minutes         |
| `45s`                      | 45 seconds                              | 45 seconds                  |
| `100ms`                    | 100 milliseconds                        | 100 milliseconds            |
| `1h 15min 10s`             | 1 hour, 15 minutes, and 10 seconds      | 1 hour, 15 minutes, 10 secs |
| `0`                        | Zero duration                           | 0                           |

### Size Format

The **Size Format** allows you to specify sizes in bytes using a combination of numeric values and unit qualifiers. This format is used in options like `-tunnel-request-size` and `-tunnel-response-size`.

#### Supported Qualifiers
- `b` - Bytes
- `kb` - Kilobytes (1 KB = 1024 bytes)
- `mb` - Megabytes (1 MB = 1024 KB)
- `gb` - Gigabytes (1 GB = 1024 MB)
- `tb` - Terabytes (1 TB = 1024 GB)

#### Examples
| Input String     | Description  | Equivalent Size                     |
|------------------|--------------|-------------------------------------|
| `1024b`          | 1024 bytes   | 1024 bytes                          |
| `1kb`            | 1 kilobyte   | 1024 bytes                          |
| `2mb`            | 2 megabytes  | 2 * 1024 * 1024 bytes               |
| `1gb`            | 1 gigabyte   | 1 * 1024 * 1024 * 1024 bytes        |
| `1tb`            | 1 terabyte   | 1 * 1024 * 1024 * 1024 * 1024 bytes |

## Author

- Created by [Emerson Capuchi Romaneli](https://github.com/ECRomaneli) (@ECRomaneli).

## Disclaimer

The Woole project, the woole.me website and all its contributors are not responsible for and do not encourage the use of this tool for any illegal activity. You as the user are solely responsible for its use. Report cybercrimes.

## License

[MIT License](https://github.com/ECRomaneli/woole/blob/master/LICENSE)
