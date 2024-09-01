# nginx-ad-auth

`nginx-ad-auth` is a Go-based authentication service for the NGINX email proxy, allowing seamless authentication of users against Active Directory using LDAP. It integrates easily with NGINX to secure email services (IMAP, SMTP, POP3), leveraging existing AD infrastructures.

## Features

- **Easy Integration:** Connects with NGINX mail proxy for seamless user authentication.
- **Supports Multiple Protocols:** IMAP, POP3, and SMTP protocols supported for full compatibility.
- **Active Directory Authentication:** Authenticate users against AD using LDAP.
- **Flexible Configuration:** Configure through command-line flags or environment variables.
- **Lightweight:** Minimal dependencies, runs as a standalone service.

## Table of Contents

- [nginx-ad-auth](#nginx-ad-auth)
  - [Features](#features)
  - [Table of Contents](#table-of-contents)
  - [TD;DR](#tddr)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
    - [Test](#test)
  - [Usage](#usage)
    - [Flags](#flags)
    - [Environment Variables](#environment-variables)
  - [Docker](#docker)
  - [Kubernetes](#kubernetes)
  - [Using the Helm Chart](#using-the-helm-chart)
  - [Configuring NGINX as an Email Proxy](#configuring-nginx-as-an-email-proxy)
  - [License](#license)
    - [Key points of the AGPLv3](#key-points-of-the-agplv3)
  - [Contributing](#contributing)
  - [Support](#support)

## TD;DR

You can run `nginx-ad-auth` using Docker in just a few steps:

```bash
docker run -p 8080:8080 \
   -e NGINX_AUTH_LDAP_URI="ldap://your-ad-server" \
   -e NGINX_AUTH_LDAP_BASE="dc=your,dc=domain" \
   -e NGINX_AUTH_AD_DOMAIN="your-domain" \
   -e NGINX_AUTH_MAIL_SERVER="your-mail-server" \
   -e NGINX_AUTH_MAIL_SERVER_PORT=143 \
   sctg/nginx-ad-auth
```

## Prerequisites

- Go 1.21 or later [(Go installation guide)](https://golang.org/doc/install)
- Access to an Active Directory server
- Docker installed for the Docker setup (optional).

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/nginx-ad-auth.git
   ```

2. Change to the project directory:

   ```bash
   cd nginx-ad-auth
   ```

3. Build the program:

   ```bash
   go build -o nginx-ad-auth
   ```

### Test

For testing you can use the provided test file:

```bash
./nginx-ad-auth -ad-domain ADDOMAIN -ldap-base "dc=ADDOMAIN,dc=WINDOWS" -ldap-uri "ldap://server.addomain.windows" -mail-server 192.168.1.1 -mail-server-port 143 -port 8080
VALIDUSER="myuser" CORRECTPASSWORD="mypassword" tests/test-nginx-ad-auth.sh
```

## Usage

Run the program with the following command:

```bash
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

1. (Optional) Build the image:

   ```bash
   docker build -t nginx-ad-auth .
   ```

2. Run the container:

   ```bash
   docker run -p 8080:8080 -e NGINX_AUTH_LDAP_URI=ldap://your-ad-server -e NGINX_AUTH_LDAP_BASE="dc=your,dc=domain" -e NGINX_AUTH_AD_DOMAIN=your-domain -e NGINX_AUTH_MAIL_SERVER="your-mail-server" -e NGINX_AUTH_MAIL_SERVER_PORT=143 sctg/nginx-ad-auth
   ```

## Kubernetes

To deploy on Kubernetes using Helm:

## Using the Helm Chart

To deploy the `nginx-ad-auth` service using the provided Helm chart, follow these steps:

1. First, ensure you have Helm installed on your local machine and configured to work with your Kubernetes cluster.

2. Update the `values.yaml` file in the `helm/nginx-ad-auth` directory to match your environment. Pay special attention to the following fields:
   - `image.repository`: Update this to your Docker registry if you've pushed a custom image.
   - `env`: Update the environment variables to match your Active Directory and mail server configuration.

3. From the root of the project, run:

   ```bash
   helm install nginx-ad-auth ./helm/nginx-ad-auth
   ```

4. To upgrade an existing deployment with new values:

   ```bash
   helm upgrade nginx-ad-auth ./helm/nginx-ad-auth
   ```

5. You can customize the installation by overriding values:

   ```bash
   helm install nginx-ad-auth ./helm/nginx-ad-auth --set replicaCount=3
   ```

Remember to configure your NGINX Ingress or other ingress controller to route traffic to the `nginx-ad-auth` service.

## Configuring NGINX as an Email Proxy

To configure NGINX as an email proxy to a mail server hosted in a private network, you can use the following NGINX configuration:

```nginx
mail {
    server_name mail.example.com;
    auth_http localhost:8080/auth;

    server {
        listen 993 ssl;
        protocol imap;
        ssl_certificate /path/to/your/certificate.crt;
        ssl_certificate_key /path/to/your/certificate.key;
        imap_capabilities "IMAP4rev1" "UIDPLUS";
    }
}
```

This configuration does the following:

- Sets up NGINX to listen on port 993 for IMAPS connections.
- Uses the `nginx-ad-auth` service running on `localhost:8080` for authentication.
- Proxies authenticated connections to the internal mail server at 192.168.1.1:143.
- Enables SSL for both the client connection and the proxy connection to the internal server.

Remember to replace `/path/to/your/certificate.crt` and `/path/to/your/certificate.key` with the paths to your SSL certificate and key files. Also, ensure that the `auth_http` URL matches the location where your `nginx-ad-auth` service is running.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPLv3).

### Key points of the AGPLv3

1. Source Code: You must make the complete source code available when you distribute the software.
2. Modifications: If you modify the software, you must release your modifications under the AGPLv3 as well.
3. Network Use: If you run a modified version of the software on a server and allow users to interact with it over a network, you must make the source code of your modified version available.
4. No Additional Restrictions: You cannot impose any further restrictions on the recipients' exercise of the rights granted by the license.

For the full license text, see the [LICENSE](LICENSE.md) file in the project repository or visit [GNU AGPL v3.0](https://www.gnu.org/licenses/agpl-3.0.en.html).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Support

If you encounter any problems or have any questions, please open an issue in the GitHub repository.
