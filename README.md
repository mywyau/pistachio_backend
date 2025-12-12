# pistachio_backend

### Get dependencies

e.g. For uuid

```
 go get <github.com/google/uuid>
```

### Access Database Locally

```
psql -h localhost -p 5432 -U pistachio -d pistachio_db
```

### Run goose sql migrations

```
goose postgres "postgres://pistachio:pistachio_pwd@localhost:5432/pistachio_db" up
```

### frontend json request body

{
"customer": {
"name": "John Smith",
"email": "john@example.com",
"phone": "123456",
"address": "12 Hill Road"
},
"title": "Leaking tap",
"description": "Kitchen sink leak",
"estimate": 45.0
}

### get a list of all jobs

```
curl http://localhost:8080/jobs/b0d92fb3-fa63-4663-b790-1146b3948e7f
```

### get a jobs notes left by tradesmen person

```
curl -X POST localhost:8080/jobs/b0d92fb3-fa63-4663-b790-1146b3948e7f/notes \
  -H "Content-Type: application/json" \
  -d '{
        "text": "Checked the water pressure, ordered replacement parts"
      }'
```

### insert photos

```
curl -X POST \
  -F "file=@/Users/michaelyau/Desktop/gnome.png" \
  http://localhost:8080/jobs/b0d92fb3-fa63-4663-b790-1146b3948e7f/photos
```

### Create a job

```
curl -X POST localhost:8080/jobs/b0d92fb3-fa63-4663-b790-1146b3948e7f/invoice \
  -H "Content-Type: application/json" \
  -d '{"amount": 120.00}'

```

### Update jobs status

```
curl -X PUT localhost:8080/jobs/b0d92fb3-fa63-4663-b790-1146b3948e7f/status \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

### Create an invoice

```
curl -X POST http://localhost:8080/invoices \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John Smith",
    "customer_email": "john.smith@example.com",
    "customer_address": {
      "line1": "12 High Street",
      "line2": "Flat 3B",
      "city": "London",
      "postcode": "SW1A 1AA",
      "country": "UK"
    },
    "items": [
      {
        "description": "Fix leaking tap",
        "quantity": 2,
        "unitPrice": 45.00
      },
      {
        "description": "Replace pipe section",
        "quantity": 1,
        "unitPrice": 120.00
      }
    ]
  }'
```
