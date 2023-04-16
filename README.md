<h2 align="center">
  <br>
  <p align="center"><img width=30% src="https://raw.githubusercontent.com/kha7iq/kc-ssh-pam/master/.github/img/logo.png"></p>
</h2>

<h4 align="center">Keycloak SSH PAM</h4>

<p align="center">
   <a href="https://github.com/kha7iq/kc-ssh-pam/releases">
   <img alt="Release" src="https://img.shields.io/github/v/release/kha7iq/kc-ssh-pam">
   <a href="https://goreportcard.com/report/github.com/kha7iq/kc-ssh-pam">
   <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/kha7iq/kc-ssh-pam">
   <a href="https://github.com/kha7iq/kc-ssh-pam/issues">
   <img alt="GitHub issues" src="https://img.shields.io/github/issues/kha7iq/kc-ssh-pam?style=flat-square&logo=github&logoColor=white">
   <a href="https://github.com/kha7iq/kc-ssh-pam/blob/master/LICENSE.md">
   <img alt="License" src="https://img.shields.io/github/license/kha7iq/kc-ssh-pam">

</p>

<p align="center">
  <a href="#install">Install</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#show-your-support">Show Your Support</a>
</p>

# Keycloak SSH PAM

**kc-ssh-pam** designed to streamline the process of user authentication and enable users to access Linux systems through SSH. The program integrates with Keycloak to obtain a password grant token based on the user's login credentials, including their username and password. If two-factor authentication is enabled for the user, the program supports OTP code as well.

Once the password grant token is obtained, the program verifies it and passes the necessary parameters so that the user can be authenticated via SSH and access the Linux systems.

```
❯ kc-ssh-pam --help   
Usage: kc-ssh-pam USERNAME PASSWORD/[OTP]

Generates a password grant token from Keycloak for the given user.

Options:
  -h, --help              Show this help message and exit
  -v, --version           Show version information

Notes:
  The program requires a configuration file named 'config.toml' to be present in the 
  current directory , installation directory, or in '/etc/config.toml', or in 
  '$HOME/.config/config.toml', in that order.

  In addition to defaults, all configuration parameters can also be provided through environment variables.

  KC_SSH_REALM   KC_SSH_ENDPOINT   KC_SSH_CLIENTID  
  KC_SSH_CLIENTSECRET  KC_SSH_CLIENTSCOPE
  
  To use the program, you must create a client in Keycloak and provide the following 
  information in the configuration file: realm, endpoint, client ID, client secret, and 
  client scope is optional.

Arguments:
  USERNAME                The username of the user is taken from $PAM_USER environment variable
  PASSWORD                The password of the user is taken from stdIn
  OTP                     (Optional) The OTP code if two-factor authentication is enabled i.e (password/otp)

  EXAMPLE                 (With otp): echo testpass/717912 | kc-ssh-pam (Only Password): echo testpass | kc-ssh-pam
```



## Configuration
  The program requires a configuration file named 'config.toml' to be present in the 
  current directory , installation directory, or in '/etc/config.toml', or in 
  '$HOME/.config/config.toml', in that order.
  ```
realm = "ssh-demo"
endpoint = "https://keycloak.example.com"
clientid = "keycloak-client-id"
clientsecret = "MIKEcHObWmI3V3pF1hcSqC9KEILfLN"
clientscop = "openid"

  ```



## Contributing

Contributions, issues and feature requests are welcome!<br/>Feel free to check
[issues page](https://github.com/kha7iq/kc-ssh-pam/issues). You can also take a look
at the [contributing guide](https://github.com/kha7iq/kc-ssh-pam/blob/master/CONTRIBUTING.md).

## Show your support

Give a ⭐️  if you like this project!

Fork it ⚙️

Make it better 🕶️