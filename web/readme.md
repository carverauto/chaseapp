<p align="center">
  <img src="https://chaseapp.tv/chaseapp-lg-bg.jpeg" />
</p>

## Nuxt.js site for ChaseApp

[![](https://img.shields.io/badge/nuxt.js-v2.15.0-04C690.svg)](https://nuxtjs.org)

### Installation

* clone from github
* run `yarn` to install all of your deps
* copy `.env.example` to `.env` and configure it to your likings
* TL;DR
 ```bash
git clone https://github.com/chase-app/web.git; cd web; yarn; cp .env.example .env;
 ```


### Local Environment
* run `yarn dev` in one terminal for our nuxt dev setup

[Setting up SSL for localhost](https://stackoverflow.com/questions/56966137/how-to-run-nuxt-npm-run-dev-with-https-in-localhost)

### Gcloud deployments

Built following the [deployment-cloud-run](https://nuxtjs.org/docs/2.x/deployment/deployment-cloud-run/) guide

Run this once:
```bash
gcloud config set run/region us-central1
```

Generate a build (choose new version number)
```bash
gcloud builds submit --tag gcr.io/chaseapp-8459b/chaseapp:1.0.6 .
```

Deploy the build (choose new version number)
```bash
gcloud run deploy --image=gcr.io/chaseapp-8459b/chaseapp:1.0.6 --platform managed --port 3000 --concurrency 80
```

Be aware that Cloud Run applications will have a default concurrency value of 80 (each container instance will handle up to 80 requests at a time). You can specify the concurrency value this way:

Run the following command to check if the deployment was created successfully:

```bash
gcloud run services list --platform managed
```

To setup CD or .env, check out [build-config](https://cloud.google.com/build/docs/build-config)



