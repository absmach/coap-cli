# Copyright (c) Abstract Machines
# SPDX-License-Identifier: Apache-2.0

CRT_LOCATION = certs
O = absmach
OU_CA = absmach_ca
OU_CRT = absmach_crt
EA = info@absmach.com
CLIENT_CN = coap-cli-client
CRT_FILE_NAME = client
CN_CA = Absmach_Self_Signed_CA

define CERT_CONFIG
[ca]
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid:always,issuer:always

[req]
req_extensions = v3_req
distinguished_name = dn
prompt = no

[dn]
CN = $(COMMON_NAME)
C  = FR
ST = FR
L  = PARIS
O  = absmach
OU = absmach

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = $(COMMON_NAME)
endef


.PHONY: all clean_certs client_cert client_cert_revoked client_cert_unknown ca

# Generate all client certificates
all: clean_certs ca client_cert client_cert_revoked client_cert_unknown

ca:
	openssl req -newkey rsa:2048 -x509 -nodes -sha512 -days 1095 \
				-keyout $(CRT_LOCATION)/ca.key -out $(CRT_LOCATION)/ca.crt -subj "/CN=$(CN_CA)/O=$(O)/OU=$(OU_CA)/emailAddress=$(EA)"

client_cert:
	# Create absmach client key and CSR.
	openssl req -new -sha256 -newkey rsa:4096 -nodes -keyout $(CRT_LOCATION)/$(CRT_FILE_NAME).key \
				-out $(CRT_LOCATION)/$(CRT_FILE_NAME).csr -subj "/CN=$(CLIENT_CN)/O=$(O)/OU=$(OU_CRT)/emailAddress=$(EA)"

	# Sign client CSR.
	openssl x509 -req -days 730 -in $(CRT_LOCATION)/$(CRT_FILE_NAME).csr -CA $(CRT_LOCATION)/ca.crt -CAkey $(CRT_LOCATION)/ca.key -CAcreateserial -out $(CRT_LOCATION)/$(CRT_FILE_NAME).crt

	# Remove CSR.
	rm $(CRT_LOCATION)/$(CRT_FILE_NAME).csr

clean_certs:
	rm -f $(CRT_LOCATION)/*.crt
	rm -f $(CRT_LOCATION)/*.key
	rm -f $(CRT_LOCATION)/*.srl
