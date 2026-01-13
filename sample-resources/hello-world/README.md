# Hello World Apache Application (state-less)

A simple Kubernetes application that deploys an Apache HTTPD web server serving a "Hello World" static page.

## Application Description

This application consists of two Kubernetes resources defined in `manifest.yaml`:

### Deployment
- **Name**: `apache-hello`
- **Replicas**: 1
- **Container**: Apache HTTPD (httpd:alpine)
- **Port**: 80
- **Functionality**: Serves a simple HTML page with "Hello World!" message and a note about running Apache HTTPD in Kubernetes

The deployment uses a command override to create an inline static HTML page that displays:
```
Hello World!
Running Apache HTTPD static page in k8s.
```

### Service
- **Name**: `apache-hello-service`
- **Type**: NodePort
- **Port**: 80
- **Target Port**: 80
- **Selector**: Routes traffic to pods with label `app: apache-hello`

The NodePort service exposes the Apache server externally, making it accessible from outside the cluster.

## Prerequisites

- Kubernetes cluster (local or remote)
- `kubectl` configured to communicate with your cluster
- Proper permissions to create deployments and services

## Deployment

To deploy the application to your Kubernetes cluster, run:

```bash
./deploy.sh
```

Or manually:

```bash
kubectl create -f manifest.yaml
```

This will create both the deployment and the service in your currently selected namespace. The manifest does not specify a target namespace, so resources will be deployed to whichever namespace is active in your current kubectl context.

### Verifying Deployment

Check the deployment status:
```bash
kubectl get deployments apache-hello
kubectl get pods -l app=apache-hello
kubectl get services apache-hello-service
```

### Accessing the Application

To access the application, get the NodePort assigned to the service:

```bash
kubectl get service apache-hello-service
```

Then access the application using:
```
http://<node-ip>:<node-port>
```

Or use port-forwarding for local access:
```bash
kubectl port-forward service/apache-hello-service 8080:80
```

Then visit `http://localhost:8080` in your browser.

## Cleanup

To remove the application from your cluster, run:

```bash
./destroy.sh
```

Or manually:

```bash
kubectl delete -f manifest.yaml
```

This will delete both the deployment and the service.
