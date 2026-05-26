FROM gomicro/goose

RUN mkdir -p /migrations/sso-service \
    && mkdir -p /migrations/order-service \
    && mkdir -p /migrations/products-service

COPY migrations/order-service/*.sql /migrations/order-service/
COPY migrations/products-service/*.sql /migrations/products-service/
COPY migrations/sso-service/*.sql /migrations/sso-service/
COPY migrations/cart-service/*.sql /migrations/cart-service/
COPY migrations/entrypoint.sh /migrations/
RUN chmod +x /migrations/entrypoint.sh
ENTRYPOINT ["/migrations/entrypoint.sh"]