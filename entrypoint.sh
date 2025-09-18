#!/bin/sh
set -e

mkdir -p ./keys

echo "$PRIVATE_PEM" > ./keys/private.pem
echo "$PUBLIC_PEM" > ./keys/public.pem
chmod 600 ./keys/*.pem

exec ./main
