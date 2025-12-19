# Admin

There are a few endpoints for admins to use to get site data. These include:

## /admin/metrics

This endpoint is for accessing metrics related to site usage. Currently it only counts total hits to the app but could be expanded to include other data. Use a GET http request to get the metrics data.

## /admin/reset

This is an internal testing tool designed to test the metric tracking logic. Sending a POST http request to this endpoint in DEV mode will reset the metrics.