# Fortune App

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

###### Start the service and ingress rules
`DEPLOYMENT=blue ENVIRONMENT=production NAMESPACE=default envsubst < k8s_fortune_service.yml | kubectl apply -f -`

##### Start the pod
`DEPLOYMENT=blue ENVIRONMENT=production NAMESPACE=default envsubst < k8s_fortune_app.yml | kubectl apply -f -`


#### Updating the API docs

`cd docs && make build`

### For help

`make help`


# Technical Decisions


### Why Golang?

It's my first language, so coding becomes fun and "produtivo", also Golang is very simple, very fast and easy to test.
And not less important the community is great.


### motivations to have a seppareted service for the "scrapper" app
  - In my previous experiences running a scheduled task in the same scope as the main application could result in data duplicity if not treated, then working with scalability  could be painful.
  - Since it's a Small json being returned by the API I have decided to index all messages to a mongodb. In case of API failure my html page will continue to work normally.
  - There's no need to add all messages everytime the app starts or pulls a repeated message, so i decided to create a checksum of the message and in case of content change a new message will be added to the database as well, also to avoid duplicated messages.
  - There's also an option to run the scrapper as a cron job, but the minumum interval is 1 minute.

### CI/CD why I decided to use Travis
  I have been using Jenkins for a long time, and I done many cool things using Jenkins, but few people feel motivated to build pipelines using Groovy or Jenkins pipeline syntax. for this test I deployed a Jenkins server, but gave up because it would take too much time.

  On the other hand, Travis is already conected to Github, Yaml is the default pipeline syntax, it has easy secrets management (Jenkins also has), so it's very close to Gitlab C.I and was much easier to configure my pipelines and connect to Kubernetes and I didn't need to run Jenkins the "Docker in Docker" way.

### Kubernetes

  I started deploying a few machines in AWS with Ansible to run everything and I was going to deploy my services using Docker's remote API, but Kubernetes makes easier to scale, manage and deploy new containers, in some of my previous experiences i had to track servers FQDN in order to deploy new containers, well I don't need to do this with Kubernetes and also i can easily perform blue green deployments.
  \< insert kubernetes repo here>

  Kubernetes was deployed using [Kops](https://github.com/kelseyhightower/kubernetes-the-hard-way) but i strongly recommend to read this [repo](https://github.com/kubernetes/kops) beforing using automatic tools.

  Also was challenging  configuring cert-manager to issue https certificates automatically

  At the moment Kubernetes is being used to run the fortune-app and fortune-scrapper only

### Deploys

As mentioned in the CI section, deploys are being made by Travis using a dedicated Service Account.

The steps are very simple, I download the kubectl to travis machine, create the config based on the environments below:

| Env              |               Description                |
| ---------------- | :--------------------------------------: |
| K8S_CERT         | Kubernetes CA certificate base64 encoded |
| K8S_CLUSTER_ADDR |          Kubernetes api address          |
| K8S_CLUSTER_NAME |         Kubernetes cluster name          |
| K8S_USERNAME     | Service Account with enough permissions  |
| K8S_CLIENT_CERT  |     Svc Account base64 encoded Cert      |
| K8S_CLIENT_KEY   |      Svc Account base64 encoded Key      |


And then I run the `deployer.sh` script which detects if a `green` or a `blue` "version" of the fortune app is running, it then creates a new deployment, after k8s performs the healthcheck defined in the Deployment spec the script checks if the number of desired pods are online, then i update the service to route all requests to the newly deployed pod and finally the old deployment is deleted.

### Why MongoDB and Why not to run on Kubernetes?

I decided to use MongoDB because there's no relationship in the data I'm storing, also it has great performance to read and write (be careful with updates) and I have experience managing mongodb replicasets.

##### Why mongodb is not running on kubernetes?

I know that kubernetes has StatefulSets and is great in keeping services separated by node and online, it just take a little bit of configuration and time.

But I never ran stateful services on kubernetes before and I think it might be dangerous if not configured very carefully, because it's not a machine or volume you can siimply delete and deploy again, sometimes many services are relying on the same mongodb and replicaset elections can be tricky, so you need to be careful with data.

Keeping mongodb in a isolated machine with ssd or io1 diks, a disk for journal, another one for logs and another one for data (to avoid write and read concurrency) seems very safe to me.

### Prometheus and Grafana

They are also stateful applications, and for this reason they are running in a separated instance on AWS as well, Prometheus configurarion is being created by Ansible and currently all scrape configs "automated".

Kubernetes scrape config is consuming Kubernetes api (also a service account for this) and MongoDB exporter and NodeExporter are being "discovered" using Ec2 discovery configurations.

### Why Ansible



# What I would like to have done

- First of all, deploying a **Graylogs** or used a managed logs service, I've been using **Graylogs** for a long time and it's great, I didn't hav a great expeience using ELK or Logentrie, but they are nice as well, but you can easily create alerts for slack, graphs, filter incoming logs, create extraction rules and also processing pipelines, it really makes easier to read logs.

- Creation a bot using hubot and connect to a chat platform (I have experience with RocketChat, Slack and mattermost) to automate some tasks, like triggering a C.I job or a deploy or even creation of review apps.

- Adding a prometheus metrics endpoint to fortune app or tracing with Jaeger, because it's easier to track metrics from the application, I mean code métrics like execution time, database wait time and so on.

- Adding performance tests to C.I, to find breakpoints, memory leaks and would be nice to diff between memory snapshots to check if the resource usage has increased or decreased.

- Static code analysis is a nice to have!

- Testing fresh docker build is important as well!

- Testing ansible roles using Molecule would be nice.


# What needs improvement?

- The deployer script, is required to have an existent deployment in order to the script be successful, also it doesn't clean broken deployments.

- Ansible shouldn't be deploying instances on AWS, Terraform is better to do this

- A default AMI with everything I need to run all sorts of apps, in order to have a faster instance creation.

- Better Kubernetes roles

- Better secret management, using Hashicorp's Vault is nice!

- Create more tests scenarios for both apps.

- Create a rate limit for saving messages and retrieving them

- 

