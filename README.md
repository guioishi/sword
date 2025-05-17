# Task Tracker App

This is an api to account for maintenance tasks performed during a working day.
There are two types of users, the manager and the technician.
The manager could see all tasks, delete them and should be notified when tasks are created of updated.
The technician could create, update and see his own tasks.

## Run local

To run app locally, you just need up the configured docker compose environment with the command:

```sh
$ docker compose up --build
```

## Setup K8S

These are the steps to configure your K8S Environment:

- Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/).
- Update .env file with your REGISTRY url.
- Configure kubectl context with your cluster info. ``kubectl config``.
- Setup RabbitMQ dependencies with the following commands:

```sh
$ kubectl apply -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml
$ kubectl get all -o wide -n rabbitmq-system
```
- Execute the bash file ``deploy.sh``.
