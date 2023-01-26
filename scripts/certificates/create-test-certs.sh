# Inspired by https://github.com/lorenzodonini/ocpp-go/blob/master/example/1.6/create-test-certificates.sh
mkdir -p ../certs/cp
cd certs
# Create CA
echo "Generating CA..."
openssl req -new -x509 -nodes -days 120 -extensions v3_ca -keyout ca.key -out ca.crt -subj "/CN=charge-point"

# Generate self-signed charge-point certificate
echo "Generating cp private key.."
openssl genrsa -out ../cp/charge-point.key 4096

echo "Creating sign request.."
openssl req -new -out ../cp/charge-point.csr -key ../cp/charge-point.key -config $1/openssl-cp.conf

echo "Creating the certificate"
openssl x509 -req -in ../cp/charge-point.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out ../cp/charge-point.crt -days 120 -extensions req_ext -extfile $1/openssl-cp.conf
