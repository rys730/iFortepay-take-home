# iFortepay-take-home

## Running
To run the service, please use the docker compose provided. Adjust the environment values accordingly. The database migration will run along with the docker compose using `goose`.

You can access the swagger on localhost:8080/docs/index.html

## Seeded data
```
INSERT INTO products (id, sku, name, price, quantity) VALUES
(1, '120P90', 'Google Home', 49.99, 10),
(2, '43N23P', 'MacBook Pro', 5399.99, 5),
(3, 'A304SD', 'Alexa Speaker', 109.50, 10),
(4, '234234', 'Raspberry Pi B', 30.00, 2);

INSERT INTO promotions (id, promotion_type, start_date, end_date) VALUES
(1, 'FREE_ITEM', '2024-01-01', '2026-12-31'),
(2, 'BUY_X_PAY_Y', '2024-01-01', '2026-12-31'),
(3, 'BULK_DISCOUNT', '2024-01-01', '2026-12-31');


INSERT INTO product_promotions (promotion_id, product_id, min_quantity, free_product_id, free_quantity) VALUES
(1, 2, 1, 4, 1);

INSERT INTO product_promotions (promotion_id, product_id, min_quantity, pay_y) VALUES
(2, 1, 3, 2);

INSERT INTO product_promotions (promotion_id, product_id, min_quantity, discount) VALUES
(3, 3, 3, 0.10);
```

## cURL example
Adjust the items id based on the seeded data id
```
curl -X 'POST' \
  'http://localhost:8080/api/checkout' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "items": [
    {
      "id": 2,
      "quantity": 1
    }
  ]
}'
```

## Architecture guide

By separating concerns into folders, making it cleaner:
- `/cmd` -> application starting point (services, jobs, etc). application will be initialized here

- `/db` -> migrations and queries

- `/docs` -> swagger docs

- `/infrastructure` -> connection to external actor

- `/internal` -> all services logic lives here
    - `/common` -> configs, utils, etc
    - `/handler` -> endpoint starts here. used for request handling and validation
    - `/model` -> entities, value objects, data transfer objects, etc
    - `/usecase` -> business logic
    - `/repository` -> DB, HTTP/External APIs

