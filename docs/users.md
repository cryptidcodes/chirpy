# Users
There are several user-related endpoints exposed by the api. Their functionality includes creating and updating user credentials, logging the user in, and assigning and revoking the user's refresh token.


## /api/users 

#### POST

Sending a POST http request to this endpoint will create a user entry in the database, hash their password and store it with the user email, and create and assign a refresh token as well as an access token. This action requires a body with this structure:

    email    string
	password string

And will return a response with this structure:

    User:
        id              uuid.UUID
        created_at      time.Time
        updated_at      time.Time
        email           string 
        is_chirpy_red   bool
    token           string
    refresh_token   string

#### PUT

Sending a PUT http request to this endpoint will update the user's credentials. It will require a similar body to the POST request but additionally require an access token to authenticate the user. Example:

    email       string
    password    string
    token       string

## /api/login

Sending a POST http request to this endpoint will log the user in and assign them a refresh token that can be used in future requests to authenticate the user. This endpoint requires the user email and password that matches the database entry for the user. Example request body:

    email string
    password string

## /api/refresh

Sending a POST request to this endpoint requires a refresh token to be present in the headers, in the `Authorization: Bearer <token>` format. A successful request will look up the token in the database. If it doesn't exist, or if it's expired, it will respond with a 401 status code. Otherwise, the response will be a 200 code and this shape:

    token UUID 

## /api/revoke

Sending a POST request to this endpoint requires a refresh token in the headers, in the `Authorization: Bearer <token>` format. 

A successful request will revoke the token in the database that matches the token that was passed in the header of the request.