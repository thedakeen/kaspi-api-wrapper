#!/bin/bash

set -e

mkdir -p ./certs

COUNTRY="KZ"
STATE="Astana"
LOCALITY="Astana"
ORGANIZATION="Test Kaspi Pay"
CA_CN="Test Kaspi CA"
SERVER_CN="mtokentest.kaspi.kz"
CLIENT_CN="test-client"

echo "=== Генерация корневого CA сертификата ==="
openssl genrsa -out ./certs/ca.key 2048
openssl req -x509 -new -nodes -key ./certs/ca.key -sha256 -days 3650 -out ./certs/ca.crt \
    -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORGANIZATION/CN=$CA_CN"

echo "=== Генерация серверного сертификата ==="
openssl genrsa -out ./certs/server.key 2048
openssl req -new -key ./certs/server.key -out ./certs/server.csr \
    -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORGANIZATION/CN=$SERVER_CN"

cat > ./certs/server.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = mtokentest.kaspi.kz
DNS.2 = mock-kaspi-api-standard
DNS.3 = mock-kaspi-api-enhanced
DNS.4 = localhost
IP.1 = 127.0.0.1
EOF

openssl x509 -req -in ./certs/server.csr -CA ./certs/ca.crt -CAkey ./certs/ca.key \
    -CAcreateserial -out ./certs/server.crt -days 825 -sha256 \
    -extfile ./certs/server.ext

echo "=== Генерация клиентского сертификата ==="
openssl genrsa -out ./certs/client.key 2048
openssl req -new -key ./certs/client.key -out ./certs/client.csr \
    -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORGANIZATION/CN=$CLIENT_CN"

cat > ./certs/client.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = clientAuth
EOF

openssl x509 -req -in ./certs/client.csr -CA ./certs/ca.crt -CAkey ./certs/ca.key \
    -CAcreateserial -out ./certs/client.crt -days 825 -sha256 \
    -extfile ./certs/client.ext

echo "=== Создание PFX-файла для клиентского сертификата ==="
openssl pkcs12 -export -out ./certs/client.pfx -inkey ./certs/client.key -in ./certs/client.crt \
    -certfile ./certs/ca.crt -passout pass:test123

# Для совместимости создаем копию CA как client-ca.crt
cp ./certs/ca.crt ./certs/client-ca.crt

echo ""
echo "=== Сертификаты успешно созданы ==="
echo ""
echo "Для использования в вашем приложении, обновите файл .env:"
echo ""
echo "KASPI_CERT_FILE=./certs/client.crt"
echo "KASPI_KEY_FILE=./certs/client.key"
echo "KASPI_KEY_PASSWORD="
echo "KASPI_PFX_FILE=./certs/client.pfx"
echo "KASPI_ROOT_CA_FILE=./certs/ca.crt"
echo ""
echo "Пароль для PFX-файла: test123"

rm -f ./certs/*.csr ./certs/*.ext ./certs/*.srl

echo ""
echo "Проверка цепочки сертификатов..."
openssl verify -CAfile ./certs/ca.crt ./certs/client.crt
openssl verify -CAfile ./certs/ca.crt ./certs/server.crt

echo ""
echo "Проверка завершена. Все сертификаты должны быть корректно подписаны одним CA."