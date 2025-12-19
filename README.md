# Chirpy
Chirpy is a social network similar to Twitter. Users can create an account, login, and start posting "chirps". There are several API endpoints exposed related to users and chirps, as well as some related to metrics for site admins.

## [users](/chirpy/docs/users.md)
There are several endpoints with user-related functionality, including creating and updating user accounts and credentials, logging in and assigning or revoking refresh tokens to the user.

## [chirps](/chirpy/docs/chirps.md)
There is only one endpoint related to chirps: `/api/chirps` - however, this endpoint has several functions depending on the request queries. 

## [admin](/chirpy/docs/admin.md)
The endpoints for `/admin` are for checking and resetting site metrics.