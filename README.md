# nginx-ad-auth

`nginx-ad-auth` is a Go program that serves as an authentication service for the NGINX email plugin. It authenticates users against Active Directory using LDAP.

## Features

- Listens on a configurable HTTP port
- Authenticates users against Active Directory
- Supports IMAP, POP3, and SMTP protocols
- Configurable via command-line flags or environment variables

## Prerequisites

- Go 1.21 or later
- Access to an Active Directory server

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/nginx-ad-auth.git
   ```

2. Change to the project directory:
   ```
   cd nginx-ad-auth
   ```

3. Build the program:
   ```
   go build -o nginx-ad-auth
   ```

## Usage

Run the program with the following command:

```
./nginx-ad-auth [flags]
```

### Flags

- `--port`: Port to listen on (default: 8080)
- `--ldap-uri`: LDAP URI
- `--ldap-base`: LDAP base
- `--ad-domain`: Active Directory domain
- `--mail-server`: Mail server address
- `--mail-server-port`: Mail server port
- `--help`: Show help message

### Environment Variables

You can also use environment variables instead of flags:

- `NGINX_AUTH_PORT`
- `NGINX_AUTH_LDAP_URI`
- `NGINX_AUTH_LDAP_BASE`
- `NGINX_AUTH_AD_DOMAIN`
- `NGINX_AUTH_MAIL_SERVER`
- `NGINX_AUTH_MAIL_SERVER_PORT`

## Docker

To build and run the Docker image:

1. Build the image:
   ```
   docker build -t nginx-ad-auth .
   ```

2. Run the container:
   ```
   docker run -p 8080:8080 -e NGINX_AUTH_LDAP_URI=ldap://your-ad-server nginx-ad-auth
   ```

## Kubernetes

To deploy on Kubernetes using Helm:

1. Change to the helm chart directory:
   ```
   cd helm/nginx-ad-auth
   ```

2. Install the chart:
   ```
   helm install nginx-ad-auth .
   ```

## License

This project is licensed under the MIT License.
