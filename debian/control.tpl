Source: nginx-ad-auth
Section: net
Priority: optional
Maintainer: Ronan LE MEILLAT <ronan@sctg.some.dom>
Build-Depends: debhelper (>= 10), pkg-config
Standards-Version: 4.5.0
Homepage: https://github.com/sctg-development/nginx-ad-auth

Package: nginx-ad-auth
Architecture: {{ ARCH }}
Depends: systemd ${misc:Depends}
Description: Nginx AD Auth helper
 Authenticates users against Active Directory using LDAP