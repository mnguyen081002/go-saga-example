# Golang Saga Pattern with DTM

Orchestration saga example with docker compose, golang and dtm

## What is DTM

DTM is a distributed transaction framework which provides cross-service eventual data consistency. It provides saga, tcc, xa, 2-phase message, outbox, workflow patterns for a variety of application scenarios.

Ref: https://github.com/dtm-labs/dtm

## Architecture

![Architecture Overview](./assets/SagaDiagram.png)

## Start with docker compose

```bash
docker-compose up
```

## Step to run

### Step 1: Create a new item

```curl --location --request POST 'http://localhost:8082/api/item/create' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Iphone 14",
    "price": 5000,
    "price_max": 5000,
    "price_min": 5000,
    "show_free_ship": true,
    "description": "Iphone 14",
    "sku": "123456789",
    "quantity": 1000,
    "stock":10,
    "discount": "0%",
    "raw_discount": 0,
    "images": ["123456789"],
    "category_id": "2",
    "variant_ids": [1]
}'
```

### Step 2: Create a new user wallet

```curl --location --request POST 'localhost:8081/api/user-wallet' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":"1",
    "balance":999999
}'
```

### Step 3: Create a new order saga

```curl --location --request POST 'localhost:8080/api/order/create-order-saga' \
--header 'Content-Type: application/json' \
--data-raw '{
    "order_items": [
        {
            "item_id": "0000022b-d387-4e7d-b0fe-506c75fc9158",
            "quantity": 1
        }
    ]
}'
```
