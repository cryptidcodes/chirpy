# Chirps
There are several api endpoints related to chirps including those for posting and deleting chirps as well as getting a list of chirps or a specific chirp.

## /api/chirps

#### POST

A POST request sent to this endpoint will create a new chirp associated with the user who made the POST request. It requires a Body and UserID in the request structured like this:

    Body   string
	UserID uuid.UUID

A successful request will store the chirp data in the database and return a response with this structure:

    ID          UUID
    CreatedAt   Time
    UpdatedAt   Time
    Body        string
    UserID:     UUID

#### GET

A GET request sent to this endpoint will attempt to retrieve a list of the most recent chirps.

Additionally, appending a ChirpID query parameter to the end of the endpoint will attempt to GET a single chirp. The endpoint then will be `/api/chirps/{chirpID}`.

Depending if there was a specified chirp to  get or not, the endpoint will return either a single chirp or list of chirps with this structure:

    ID:        UUID
    CreatedAt: Time
    UpdatedAt: Time
    Body:      string
    UserID:    UUID

#### DELETE

A DELETE request sent to this endpoint will authenticate and check if the user is authorized to make a DELETE request, and if so, will delete the chirp from the database.