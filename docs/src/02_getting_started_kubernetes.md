# Kubernetes

> [!WARNING]
> The Helm Chart is still work in progress.

In this guide, we will deploy Ory Hydra and ghostidp on a Kubernetes cluster.

This assumes you already have a cluster up and running, such as a [Kind](https://kind.sigs.k8s.io/) cluster.

Note that we here deploy Hydra and ghostidp together, with Hydra being a dependency of ghostidp.
If you need more freedom/configuration options when deploying ghostidp (e.g. if you deploy it to EKS or GKE), you might want to deploy Hydra and ghostidp separately.

## Deploying the chart

Start by creating a value file:

```yaml
# values.yaml
ghostidp:
  usersConfig: |
    users:
      - username: alice
        description: A demo user.
        email: alice@example.com
        given_name: Alice
        family_name: Smith
        custom_claims:
          department: engineering
          employed_id: 12345
          roles:
            - admin
            - user
      - username: bob
        description: A demo user.
        email: bob@example.com
        given_name: Bob
        family_name: Marvel
        custom_claims:
          department: engineering
          employed_id: 6587
          roles:
            - user

hydra:
  serve:
    cookies:
      same_site_mode: Lax
  hydra:
    dev: true # Required to use Hydra without https
    config:
      dsn: "memory"
      urls:
        self:
          public: "http://127.0.0.1:4444"
          issuer: "http://127.0.0.1:4444"         
        login: "http://127.0.0.1:8080/login"      
        consent: "http://127.0.0.1:8080/consent"
```

Install the Helm Chart to the cluster:

```bash
helm install -f values.yaml ghostidp ./helm/ghostidp
```

Check that everything is running, it should look similar to:

```bash
NAME                    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
ghostidp                ClusterIP   10.96.25.18     <none>        8080/TCP   2m7s
ghostidp-hydra-admin    ClusterIP   10.96.187.133   <none>        4445/TCP   2m7s
ghostidp-hydra-public   ClusterIP   10.96.173.142   <none>        4444/TCP   2m7s
kubernetes              ClusterIP   10.96.0.1       <none>        443/TCP    16m
```

## Create the client 

```bash
kubectl get pods -l app.kubernetes.io/name=hydra -o name | xargs -I{} kubectl exec {} -- \ 
  hydra create client \
  --name "Demo Client" \
  --endpoint http://127.0.0.1:4445 \
  --grant-type authorization_code,refresh_token \
  --response-type code,id_token \
  --format json \
  --scope openid --scope offline \
  --redirect-uri http://127.0.0.1:5555/callback
```

## Access the Identity Provider

Forward the Hydra's public service and the ghostidp one:

```bash
kubectl port-forward service/ghostidp-hydra-public 4444:4444
```
and
```bash
kubectl port-forward service/ghostidp 8080:8080
```

As in the [Docker Example](02_getting_started_docker.md), the URLs are:
- Auth URL: [http://127.0.0.1:4444/oauth2/auth](http://127.0.0.1:4444/oauth2/auth).
- Token URL: [http://127.0.0.1:4444/oauth2/auth](http://127.0.0.1:4444/oauth2/auth).
