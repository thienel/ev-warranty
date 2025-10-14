# API ENDPOINTS

## BACKEND-GO

## Authentication & Authorization

| Method  | Endpoint                 | Description                   |
|---------|--------------------------|-------------------------------|
| POST    | `/auth/login`            | User login                    |
| POST    | `/auth/logout`           | User logout                   |
| POST    | `/auth/token`            | Refresh access token          |
| GET     | `/auth/token`            | Validate access token         |
| GET     | `/auth/google`           | Login with Google             |
| GET     | `/auth/google/callback`  | Handle callback from Google   |

## User

| Method  | Endpoint      | Description             |
|---------|---------------|-------------------------|
| POST    | `/users`      | Create a new user       |
| GET     | `/users`      | Get a list of users     |
| GET     | `/users/:id`  | Get user details        |
| PUT     | `/users/:id`  | Update user information |
| DELETE  | `/users/:id`  | Delete a user           |

## Office

| Method  | Endpoint       | Description               |
|---------|----------------|---------------------------|
| POST    | `/offices`     | Create a new office       |
| GET     | `/offices`     | Get a list of offices     |
| GET     | `/offices/:id` | Get office details        |
| PUT     | `/offices/:id` | Update office information |
| DELETE  | `/offices/:id` | Delete an office          |

## Warranty Claim

### Claims

| Method | Endpoint                  | Description                          |
|--------|---------------------------|--------------------------------------|
| GET    | `/claims`                 | Get a list of claims                 |
| POST   | `/claims`                 | Create a new claim                   |
| GET    | `/claims/:id`             | Get claim details                    |
| PUT    | `/claims/:id`             | Update claim                         |
| DELETE | `/claims/:id`             | Delete claim                         |
| POST   | `/claims/:id/submit`      | Submit claim                         |
| POST   | `/claims/:id/review`      | Review claim                         |
| POST   | `/claims/:id/requestinfo` | Request more info                    |
| POST   | `/claims/:id/cancel`      | Cancel claim                         |
| POST   | `/claims/:id/complete`    | Approve / Reject / Partially Approve |

### History

| Method | Endpoint              | Description        |
|--------|-----------------------|--------------------|
| GET    | `/claims/:id/history` | Get claim history  |

### Claim Items

| Method | Endpoint                             | Description                 |
|--------|--------------------------------------|-----------------------------|
| GET    | `/claims/:id/items`                  | Get all items of a claim    |
| POST   | `/claims/:id/items`                  | Add a new item to a claim   |
| PUT    | `/claims/:id/items/:itemID`          | Update claim item details   |
| DELETE | `/claims/:id/items/:itemID`          | Delete an item from a claim |
| POST   | `/claims/:id/items/:itemID/approve`  | Approve item                |
| POST   | `/claims/:id/items/:itemID/reject`   | Reject item                 |

### Attachments

| Method   | Endpoint                                | Description               |
|----------|-----------------------------------------|---------------------------|
| GET      | `/claims/:id/attachments`               | Get a list of attachments |
| POST     | `/claims/:id/attachments`               | Add an attachment         |
| GET      | `/claims/:id/attachments/:attachmentID` | Get attachment by ID      |
| DELETE   | `/claims/:id/attachments/:attachmentID` | Delete an attachment      |


