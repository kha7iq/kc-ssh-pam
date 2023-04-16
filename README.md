# Keycloak SSH PAM


**kc-ssh-pam** designed to streamline the process of user authentication and enable users to access Linux systems through SSH. The program integrates with Keycloak to obtain a password grant token based on the user's login credentials, including their username and password. If two-factor authentication is enabled for the user, the program supports OTP code as well.

Once the password grant token is obtained, the program verifies it and passes the necessary parameters so that the user can be authenticated via SSH and access the Linux systems.