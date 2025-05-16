<h1 align="center">JustDeploy 🛵</h1>
<p align="center">Deploy what you want where you want !</p>

<p align="center">
  <a href="https://www.producthunt.com/posts/justdeploy-2?embed=true&utm_source=badge-featured&utm_medium=badge&utm_souce=badge-justdeploy&#0045;2" target="_blank">
    <img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=478298&theme=light" alt="JustDeploy - Deploy&#0032;your&#0032;app&#0032;with&#0032;no&#0032;lock&#0045;in&#0032;&#0038;&#0032;extra&#0032;cost | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" />
  </a>
</p>

![ScreenshotMain](https://raw.githubusercontent.com/cchalop1/JustDeploy/main/images/ScreenProjectMode.png)

<br>

JustDeploy is a PaaS tool designed to simplify the lives of developers. Install it on your server to easily deploy your projects and databases. It's based on Docker, so JustDeploy fetches your GitHub repository and deploys your application using your Docker and Docker Compose configurations. All while deploying to any VPS of your choice. Forget about vendor lock-in—deploy wherever you want.

<br>

## Features

- 🚀 **Deploy full-stack apps in one click**
- ⚙️ **Manage environments seamlessly**
- 🗂️ **Project Management**
- 🌐 **Deploy on any VPS**
- 🔓 **No vendor lock-in—deploy wherever you want**
- 🐙 **Github integration push to deploy**

<br>

## ![Screenshot](https://raw.githubusercontent.com/cchalop1/JustDeploy/main/images/ScreenCreateServices.png)

## Prerequisites

Before getting started, your server need to be debian base or have docker.

---

## Installation

To install JustDeploy on your server, run the following command:

Your server needs to be running Debian.

```bash
curl -fsSL https://get.justdeploy.app | bash
```

If you prefer, you can also run JustDeploy in Docker. With the following command, you will pull and run JustDeploy.

```bash
docker run -d --name justdeploy \
  -p 5915:5915 \
  -v /var/lib/justdeploy:/app/data \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --restart=unless-stopped \
  cchalop1/justdeploy
```

---

## Usage

After installation, JustDeploy runs on your server and you can access it by clicking on the URL displayed at the end of the installation process. You can then access the web interface from your local computer.

With this, you'll be able to:

- Connect your GitHub account
- Set up new databases in just one click
- Manage and configure environments for your applications
- Deploy your project to any VPS of your choice without being locked into a particular provider
- Enjoy zero-config deployment as JustDeploy extracts information from your GitHub repo to deploy your full-stack application

---

## Why Choose JustDeploy?

- **Full Control**: Deploy where you want without extra costs or limitations.
- **Developer-Friendly**: Intuitive and designed to simplify everyday tasks for developers.

---

## Join the Community

- ⭐ Star the repository on [GitHub](https://github.com/cchalop1/JustDeploy)
- 📖 Follow [discussions](https://github.com/cchalop1/JustDeploy/discussions)
- 🐛 Report issues [here](https://github.com/cchalop1/JustDeploy/issues)
- 💬 Follow me on [Twitter](https://x.com/ChalopinClement)
- 🎮 Join our [Discord](https://discord.gg/RteyWyKjz4)

---

## License

This project is licensed under the AGPL-3.0 License. See the [LICENSE](https://github.com/cchalop1/JustDeploy/blob/main/LICENSE) file for more details.
