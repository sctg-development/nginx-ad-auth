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

## Using the Helm Chart

To deploy the `nginx-ad-auth` service using the provided Helm chart, follow these steps:

1. First, ensure you have Helm installed on your local machine and configured to work with your Kubernetes cluster.

2. Update the `values.yaml` file in the `helm/nginx-ad-auth` directory to match your environment. Pay special attention to the following fields:
   - `image.repository`: Update this to your Docker registry if you've pushed a custom image.
   - `env`: Update the environment variables to match your Active Directory and mail server configuration.

3. From the root of the project, run:
   ```
   helm install nginx-ad-auth ./helm/nginx-ad-auth
   ```

4. To upgrade an existing deployment with new values:
   ```
   helm upgrade nginx-ad-auth ./helm/nginx-ad-auth
   ```

5. You can customize the installation by overriding values:
   ```
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