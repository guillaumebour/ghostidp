# ghostidp
A mock Identity Provider to support development.

## Getting started

GhostIdP depends on [Ory Hydra](https://github.com/ory/hydra), which we will deploy either in Docker or in Kubernetes.

### Outside docker

In this scenario, we run `ghostidp` on the host, and Hydra in Docker.

1. Start Hydra with `task start-hydra`
2. Run `ghostidp` with `go run cmd/idp/main.go -debug -users-file config/users.yaml`
3. Create a client for your application (omit `id` and `secret` if you want `hydra` to generate them for you):
```bash
docker exec ghostidp_hydra hydra create client \
  --name "Demo client" \
  --endpoint http://127.0.0.1:4445 \
  --grant-type authorization_code,refresh_token \
  --response-type code,id_token \
  --format json \
  --scope openid --scope offline \
  --redirect-uri http://127.0.0.1:5050/callback \
  --skip-consent \ # Whether we trust the client and can skip the consent page.
  --id "[id]" \
  --secret "[secret]"
```
4. Use the resulting `client_id` and `client_secret` values in your application.

The token authorize path is `http://127.0.0.1:4444/oauth2/auth`.

To clean up, run `task clean-idp`


### Docker

In this scenario, we run both `ghostidp` and Hydra in Docker.

1. Start both service with `task start-idp`
2. Generate the client for your application like described above, and use the resulting `client_id` and `client_secret`.

The token authorize path is still `http://127.0.0.1:4444/oauth2/auth`.


### Kubernetes

Let's now run `ghostidp` in a local Kubernetes cluster. To make things a bit more interesting, we will add a Gateway that takes care of the TLS Termination.

1. Make sure that the domains `idp.dev.local` resolves to 127.0.0.1. On Unix, add the following entries to your `/etc/hosts`:
```
# Local Development
::1             idp.dev.local
127.0.0.1       idp.dev.local
```

2. Spin up the local Kubernetes infrastructure by running `task boostrap-dev-cluster`.
3. Verify that the cluster is up and running with `kubectl get pods --all-namespaces`.
4. Build the latest Docker image and push it to the local registry
```
docker build -t localhost:5005/ghostidp:latest -f Dockerfile . 
docker push localhost:5005/ghostidp:latest  
```

5. If you want the default OAuth2 client to be created automatically (`hydra.maester.enabled` is `true`), you also need to build the `hydra-maester` controller image
```
# From your project folder
git clone git@github.com:guillaumebour/hydra-maester.git
git checkout feat/155-allow-user-defined-credentials
docker build -t localhost:5005/controller:latest -f .docker/Dockerfile-build .    
docker push localhost:5005/controller:latest
```

6. Install GhostIdP to the cluster
```
helm install -f ./helm/values.https.yaml ghostidp ./helm/ghostidp 
```
