package main

const USAGE = `
# Usage - tfdoc

## Get Providers list

- Method: GET

- URL: /provider/

- Example:

	curl -v http://127.0.0.1:8080/provider

## Get Resources list

- Method: GET

- URL: /provider/<provider>

- Example:

	curl -v http://127.0.0.1:8080/provider/tencentcloud

## Get Arguments list

- Method: GET

- URL: /provider/<provider>/<resource>

- Example:

	curl -v http://127.0.0.1:8080/provider/tencentcloud/tencentcloud_security_group
`
