It is a proof of concept based in a fictional e-commerce platform. There are two microservices:

1) Catalog: It is written in Go and exposes a REST API that returns a product list.
2) Discount: It is consumed by Catalog, written in Python and returns the received product with 10% discount applied.

The full text about it can be found in my blog https://gustavohenrique.com.

## Getting Started

```
# clone this repo
cd $HOME
git clone https://github.com/gustavohenrique/microservices-grpc-go-python.git
cd microservices-grpc-go-python

# rename the keys to work with Docker
rm keys/cert.pem keys/private.key
mv keys/discount.pem keys/cert.pem
mv keys/discount.key keys/private.key

# run docker compose
docker-compose up -d
```

The ports used are 11443 and 11080. 
To test using curl:

```
curl -H 'X-USER-ID: 1' http://localhost:11080/products
```
