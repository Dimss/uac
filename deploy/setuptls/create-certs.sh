#!/usr/bin/env bash
init () {
#    if [ -d "$WORK_DIR" ]; then
#        rm -fr ${WORK_DIR}
#    fi
#    mkdir ${WORK_DIR}
    cd ${WORK_DIR}
    cat << EOF > conf
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, serverAuth
EOF
}

create_ca () {
    openssl genrsa -out ca.key 2048
    openssl req -x509 -new -nodes -key ca.key -days ${EXPIRATION_DAYS} -out ca.crt -subj "/CN=admission_ca"
}

create_server_crts () {
    openssl genrsa -out server.key 2048
    openssl req -new -key server.key -out server.csr -subj "/CN=${COMMON_NAME}" -config conf
    openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days ${EXPIRATION_DAYS} -extensions v3_req -extfile conf
}

print_base64_certs (){
    echo -e "base64 encoded ca.crt\n"
    base64 -i /tmp/webhook_deployment/ca.crt
    echo -e "\n"
    echo -e "base64 encoded server.crt\n"
    base64 -i /tmp/webhook_deployment/server.crt
    echo -e "\n"
    echo -e "base64 encoded server.key\n"
    base64 -i /tmp/webhook_deployment/server.key
    echo -e "\n"
}
if [ "$#" -ne 1 ]; then
    echo "Missing certificate common name (CN). Example usage: ./create-certs.sh uac.bnhp-system.svc.cluster.local"
    exit 1
fi

WORK_DIR="/tmp/webhook_deployment"
COMMON_NAME=$1
EXPIRATION_DAYS=36500

echo ${WORK_DIR}
echo ${COMMON_NAME}
init
create_ca
create_server_crts
print_base64_certs
