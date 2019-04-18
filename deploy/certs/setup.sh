#!/usr/bin/env bash
# https://docs.okd.io/latest/architecture/additional_concepts/dynamic_admission_controllers.html
# https://github.com/stevesloka/validatingwebhook/blob/master/deployment/02-webhook.yaml

# Generate self sign certificate

## Create CA
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -days 100000 -out ca.crt -subj "/CN=admission_ca"

## Generate server
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj "/CN=172.20.10.5" -config conf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 100000 -extensions v3_req -extfile conf

#####  conf file example #####
#
#[req]
#req_extensions = v3_req
#distinguished_name = req_distinguished_name
#[req_distinguished_name]
#[ v3_req ]
#basicConstraints = CA:FALSE
#keyUsage = nonRepudiation, digitalSignature, keyEncipherment
#extendedKeyUsage = clientAuth, serverAuth
#subjectAltName = IP:172.20.10.5