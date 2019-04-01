# Fortune App  ![alt](https://travis-ci.org/thiagotrennepohl/fortune-scrapper.svg?branch=master)

The goal of this app is to consume an API and show random fortune messages to the user.


## API reference

[localhost:4001/docs](localhost:4001/docs) OR [https://app.fortune.opsadventures.com/docs](https://app.fortune.opsadventures.com/docs)


### Running locally
 > Note: to populate the database you can also start the scrapper app, which can be found [here](https://github.com/thiagotrennepohl/fortune-scrapper).

`docker-compose -f dev_fortune_app.yml up -d`

### Running tests

`make local-test`

### Docker build

`make build`

### Docker push
> DOCKER_LOGIN, DOCKER_PASSWORD and DOCKER_HUB_NAMESPACE environments are required

### Running on Travis C.I:
 > Note: Be careful with this command, it might apply changes to your Kubernetes cluster.

This repo already contains a `.travis.yml`

`make ci`

This command will install kubectl, configure  a new Kubeconfig file under `~/.kube` and run `./scripts/deployer.sh`. The `deployer.sh` will attempt to perform a blue green deployment.

 > Note: To the blue green deployment be succesful you must have an existent deployment named in this format <production|staging>-fortune-app-<green|blue>.

In order to this command be succesfull you wil have to set up the following environment variables

| Env              |               Description                |
| ---------------- | :--------------------------------------: |
| K8S_CERT         | Kubernetes CA certificate base64 encoded |
| K8S_CLUSTER_ADDR |          Kubernetes api address          |
| K8S_CLUSTER_NAME |         Kubernetes cluster name          |
| K8S_USERNAME     | Service Account with enough permissions  |
| K8S_CLIENT_CERT  |     Svc Account base64 encoded Cert      |
| K8S_CLIENT_KEY   |      Svc Account base64 encoded Key      |


### Performing the first deploy

 > Note: If you already have cert-manager running on your cluster, let's encrypt certs will be automatically created, if you don't you can follow this [guide](https://cert-manager.readthedocs.io/en/latest/getting-started/index.html).

##### Start the service and ingress rules
`DEPLOYMENT=blue ENVIRONMENT=production NAMESPACE=default envsubst < k8s_fortune_service.yml | kubectl apply -f -`

##### Start the pod
`DEPLOYMENT=blue ENVIRONMENT=production NAMESPACE=default envsubst < k8s_fortune_app.yml | kubectl apply -f -`


#### Updating the API docs

`cd docs && make build`

### For help

`make help`


# Technical Decisions


### Why Golang?

It's my first language, so coding becomes fun and productive, also Golang is very simple, very fast and easy to test.
And not less important the community is great.


### motivations to have a seppareted service for the "scrapper" app
  - In my previous experiences running a scheduled job in the same scope as the main application could result in data duplicity. If the cron job is not configurable working with scalability could be painful.
  - Since it's a small json being returned by the API, I have decided to index all messages to a mongodb. If the source API goes through an outage, my html page will continue to work normally.
  - To avoid repeated messages being added to the database, the app creates a checksum of the message and in case of content change a new message will be added.
  - There's also an option to run the scrapper as a cron job, but the minumum interval is 1 minute.

### CI/CD why I decided to use Travis
  I have been using Jenkins for a long time, and Jenkins helped me to create many cool things, but is not everyone that feels motivated to build pipelines using Groovy or Jenkins pipeline syntax. 
  
  For this test I deployed a Jenkins server, but after a few hours I gave up and configured this project to be built by Travis, because it would take too much time to configure all pipelines and plugins.

  On the other hand, Travis is already connected to Github, Yaml is the default pipeline syntax, it has easy secrets management (Jenkins also has) and it's very close to Gitlab C.I. It was much easier to configure my pipelines and connect to Kubernetes, the best part is that I didn't need to run Jenkins the "Docker in Docker" way.

### Kubernetes

  I started deploying a few instances in AWS with Ansible, I was considering to run everything with Ansible/Teraform and the deploys were going to be using Docker's remote API, but Kubernetes makes easier to scale, manage and deploy new containers, in some of my previous experiences tracking servers FQDN was required in order to deploy new containers, well I don't need to do this with Kubernetes and also blue green deployments can be easily created.
  \< insert kubernetes repo here>

  Kubernetes was deployed using [Kops](https://github.com/kelseyhightower/kubernetes-the-hard-way) and it's strongly recommended to read this [repo](https://github.com/kubernetes/kops) beforing using automatic tools.

  Also the cert-manager (for automatic certificate creation) was very challenging.

  At the moment Kubernetes is being used to run the fortune-app and fortune-scrapper only.

### Deploys

As mentioned in the CI section, deploys are being made by Travis using a dedicated Kubernetes Service Account.

The steps are very simple, the command `make ci` does the following
1. Download and install kubectl binary
2. Create Kubeconfig file under `~/.kube/config`
3. Run `deployer.sh`

The below environments are used to create the Kubeconfig file.

| Env              |               Description                |
| ---------------- | :--------------------------------------: |
| K8S_CERT         | Kubernetes CA certificate base64 encoded |
| K8S_CLUSTER_ADDR |          Kubernetes api address          |
| K8S_CLUSTER_NAME |         Kubernetes cluster name          |
| K8S_USERNAME     | Service Account with enough permissions  |
| K8S_CLIENT_CERT  |     Svc Account base64 encoded Cert      |
| K8S_CLIENT_KEY   |      Svc Account base64 encoded Key      |


The `deployer.sh` script which detects if a `green` or a `blue` "version" of the fortune app is running. after that a new deployment is created, then the script waits 30 seconds for any pre-configured healthchecks (kubernetes healthcheck configuration),after that the number of desired and ready pods is checked. and then the service is updated to route all requests to the newly deployed pod and finally the old deployment is deleted.

### Why MongoDB and Why not to run on Kubernetes?

I decided to use MongoDB because there's no relationship in the data being stored, also because of the befits of great performance to read and write  data(be careful with updates) and some my previous experiences managing mongodb replicasets.

##### Why mongodb is not running on kubernetes?

I know that kubernetes has StatefulSets and is great in keeping services separated by node and online, it just take a little bit of configuration and time and you are able to run databases on K8S.

But I never ran stateful services on kubernetes before and I think it might be dangerous if not configured very carefully, because it's not a machine or volume you can simply delete and deploy again, sometimes many services are relying on the same mongodb. And not less important, replicaset elections can be tricky, so you need to be careful to not ruin your data.

Keeping mongodb in a isolated machine with ssd or io1 diks is safe, because you can easily set a disk for journal, one for logs and another one for data (to avoid write and read concurrency).

### Prometheus and Grafana

They are also stateful applications, and for this reason they are running in a separated instance on AWS as well, Prometheus configurarion is being created by Ansible and currently all scrape configs are "automated".

Kubernetes scrape config is consuming Kubernetes API (also a service account for this) and MongoDB exporter and NodeExporter are being "discovered" using Ec2 discovery configurations.


# Fortune Scrapper

More info about the project can be found [here](https://github.com/thiagotrennepohl/fortune-scrapper)


# Fortune Infra / Ansible

More info about infrastructure can be found [here](https://github.com/thiagotrennepohl/fortune-infra)


# What I would like to have done

- First of all, deploying a **Graylogs** or have used a managed logs service. I've been using **Graylogs** for more than 3 year and it's very easy to use, I didn't hav a very good expeience using ELK or Logentries, but they are nice as well. With Graylogs you can easily create alerts for slack, graphs, filter incoming logs, create extraction rules and also processing pipelines, it really makes easier to read logs.

- I woud like to have created a bot using Hubot and connect him to a chat platform (I already did this with RocketChat, Slack and mattermost) to automate some tasks, like triggering a C.I job or a deploy or even creation of review apps.

- Adding a prometheus metrics endpoint to fortune app or add tracing using Jaeger, because it's easier to track metrics from the application, I mean code metrics like execution time, database wait time and so on.

- Adding performance tests to C.I pipeline  to find load breakpoints, memory leaks and would be nice to see difference between memory snapshots.

- Static code analysis is a nice to have!

- Testing fresh docker build is important as well!

- Testing ansible roles using Molecule would be nice.


# What needs improvement?

- The deployer script requires to have an existent deployment in order to the script be successful, also it doesn't clean broken deployments.

- Ansible shouldn't be deploying instances on AWS, Terraform is better to do this

- A default OS snapshot with everything I need to run all sorts of apps, in order to have a faster instance creation.

- Better Kubernetes roles.

- Use Hashicorp Vault for better secret management.

- Create more tests scenarios for both apps.

- Create a rate limit for saving and retrieving messages.

- Create automated backups for mongodb, prometheus and Grafana.
  
- Create alerts for slack and pushover.




