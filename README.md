# iFortepay-take-home

You can access the swagger on localhost:8080/docs/index.html

## Seeded data
```
INSERT INTO products (id, sku, name, price, quantity) VALUES
(1, '120P90', 'Google Home', 49.99, 10),
(2, '43N23P', 'MacBook Pro', 5399.99, 5),
(3, 'A304SD', 'Alexa Speaker', 109.50, 10),
(4, '234234', 'Raspberry Pi B', 30.00, 2);
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