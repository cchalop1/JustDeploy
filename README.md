# JustDeploy (work in progress)

JustDeploy is an open-source application designed to simplify the deployment process for developers. It allows you to deploy your simple applications on your own server with ease. JustDeploy handles server connection, Docker installation, and secure certificate generation, making the deployment process seamless and efficient.

JustDeploy is that it **doesn't install anything on your server other than your application.** This makes it an ideal choice for small servers with limited resources or development environments

![Screenshot4](https://raw.githubusercontent.com/cchalop1/JustDeploy/main/images/Screenshot4.png)

## Install

```bash
curl -fsSL https://justdeploy.app/install | bash
```

## Usage

To create a new deployment you must run the justdeploy command at the bottom of your project where your dockerfile is.

```bash
justdeploy
```

JustDeploy will create a `./just` deploy folder at the root of your project to be able to store the deployment information and the certificates that certify to communicate with the docker engine

For now JustDeploy only support debian base VMs.

After running `justdeploy`, follow the prompts to connect to your server and deploy your application.

it should open your browser on this page.

In this step you can connect you server and when you click on the button JustDeploy install and setup everything he needs.

![Screenshot](https://raw.githubusercontent.com/cchalop1/JustDeploy/main/images/Screenshot.png)

In this 2nd step you can configure your application with name and enable tls if you want.

![Screenshot2](https://raw.githubusercontent.com/cchalop1/JustDeploy/main/images/Screenshot2.png)

### Prerequisites

- Git
- Make
- golang

Before deploying your application using JustDeploy, it is highly recommended that a domain name is associated with your server. This is a necessary step to generate the certificate, even if you disable

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

### Build in local

1. Clone the repository

```bash
git clone https://github.com/cchalop1/JustDeploy.git
```

2. Navigate to the project directory

```bash
cd JustDeploy
```

3. Run Make to build the project

```bash
make
```

4. Run JustDeploy

```bash
./bin/just-deploy
```

## License

This project is licensed under the AGPLv3 License - see the [LICENSE](LICENSE) file for details.

## Contact

For any questions or suggestions, please open an issue or contact the project maintainer (cchalop1).

<!--
## Features

- Connects to your server
- Installs Docker
- Generates secure certificates
- Deploys your application

## TODO

- [x] password auth
- [x] DNS setting process
- [x] persistent data
- [x] buttons on the run part
- [x] logs of the containers
- [x] embed the web build in gobinary
- [x] packages for release it
- [x] install script
- [ ] git hooks post-commit
- [ ] socket
<!-- - [ ] update status with -->

<!-- - [ ] Usage graph on the sucess deploy page
- [ ] From github Url -->
