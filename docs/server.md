[<- Go back to "README"](../README.md)

# Server Documentation

Woole allows the creation of custom servers. Before setting up a Woole Server, ensure that the necessary ports (HTTP, HTTPS, and Tunnel) are open and properly configured in the firewall, if applicable. Refer to the server provider’s documentation for specific configuration instructions.

Please note that domains and hosting services are not included with Woole Server.

## Basic Usage

```sh
./woole-server 
```

Output:
```
===============
   HTTP listening on http://{client}.mywebsite.com
  HTTPS listening on https://{client}.mywebsite.com
 Tunnel listening on grpc://mywebsite.com:9653
===============
```

*To provide an HTTPS server, the server must be certified. Consult the [Using HTTPS](#using-https) section for more details.*

## Available Options

| Option                      | Description                                                                |
|-----------------------------|----------------------------------------------------------------------------|
| `-pattern`                  | Set the server hostname pattern. Example: `{client}.mysite.com`            |
| `-http`                     | Port on which the server listens for HTTP requests (default `80`)          |
| `-https`                    | Port on which the server listens for HTTPS requests (default `443`)        |
| `-tunnel`                   | Port on which the gRPC tunnel listens (default `9653`)                     |
| `-seed`                     | Key used to hash the client bearer                                         |
| `-tls-cert`                 | Path to the TLS certificate or fullchain file                              |
| `-tls-key`                  | Path to the TLS private key file                                           |
| `-priv-key`                 | Path to the ECC private key used to validate clients (default disabled)    |
| `-log-level`                | Level of detail for the logs to be displayed (default `INFO`)              |
| `-log-remote-addr`          | Log the request remote address                                             |
| `-tunnel-reconnect-timeout` | Timeout to reconnect the stream when the connection is lost. [Duration format](special-types.md#duration-format) (default `10s`) |
| `-tunnel-request-size`      | Tunnel maximum request size. [Size format](special-types.md#size-format) (default `2GB`, limited by gRPC)  |
| `-tunnel-response-size`     | Tunnel maximum response size. [Size format](special-types.md#size-format) (default `2GB`, limited by gRPC) |
| `-tunnel-response-timeout`  | Timeout to receive a client response. [Duration format](special-types.md#duration-format) (default `10s`)  |
| `-tunnel-connection-timeout`| Timeout for client connections. [Duration format](special-types.md#duration-format) (default `unset`)      |

## Hostname Pattern

The `pattern` is used to define the host format and where the [Client ID](client.md#client-id) will be displayed in the URL. Example, `{client}.pattern.here` will generate URLs such as:
- https://clientid.pattern.here;
- https://test.pattern.here
- https://l2rhwi87aira.pattern.here;

*If using a host, configure it to allow the `*.pattern.here` DNS records.*

### Custom URL Rules

The [Client ID](client.md#client-id) will be used for the first attached tunnel and the subsequents will be appended with a 5 digits hash. The [Client ID](client.md#client-id) will become available again once the tunnel dettach.

If no name is provided, an 8 digits hash will be returned instead.

### Example

Using the server pattern "https://{client}.pattern.here" and the [Client ID](client.md#client-id) `test` will return the following URL:
- https://test.pattern.here, if the name test is not in use right now OR
- https://test-3ld8f.pattern.here, with a 5 digits hash.

Otherwise, if the name is not provided, an 8 digits hash will be used instead:
- https://2hv9e4lf.pattern.here

## Using HTTPS

The HTTPS URL is only available for certified servers. Provide the certification path and the key path using `-tls-cert` and `-tls-key` respectively. The HTTPS port can be changed using the `-https` option.

### Example

```sh
    ./woole-server \
        -tls-cert "/etc/tls/domain/fullchain.pem" \
        -tls-key "/etc/tls/domain/privkey.pem"
```

## Server with ECC Authentication

To use the server authentication, generate a pair of ECC (Elliptic Curve Cryptography) keys. Follow the steps below:

### **Using OpenSSL**

1. **Generate the Private Key**:
   ```sh
   openssl ecparam -genkey -name prime256v1 -noout -out private_key.pem
   ```
   - This command generates a private key using the `prime256v1` curve.

2. **Generate the Public Key**:
   ```sh
   openssl ec -in private_key.pem -pubout -out public_key.pem
   ```
   - This command extracts the public key from the private key.

### **Key Usage**
- **Server**:
  - Use the private key (`private_key.pem`) with the `-priv-key` option to validate client connections.
- **Client**:
  - Share the public key with the allowed clients and use the `-server-key` option to authenticate with the server.
