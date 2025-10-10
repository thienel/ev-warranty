# API ENDPOINTS

## BACKEND-GO

## Authentication & Authorization

| Method | Endpoint                | Description                 |
| ------ | ----------------------- | --------------------------- |
| POST   | `/auth/login`           | User login                  |
| POST   | `/auth/logout`          | User logout                 |
| POST   | `/auth/token`           | Refresh access token        |
| GET    | `/auth/token`           | Validate access token       |
| GET    | `/auth/google`          | Login with Google           |
| GET    | `/auth/google/callback` | Handle callback from Google |

## User

| Method | Endpoint     | Description             |
| ------ | ------------ | ----------------------- |
| POST   | `/users`     | Create a new user       |
| GET    | `/users`     | Get a list of users     |
| GET    | `/users/:id` | Get user details        |
| PUT    | `/users/:id` | Update user information |
| DELETE | `/users/:id` | Delete a user           |

## Office

| Method | Endpoint       | Description               |
| ------ | -------------- | ------------------------- |
| POST   | `/offices`     | Create a new office       |
| GET    | `/offices`     | Get a list of offices     |
| GET    | `/offices/:id` | Get office details        |
| PUT    | `/offices/:id` | Update office information |
| DELETE | `/offices/:id` | Delete an office          |

## Warranty Claim

### Claims

| Method | Endpoint             | Description          |
| ------ |----------------------| -------------------- |
| GET    | `/claims`            | Get a list of claims |
| GET    | `/claims/:id`        | Get claim details    |
| POST   | `/claims`            | Create a new claim   |
| PUT    | `/claims/:id`        | Update claim         |
| DELETE | `/claims/:id`        | Delete claim         |
| PATCH  | `/claims/:id/status` | Update claim status  |

### Claim Items

| Method | Endpoint                             | Description                 |
| ------ |--------------------------------------| --------------------------- |
| GET    | `/claims/:id/items`                  | Get all items of a claim    |
| POST   | `/claims/:id/items`                  | Add a new item to a claim   |
| PUT    | `/claims/:id/items/{item_id}`        | Update claim item details   |
| DELETE | `/claims/:id/items/{item_id}`        | Delete an item from a claim |
| PATCH  | `/claims/:id/items/{item_id}/status` | Update item status          |

### Attachments

| Method | Endpoint                                   | Description               |
| ------ | ------------------------------------------ | ------------------------- |
| GET    | `/claims/:id/attachments`                 | Get a list of attachments |
| POST   | `/claims/:id/attachments`                 | Add an attachment         |
| DELETE | `/claims/:id/attachments/{attachment_id}` | Delete an attachment      |


### Allocations

| Method | Endpoint                                           | Description               |
| ------ | -------------------------------------------------- | ------------------------- |
| GET    | `/claims/:id/allocations`                         | Get a list of allocations |
| POST   | `/claims/:id/allocations`                         | Create an allocation      |
| PATCH  | `/claims/:id/allocations/{allocation_id}/confirm` | Confirm an allocation     |

### History

| Method | Endpoint                | Description          |
| ------ | ----------------------- | -------------------- |
| GET    | `/claims/:id/history`  | Get claim history    |

