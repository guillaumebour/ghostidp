# ghostidp

A mock Identity Provider to support development.

## Concept

_ghostidp_ was created to fulfill a development need: quickly spin up an OpenID Connect-compliant Identity Provider (IdP) with hard-coded demo users.

The goal is to provide an easy-to-setup and configure IdP tailored for development and testing. By preloading users with customizable claims, developers can simulate different roles, permissions, or even multiple identity providers.

Under the hood, _ghostidp_ is a custom UI for the login and consent screens of [Ory Hydra](https://github.com/ory/hydra).

There are no passwords in _ghostidp_. Instead, users are selected from a list, enabling fast switching between identities. 
Combined with a stateless (sessionless) design, this makes it simple to test scenarios involving different users without any friction.

![ghostidp screenshots](screenshots.png)

## Getting started

In the examples below, Hydra is deployed using Docker or Kubernetes.

### Configuring Users

Available users are provided to ghostidp via a config file such as:

```yaml
users:
  - username: alice
    display:
      description: A demo user called Alice
    email: alice@example.com
    given_name: Alice
    family_name: Smith
    custom_claims:
      roles:
        - admin
        - user
      department: engineering
      employee_id: "12345"
```

`username`, `email`, `given_name` and `family_name` are used to populate the standard claims of the JWT token.
Everything under `custom_claims` will be added as-is to the JWT token. 

You can use `custom_claims` to provide claims that will fit your application's needs.

The `description` field is used to provide a description of the user in the UI users' list.

### Docker

To run both `ghostidp` and Hydra in Docker:

1. Start both services with `task start-idp`
2. Create a client for your application (omit `id` and `secret` if you want `hydra` to generate them for you):
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
3. Use the resulting `client_id` and `client_secret` values in your application.

The Auth URL is `http://127.0.0.1:4444/oauth2/auth`, and the Token URL is `http://127.0.0.1:4444/oauth2/token`.

To clean up, run `task clean-idp`

### Kubernetes

To run `ghostidp` in a local Kubernetes cluster:

1. Make sure that the domain `idp.dev.local` resolves to 127.0.0.1. On Unix, add the following entries to your `/etc/hosts`:
```
# Local Development
::1             idp.dev.local
127.0.0.1       idp.dev.local
```

2. Spin up the local Kubernetes infrastructure by running `task boostrap-dev-cluster`.

> [!NOTE]  
> To make things a bit more interesting, this script also setups a Gateway that takes care of the TLS Termination.

3. Verify that the cluster is up and running with `kubectl get pods --all-namespaces`.

4. If you want the default OAuth2 client to be created automatically (`hydra.maester.enabled` is `true`), you also need to build the `hydra-maester` controller image (see https://github.com/ory/hydra-maester/pull/159).
```
# From your project folder
git clone git@github.com:guillaumebour/hydra-maester.git
git checkout feat/155-allow-user-defined-credentials
docker build -t localhost:5005/controller:latest -f .docker/Dockerfile-build .    
docker push localhost:5005/controller:latest
```

5. Install ghostidp to the cluster
```
helm install -f ./helm/values.https.yaml ghostidp ./helm/ghostidp 
```

### Running outside docker

In this scenario, mostly relevant for developing ghostidp, we run `ghostidp` on the host, and Hydra in Docker.

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

The Auth URL is `http://127.0.0.1:4444/oauth2/auth`, and the Token URL is `http://127.0.0.1:4444/oauth2/token`.

To clean up, run `task clean-idp`
