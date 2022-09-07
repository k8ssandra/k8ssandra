---
title: "Develop with Stargate APIs"
linkTitle: "Develop"
weight: 5
description: Develop clients to interact with Apache Cassandra® data via Stargate REST, GraphQL and document APIs.
---

[Stargate](https://stargate.io/) is an open-source data gateway providing common API interfaces for backend databases. With K8ssandra, Stargate may be deployed in front of the Apache Cassandra® cluster providing CQL, REST, GraphQL, and document-based API endpoints. Stargate itself may be scaled horizontally within the cluster as needed. This scaling is done independently from the data layer.

This topic provides information about accessing the various API endpoints provided by Stargate.

While this document will help get you up and going quickly with Stargate, more detailed information about using Stargate can be found in the [Stargate docs](https://stargate.io/docs/latest/quickstart/quickstart.html).

## Tools

* HTTP client (cURL, Postman, etc.)
* Web Browser

## Prerequisites

1. [K8ssandra Cluster]({{< relref "quickstarts#install-k8ssandra" >}})
1. [Ingress]({{< relref "/tasks/connect/ingress" >}}) configured to expose each of the Stargate services (Auth, REST, GraphQL)
1. DNS names configured for the exposed Stargate services, relreferred to as `STARGATE_AUTH_DOMAIN`, `STARGATE_REST_DOMAIN`, and `STARGATE_GRAPHQL_DOMAIN` below.

## Access Auth API

Before accessing any of the provided Stargate data APIs, an auth token must be generated and provided to subsequent data API requests.  Use the auth API to generate a token.

The default port exposed by Stargate for the auth API is `8081`, these examples will assume that is the port exposed by the cluster ingress configuration for access.

The authorization API can be accessed at: <http://STARGATE_AUTH_DOMAIN/v1/auth>

Replace `STARGATE_AUTH_DOMAIN` in the example above with the DNS name and port. For example, when running on your localhost: 

http://localhost:8081/v1/auth

Detailed information about the Stargate auth API can be found in the [Stargate docs](https://stargate.io/docs/latest/secure/auth.html).

### Extracting Cassandra username/password Secrets

The auth API requires the Cassandra username and password to be provided to it.  Those values can be extracted from the K8ssandra cluster through the following commands (replace `k8ssandra` with the name configured for your running cluster).

Extract and decode the username secret:

```bash
kubectl get secret k8ssandra-superuser -o jsonpath="{.data.username}" | base64 --decode
```

Extract and decode the password secret:

```bash
kubectl get secret k8ssandra-superuser -o jsonpath="{.data.password}" | base64 --decode
```

### Generating Auth Tokens

Next, use the extracted and decoded secrets to request a token from the Stargate auth API.

```bash
curl -L -X POST 'http://_STARGATE_DOMAIN_/v1/auth' -H 'Content-Type: application/json' --data-raw '{"username": "k8ssandra-superuser", "password": "1LI8TebjjHYrqUk9xYbJnbYJheX3Ckq250byd2ePDPXNtweaYgznmg"}'
```

This request will return a response similar to the following. The value given for `authToken` will be required when making requests to the Stargate data APIs.

```json
{"authToken":"e4b34bbc-0ebc-4e2a-86ca-04793ca658a7"}
```

### Using Auth Tokens

Stargate supports authorization within the data APIs through a custom HTTP header `x-cassandra-token`, which must be populated with the token given by the auth API.

## Access Document Data API

The Stargate document APIs provide a way schemaless way to store and interact with data inside of Cassandra. The first step is to [create a namespace](https://stargate.io/docs/latest/quickstart/qs-document.html#creating-schema). That can be done with a request to the `/v2/schemas/namespaces` API:

```bash
curl --location --request POST 'http://STARGATE_REST_DOMAIN/v2/schemas/namespaces' \
--header "x-cassandra-token: e4b34bbc-0ebc-4e2a-86ca-04793ca658a7" \
--header 'Content-Type: application/json' \
--data '{
    "name": "mynamespace"
}'
```

Replace `STARGATE_REST_DOMAIN` in the example above with the DNS name and port. For example, when running on your localhost, specify as part of the `curl` command:

http://localhost:8082/v2/schemas/namespaces

The POST request will use the auth token previously generated to request the creation of a namespace called `mynamespace`. The server should return a response like:

```json
{"name":"mynamespace"}
```

Additional information related to using the Document APIs can be found in the Stargate [docs](https://stargate.io/docs/stargate/1.0/quickstart/quick_start-document.html).

## Access REST Data API

The Stargate REST APIs provide a RESTful way to store and interact with data inside of Cassandra that should feel familiar to developers. Unlike the document APIs, some understanding of Cassandra data modeling will be required. The first step is to [create a keyspace](https://stargate.io/docs/latest/quickstart/qs-rest.html#creating-schema). That can be done with a request to the `/v2/schemas/keyspaces` API:

```bash
curl --location --request POST 'http://STARGATE_REST_DOMAIN/v2/schemas/keyspaces' \
--header "x-cassandra-token: e4b34bbc-0ebc-4e2a-86ca-04793ca658a7" \
--header 'Content-Type: application/json' \
--data '{
    "name": "mykeyspace"
}'
```

Replace `STARGATE_REST_DOMAIN` in the example above with the DNS name and port. For example, when running on your localhost, specify as part of the `curl` command: 

http://localhost:8082/v2/schemas/keyspaces

The POST request will use the auth token previously generated to request the creation of a keyspace called `mykeyspace`. The server should return a response like:

```json
{"name":"mykeyspace"}
```

Additional information related to using the Document APIs can be found in the [Stargate docs](https://stargate.io/docs/latest/quickstart/qs-rest.html#creating-schema).

## Access GraphQL Data API

The Stargate GraphQL APIs provide a way to store and interact with data inside of Cassandra using the powerful GraphQL query language and tooling ecosystem. Like the REST APIs, this does require some additional Cassandra data modeling understanding. Like the REST APIs, The first step to using the GraphQL APIs is to [create a keyspace](https://stargate.io/docs/latest/quickstart/qs-graphql-cql-first.html#create-a-keyspace).

The easiest way to get started with the GraphQL APIs is to use the built-in GraphQL playground described in the next section.

Additional information related to using the Document APIs can be found in the [Stargate docs](https://stargate.io/docs/latest/quickstart/qs-graphql-cql-first.html).

### Access GraphQL playground

Stargate's GraphQL service provides an interactive "playground" application that can be used to interact with the GraphQL APIs.

The playground application can be accessed at <http://STARGATE_GRAPHQL_DOMAIN/playground>.

Replace `STARGATE_GRAPHQL_DOMAIN` in the example above with the DNS name and port. For example, when running on your localhost: 

http://localhost:8080/playground

Detailed information related to using the GraphQL playground can be found in the [Stargate docs](https://stargate.io/docs/latest/develop/tooling.html#using-the-graphql-playground).

## Next steps

* For comprehensive information about Stargate, visit the [stargate.io](https://stargate.io/) site.
* For details on the API calls, see the Stargate [API reference](https://stargate.io/docs/latest/api.html).
* For information about using a superuser and secrets with Stargate authentication, see [Stargate security]({{< relref "/tasks/secure/#stargate-security" >}}).
* Also see the topics covering other [components]({{< relref "/components/" >}}) deployed by K8ssandra. 
* For information on using additional deployed components, see the [Tasks]({{< relref "/tasks/" >}}) topics.
