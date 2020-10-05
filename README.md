# Geocheck service

This microservice exposes a RESTful API matching a [Swagger](http://swagger.io) definition.

## Building locally

To get started building this web application locally, you can run the application natively. 

Assuming Go is installed properly on your machine - you can use the run-dev.sh script to run the server locally.
The configuration has defaults which are suited to run the application in this mode.
The build-linux.sh script will build a Linux-compatible executable in the /bin directory.
 
The following is a list of the environment variables used by this service:

- HUMANREADABLE - if set turns on human readable flag and dev mode settings for the logs. Defaults to true.
- DB_DIR - the directory that contains the maxmind location database. Defaults to 'data/'
- DB_NAME - the base name of the maxmind location database file. Defaults to 'GeoLite2-Country.mmdb'
- PORT - the port that the server should listen on. Defaults to '8080'.

The endpoint for this service is: `/geocheck` - See the Swagger endpoint below for the details on the url parameters
required to invoke this service.

In addition to the `/geocheck' endpoint, this service comes with the following capabilities:
- [Swagger UI](http://swagger.io/swagger-ui/) running on: `/explorer`
- An OpenAPI 2.0 definition hosted on: `/swagger/api`
- A Healthcheck: `/health`

## Current limitations / future changes

### Dynamic Update of Maxmind file

The current implementation requires a repackaging / re-deployment of the application in order to update the MaxMind 
database. Because you cannot simply replace the file within a running application, functionality must be developed to allow
a new database file to be downloaded so that the application can switch to that newer version. These files are typically
updated on a weekly cadence - so a once-a-day update check would be more than adequate to ensure the file is kept up to date.

The simplest approach would have each running container be responsible for periodically pulling down an updated copy and 
then have the application dynamically switch to the new version. Once switched - the old version could be deleted. The
disadvantage of this approach, beyond the obvious overhead of each container downloading its own copy,
is that you would need to synchronize the update across the containers
to avoid inconsistent results from being returned. It also opens up more potential for failures since each container
is performing this update process.

While more complex, it may be preferable for the containers to use shared storage - such as an NFS mounted filesystem, 
for storing and updating these files. A single update service could be responsible for downloading new versions of the
files on a periodic basis. Various mechanisms could be used to notify the service instances using the file of a new update. And 
likewise the services could register/deregister their usage of a given file - enabling the central update process 
to safely clean up the old files no longer in use.

### Specific overrides

For testing and for dealing with situations where the MaxMind file may be generating incorrect results for a given IP or range of IP's,
the service could support the maintenance and use of specific IP overrides - associating either a specific IP or IP network
to a specific country. 

## License

This sample application is licensed under the Apache License, Version 2. Separate third-party code objects invoked within this code pattern are licensed by their respective providers pursuant to their own separate licenses. Contributions are subject to the [Developer Certificate of Origin, Version 1.1](https://developercertificate.org/) and the [Apache License, Version 2](https://www.apache.org/licenses/LICENSE-2.0.txt).

[Apache License FAQ](https://www.apache.org/foundation/license-faq.html#WhatDoesItMEAN)
