# UAC - User access controller  

This app syncing users permissions between corporate AD server and OCP cluster 

## Generate TLS certificates 
You'll need to provide `ca.crt`, `server.crt` and `server.key` base64 encoded strings while deploying the UAC admission controller web hook. 
Run following commands to generate all required certificates.  
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

## UAC Admission control webhook configuration
The webhook can be configured by `config.json` file or by environment variables.
Each entry in `config.json` can be override with environment variable. 

Example `config.json` 
```json
{
  "http": {
    "crt": "/path/to/server.crt",
    "key": "/path/to/server.key"
  },
  "ad": {
    "host": "10.2.3.4",
    "port": 389,
    "baseDN": "dc=ad,dc=lab",
    "bindUser": "admin1",
    "bindPass": "admin1"
  }
}
```
Environment variables example
```bash
UAC_HTTP_CRT="base64 encoded cert"
UAC_HTTP_KEY="base64 encoded key"
UAC_AD_HOST="10.1.2.3"
UAC_AD_PORT=389
UAC_AD_BASEDN="dc=ad,dc=lab"
UAC_AD_BINDUSER="admin1"
UAC_AD_BINDPASS="admin1"
```   

