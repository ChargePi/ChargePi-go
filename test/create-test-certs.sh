#Copied from https://github.com/lorenzodonini/ocpp-go/blob/master/example/1.6/create-test-certificates.sh
mkdir -p certs/cp
cd certs
# Create CA
openssl req -new -x509 -nodes -days 120 -extensions v3_ca -keyout ca.key -out ca.crt -subj "/CN=ocpp-go-example"
# Generate self-signed charge-point certificate
openssl genrsa -out cp/charge-point.key 4096
openssl req -new -out cp/charge-point.csr -key cp/charge-point.key -config ../openssl-cp.conf
openssl x509 -req -in cp/charge-point.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out cp/charge-point.crt -days 120 -extensions req_ext -extfile ../openssl-cp.conf