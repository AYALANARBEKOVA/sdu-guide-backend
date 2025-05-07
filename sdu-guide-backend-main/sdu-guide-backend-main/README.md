# SDU Guide Backend

## Описание
`SDU Guide Backend` — это серверная часть проекта, предназначенная для управления данными и взаимодействия с клиентами.

## Установка и запуск

### 1. Установите и запустите Redis

Перед запуском сервера необходимо установить и запустить Redis.

#### Установка Redis (Linux/macOS)
```sh
sudo apt update && sudo apt install redis -y  # Для Ubuntu/Debian
brew install redis  # Для macOS (Homebrew)
```

#### Запуск Redis
```sh
redis-server
```

### 2. Установите и запустите MongoDB

Проект также использует MongoDB. Убедитесь, что MongoDB установлена и запущена.

#### Установка MongoDB (Linux/macOS)
```sh
sudo apt update && sudo apt install -y mongodb  # Для Ubuntu/Debian
brew tap mongodb/brew && brew install mongodb-community@6.0  # Для macOS (Homebrew)
```


### 3. Запуск Backend-сервера

После запуска Redis и MongoDB выполните следующую команду для запуска backend-сервера:
```sh
go run ./cmd/main.go
```


## Технологии
- Golang
- Redis
- MongoDB
- Gin 




## User Endpoints

### Sign Up
**Endpoint:** `POST /sign-up`
**Description:** Registers a new user.
**Request Body:**
```json
{
  "login": "string",
  "email": "string",
  "password": "string"
}
```
**Response:**
```json
{
  "Status": "User successfully created"
}
```

### Sign In
**Endpoint:** `POST /sign-in`
**Description:** Authenticates a user.
**Request Body:**
```json
{
  "email": "string",
  "password": "string"
}
```
**Response:**
```json
{
  "Status": "User successfully signed in"
}
```

### Logout
**Endpoint:** `GET /logout`
**Description:** Logs out the user.
**Response:**
```json
{
  "Status": "User successfully logged out from the system"
}
```

### Get Profile
**Endpoint:** `GET /profile`
**Description:** Retrieves the user profile.
**Response:**
```json
{
  "data": {
    "_id": "number",
    "username": "string",
    "email": "string",
    "firstName": "string",
    "lastName": "string",
    "registrationDate": "string (ISO 8601 format)",
    "lastLogin": "string (ISO 8601 format)",
    "imageHash": "string"
  }
}

```

### Update User
**Endpoint:** `PUT /update-user`
**Description:** Updates user information.
**Request Body:**
```json
{
  "username": "string",
  "email": "string",
  "firstName": "string",
  "lastName": "string",
  "age": "number",
  "gender": "string",
  "imageHash": "string"
}
```
**Response:**
```json
{
  "Status": "User successfully updated"
}
```

### Get Translations
**Endpoint:** `GET /translations`
**Description:** Retrieves translations based on the `lang` query parameter.

---

## Protected Endpoints (Require Authentication)
*For the following requests, the user must be authenticated, otherwise a `401 Unauthorized` response will be returned.*


### Ping
**Endpoint:** `GET /ping`
**Description:** Checks system status.
**Response:**
```json
{
  "Status": "Check"
}
```

All endpoints inside the / protected group require authentication:

/upload-XLSX

/upload-image

/create-room

/update-room

/getAll-rooms

/update-user

/delete-room/{id}

/create-event

/update-event

/delete-event/{id}



## Room/File Endpoints

### File Upload and Retrieval

#### 1. Upload XLSX File
**Endpoint:** `POST /upload-XLSX`  
**Authorization:** Required  
**Headers:**  
- `Content-Type: multipart/form-data`

**Form Data:**  
- `xlsx`: File (XLSX format)

**Response:**
```json
{
  "hash": "generated_file_hash"
}
```

#### 2. Upload Image File
**Endpoint:** `POST /upload-image`  
**Authorization:** Required  
**Headers:**  
- `Content-Type: multipart/form-data`

**Form Data:**  
- `image`: File (Image format)

**Response:**
```json
{
  "hash": "generated_file_hash"
}
```

#### 3. Get XLSX File by Hash
**Endpoint:** `GET /xlsx/:hash`  
**Authorization:** Not required  
**Response:** Binary XLSX file

#### 4. Get Image File by Hash
**Endpoint:** `GET /image/:hash`  
**Authorization:** Not required  
**Response:** Binary Image file

### Room Management

#### 5. Create Room
**Endpoint:** `POST /create-room`  
**Authorization:** Required  
**Headers:**  
- `Content-Type: application/json`

**Body:**  
```json
{
  "id": 1,
  "block": "A",
  "number": 101,
  "sef": "room-101",
  "hash": "some_hash",
  "deleted": false,
  "updated": "2024-03-14T12:00:00Z"
}
```

**Response:**
```json
{
  "Status": "Room successfully created"
}
```

#### 6. Get Room by ID
**Endpoint:** `GET /room/:id`  
**Authorization:** Not required  

**Response:**
```json
{
  "data": {
    "id": 1,
    "block": "A",
    "number": 101,
    "sef": "room-101",
    "hash": "some_hash",
    "deleted": false,
    "updated": "2024-03-14T12:00:00Z"
  }
}
```

#### 7. Update Room
**Endpoint:** `PUT /update-room`  
**Authorization:** Required  
**Headers:**  
- `Content-Type: application/json`

**Body:**  
```json
{
  "id": 1,
  "block": "B",
  "number": 102,
  "sef": "room-102",
  "hash": "updated_hash",
  "deleted": false,
  "updated": "2024-03-14T12:30:00Z"
}
```

**Response:**
```json
{
  "Status": "Room successfully updated"
}
```

#### 8. Get All Rooms
**Endpoint:** `GET /getAll-rooms`  
**Authorization:** Required  
**Query Parameters:**  
- `block` (optional)
- `number` (optional)

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "block": "A",
      "number": 101,
      "sef": "room-101",
      "hash": "some_hash",
      "deleted": false,
      "updated": "2024-03-14T12:00:00Z"
    }
  ]
}
```

#### 9. Delete Room
**Endpoint:** `DELETE /delete-room/:id`  
**Authorization:** Required  

**Response:**
```json
{
  "Status": "Room successfully deleted"
}
```

#### 10. Get Room Schedule
**Endpoint:** `GET /schedule/:sef`  
**Authorization:** Not required  

**Response:** Binary XLSX file containing the room schedule

---

## Notes
- Ensure authentication for protected routes.
- Use proper file handling when uploading XLSX and image files.
- Query parameters in `getAll-rooms` allow filtering by `block` and `number`.
- `getSchedule` fetches the XLSX schedule using a SEF identifier.

This documentation serves as a quick reference for API interaction via Postman.



## Event Endpoints

### Create Event
**Endpoint:** `POST /create-event`
**Description:** Creates an event.
**Request Body:**
```json
{
  "name": "string",
  "date": "string (ISO 8601 format)",
  "place": "string",
  "startTime": "string",
  "endTime": "string",
  "ended": "boolean",
  "hash": "string",
  "shortName": "string"
}
```
**Response:**
```json
{
  "Status": "Event successfully created"
}
```

### Get Event
**Endpoint:** `GET /get-event/:id`
**Description:** Retrieves an event by ID.
**Response:**
```json
{
  "data": {
    "id": "number",
    "name": "string",
    "date": "string (ISO 8601 format)",
    "place": "string",
    "startTime": "string",
    "endTime": "string",
    "ended": "boolean",
    "hash": "string",
    "shortName": "string"
  }
}
```

### Update Event
**Endpoint:** `PUT /update-event`
**Description:** Updates event details.
**Request Body:**
```json
{
  "id": "number",
  "name": "string",
  "date": "string (ISO 8601 format)",
  "place": "string",
  "startTime": "string",
  "endTime": "string",
  "ended": "boolean",
  "hash": "string",
  "shortName": "string"
}
```
**Response:**
```json
{
  "Status": "Event successfully updated"
}
```

### Get All Events
**Endpoint:** `GET /getAll-events`
**Description:** Retrieves all events with optional filters.
**Response:**
```json
{
  "data": [
    {
      "id": "number",
      "name": "string",
      "date": "string (ISO 8601 format)",
      "place": "string",
      "startTime": "string",
      "endTime": "string",
      "ended": "boolean",
      "hash": "string",
      "shortName": "string"
    }
  ]
}
```

### Get Events for Calendar
**Endpoint:** `GET /getAll-events-calendar`
**Description:** Retrieves events for the current month.
**Response:**
```json
{
  "data": [
    {
      "id": "number",
      "name": "string",
      "date": "string (ISO 8601 format)",
      "place": "string",
      "startTime": "string",
      "endTime": "string",
      "ended": "boolean",
      "hash": "string",
      "shortName": "string"
    }
  ]
}
```

### Delete Event
**Endpoint:** `DELETE /delete-event/:id`
**Description:** Deletes an event by ID.
**Response:**
```json
{
  "Status": "Event successfully deleted"
}
```

**Note:** Ensure that a valid authentication token is sent in the request headers to access protected routes.

