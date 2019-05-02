# UAC - User access controller  

This app syncing users permissions between corporate AD server and OCP cluster 

## Generate TLS certificates 
You'll need to provide `ca.crt`, `server.crt` and `server.key` base64 encoded strings while deploying the UAC admission controller web hook. 
Run following command to generate all required certificates.  
 - `UAC_WEBHOOK_SERVICE_NAME=uac.bnhp-system.svc.cluster.local make generate-tls`
 - `cd /tmp/webhook_deployment` 
 - Get `ca.crt` by executing `base64 -i /tmp/webhook_deployment/ca.crt`
 - Get `server.crt` by executing `base64 -i /tmp/webhook_deployment/server.crt`
 - Get `server.key` by executing `base64 -i /tmp/webhook_deployment/server.key`

## Local build with S2I Build
 - Install Golang and S2I on you machine 
 - Run local build: `DOCKER_REGISTRY=your.docker.registry/repo/uac make build-docker`

## Build UAC image inside OpenShift  
 - Add S2I golang build image `oc create -f deploy/ocp/golang-is.yaml` 
 - Import OpenShift template either from UI or cli `oc create -f deploy/ocp/build-template.yaml`
 - Start new build either from UI or cli `oc start-build uac-build --follow`

## Deploy UAC Admission control webhook 
 - Create UAC webhook by creating OpenShift template either by UI or CLI `oc create -f deploy/ocp/uac-template.yaml`
 - Deploy UAC webbook either from CLI or UI `oc process -f deploy/ocp/uac-template.yaml | oc create -f -`
