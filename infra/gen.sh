#!/usr/bin/env bash

openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:secp384r1 -days 3650 \
  -nodes -keyout key.pem -out crt.pem -subj "/CN=dev.local" \
  -addext "subjectAltName=DNS:*.dev.local,DNS:*.idp.dev.local"

echo "Copy the following values in gateway-tls-cert-secret.yaml:"
echo "tls.crt: $(cat crt.pem | base64)"
echo "tls.key: $(cat key.pem | base64)"