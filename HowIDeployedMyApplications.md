# DeployMate

Create a vps with a public ip.

Make sure your domain name have a A record to your IP and CNAME to itself.

1.1 create ssh key and

1. disable ssh password connection:

```shell
sudo vi /etc/ssh/sshd_config
```

PasswordAuthentication yes

```shell
sudo service ssh restart
```

2. install docker

### Add Docker's official GPG key:

```shell
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
```

### Add the repository to Apt sources:

```shell
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
```

### Install Docker and plugins

```
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

3. Set global env variables

For the treafik configuration and tls we need a email and a domain name.

```shell
export DOMAIN=dog.example.com
export EMAIL=name@email.com
```

4. Run docker compose treafik router

```shell
docker compose -f treafik.yml up
```

Content of `treafik.yml` file.

```yml
version: "3.3"
services:
  traefik:
    image: "traefik:v2.10"
    container_name: "traefik"
    command:
      #- "--log.level=DEBUG"
      - "--api.insecure=false"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.myresolver.acme.tlschallenge=true"
      #- "--certificatesresolvers.myresolver.acme.caserver=https://acme-staging-v02.api.letsencrypt.org/directory"
      - "--certificatesresolvers.myresolver.acme.email=${EMAIL}"
      - "--certificatesresolvers.myresolver.acme.storage=/letsencrypt/acme.json"
    ports:
      - "443:443"
    volumes:
      - "./letsencrypt:/letsencrypt"
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
```

5. Run your service with treafik config

For exemple we take the simple whoami service

Content of `whoami.yml` file.

```yml
version: "3.3"
services:
  whoami:
    image: "traefik/whoami"
    container_name: "whoami"
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.whoami.rule=Host(`whoami.${DOMAIN}`)"
      - "traefik.http.routers.whoami.entrypoints=websecure"
      - "traefik.http.routers.whoami.tls.certresolver=myresolver"
```

## Deploy from a dockerfile to a remote

./deploy --domain inizio.app

## Add build and deploy pipline

To add a ci to build and deploy your app you need to create

- docker registry and add USERNAME and PASSWORD in github ci env
- add your server ssh key to your server

Create a `github/workflows/deploy.yml`

```yml
name: Docker Build and Deploy

on:
  push:
    branches:
      - main

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout Repository
        uses: actions/checkout@v2

      - name: Build Docker Image
        run: docker build -t your-image-name:latest .

      - name: Set up QEMU
        uses: docker/setup-qemu-action@v2

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v1

      - name: Login to Docker Hub
        uses: docker/login-action@v1
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Push Docker Image
        run: |
          docker tag your-image-name:latest your-docker-registry/your-image-name:latest
          docker push your-docker-registry/your-image-name:latest

  deploy:
    runs-on: ubuntu-latest

    needs: build

    steps:
      - name: SSH into Server and Update Containers
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USERNAME }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            cd /path/to/your/docker-compose
            docker-compose -f service.yml pull
            docker-compose -f service.yml stop
            docker-compose -f service.yml rm -f
            sed -i 's/image: your-docker-registry\/your-image-name:.*$/image: your-docker-registry\/your-image-name:latest/' service.yml
            docker-compose -f service.yml up -d
```

```bash
docker run
  -d
  --name traefik
  -p 443:443
  -p 80:80
  -p 8080:8080
  -v $(pwd)/letsencrypt:/letsencrypt
  -v /var/run/docker.sock:/var/run/docker.sock:ro traefik:v2.11
  --api.insecure=true
  --providers.docker=true
  --providers.docker.exposedbydefault=false
  --entrypoints.web.address=:80
  --entrypoints.websecure.address=:443
  --certificatesresolvers.myresolver.acme.tlschallenge=true
  --certificatesresolvers.myresolver.acme.email=clement.chalopin@gmail.com
  --certificatesresolvers.myresolver.acme.storage=/letsencrypt/acme.json
```
