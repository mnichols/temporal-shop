# Temporal Shop

This repo uses go workspaces.
But there are issues with `go mod tidy` in underlying go modules:
https://github.com/golang/go/issues/50750

### Deployment

**UI**
1. Prerender static svelte app
1. Deploy to s3 bucket `s3://temporal-sa/temporal-shop/app`
   1. Note the bucket policy [here](https://s3.console.aws.amazon.com/s3/buckets/temporal-sa?region=us-east-1&tab=permissions) that allows objects having tag `security:public` (key:value) to be publicly accessible.
   1. The actual command to put tags on the relevant objects looks like [this](#tagging-s3-objects-for-public-access):
1. Configure traefik to point to s3 bucket
   1. A `ui` [ExternalName service](kustomize/web/overlays/prod/service-ui.yml) is created for traefik to use in the [router](kustomize/web/overlays/prod/ingress-route.yml) . This effectively acts as a CNAME to our s3 bucket for delivery of objects.Idea here https://stackoverflow.com/questions/62242278/traefik-as-a-reverse-proxy-for-s3-static-website
      1. Idea [here](https://stackoverflow.com/questions/62242278/traefik-as-a-reverse-proxy-for-s3-static-website)
      1. https://stackoverflow.com/questions/72218644/traefik-proxy-to-s3-site

### Tagging S3 objects for public access

```shell
aws s3api list-objects --bucket temporal-sa \
--prefix temporal-shop/app \
--query 'Contents[].[Key]' --output text | xargs -n 1 -P 10 aws s3api put-object-tagging --bucket temporal-sa --tagging 'TagSet=[{Key=security,Value=public}]' --key
```

*References*

- https://aws.amazon.com/premiumsupport/knowledge-center/read-access-objects-s3-bucket/
- https://www.learnaws.org/2022/08/22/tag-objects-s3/

## Development

### Running The Things

**Golang Services**

`cd services/go && go run cmd/temporal_shop/main.go`

**Golang Web BFF**

`cd web && go run bff/cmd/bff/main.go`

**Svelte UI**

`cd web/ui && npm run dev`

### Generating The Things

**Protobufs API**

`make genapi`

**Graphql API**

`make gengql`

### Debugging The Things

**grpcui** 

You can interact with grpc services the workflow(s) consume with this tool.

- [installation](go install github.com/fullstorydev/grpcui/cmd/grpcui@latest)
- `cd services/go && go run cmd/temporal_shop/main.go`
- `grpcui -plaintext localhost:9000`

**graphiql dashboard** 

`gqlgen` exposes a http handler for Graphiql that allows you to dive into our GraphQL models.

1. Make `web/bff/.env` environment variable `HTTP_SERVER_SHOWS_GRAPHQL_PLAYGROUND=true`
2. Restart the bff server
3. Visit `localhost:8080/api/gql` in your browser
