import sys
import time
import os
import grpc
import decimal
import ecommerce_pb2
import ecommerce_pb2_grpc
from concurrent import futures


class Ecommerce(ecommerce_pb2_grpc.DiscountServicer):
    def ApplyDiscount(self, request, content):
        customer = request.customer
        product = request.product
        discount = ecommerce_pb2.DiscountValue()
        if customer.id == 1 and product.price_in_cents > 0:
            percentual = decimal.Decimal(10) / 100  # save 10%
            price = decimal.Decimal(product.price_in_cents) / 100
            new_price = price - (price * percentual)
            value_in_cents = int(new_price * 100)
            discount = ecommerce_pb2.DiscountValue(pct=percentual, value_in_cents=value_in_cents)

        new_product = ecommerce_pb2.Product(id=product.id,
            slug=product.slug,
            description=product.description,
            price_in_cents=product.price_in_cents,
            discount_value=discount)
        return ecommerce_pb2.DiscountResponse(product=new_product)


def get_server(host):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    keys_dir = os.path.abspath(os.path.join('.', os.pardir, 'keys'))
    with open('%s/private.key' % keys_dir, 'rb') as f:
        private_key = f.read()
    with open('%s/cert.pem' % keys_dir, 'rb') as f:
        certificate_chain = f.read()
    server_credentials = grpc.ssl_server_credentials(((private_key, certificate_chain),))
    server.add_secure_port(host, server_credentials)
    ecommerce_pb2_grpc.add_DiscountServicer_to_server(Ecommerce(), server)
    return server


if __name__ == '__main__':
    port = sys.argv[1] if len(sys.argv) > 1 else 443
    host = '[::]:%s' % port
    server = get_server(host)
    try:
        server.start()
        print('Running Discount service on %s' % host)
        while True:
            time.sleep(1)
    except Exception as e:
        print('[error] %s' % e)
        server.stop(0)

