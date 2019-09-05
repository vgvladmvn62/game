# Bullseye UI

This is a UI for [Bullseye](https://github.wdf.sap.corp/Magikarpie/bullseye) made with React. 
For info on developing with react read the wiki.

To launch you need to start the [backend](https://github.wdf.sap.corp/Magikarpie/bullseye).

## Debugging

There are few debugging tools that will help you in your jurney.
You should seriously consider using them. 

- [redux dev tools for Chrome](http:github.com/zalmoxisus/redux-devtools-extension) [with src and instructions](http:github.com/zalmoxisus/redux-devtools-extension)
- [react dev tools for Chrome](https://chrome.google.com/webstore/detail/react-developer-tools/fmkadmapgofadopljbjfkapdkoienihi?hl=en) which work without additional setup

## Running as docker image

1) docker build -t bullseye-ui:latest .
2) docker run -d -p 80:5000 bullseye-ui:latest

Now it will be running on http://localhost

## Configuring backend URL

Due to React reading environmental variables at build time rather than start time, we had to store backend URL in file in
order to run on Kubernetes.

You can set backend URL in **public/config/config.js**. Default value is http://localhost:8080.

## Kubernetes


### (Optional) build Docker image and push it to artifactory repo

Make sure that you are in factory-dockerysf-deploy group, which has deploy permissions to the docker-kyma repository.

Then login into repository:
```shell
docker login repository.hybris.com:5011 
```

There is already bullseye Docker image in artifactory, but if you make some changes you can build and push new image by running:

```shell
./build.sh
```

**Warning:** It will be tagged as bullseye-ui:latest and will override previous image. Use with caution.

### Deployment on Kubernetes cluster

First make sure that you have kubectl. Then you need to point KUBECONFIG to your Kubernetes cluster. This
topic is not in scope of this guide. We assume that if you want to deploy on Kubernetes you have the necessary knowledge.

Update backend URL in **k8s/configmap.yaml**.

Then you need to apply yamls:

```shell
kubectl apply -f k8s/regcred.yaml -n bullseye
kubectl apply -f k8s/configmap.yaml -n bullseye
kubectl apply -f k8s/deployment.yml -n bullseye
kubectl apply -f k8s/service.yaml -n bullseye
```

### (Optional next step) Deployment on Kubernetes cluster with ISTIO (like Kyma)

You need to change *hosts* and *http* sections in **k8s/virtualservice.yaml** to your cluster host and apply it:

```shell
kubectl apply -f k8s/virtualservice.yaml -n bullseye
```