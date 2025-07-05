# Docker

In this guide, we will deploy Ory Hydra and _ghostidp_ with Docker.

## Hydra Configuration

Start by creating the configuration file for Hydra:

```yaml
# config/hydra/hydra.yaml
serve:
  cookies:
    same_site_mode: Lax

urls:
  self:
    issuer: http://127.0.0.1:4444
  consent: http://127.0.0.1:8080/consent # ghostidp Login Endpoint
  login: http://127.0.0.1:8080/login     # ghostidp Consent Endpoint

secrets:
  system:
    - youReallyNeedToChangeThis

oidc:
  subject_identifiers:
    supported_types:
      - pairwise
      - public
    pairwise:
      salt: youReallyNeedToChangeThis
```

## GhostIdP Configuration

Create a config for _ghostidp_, containing the hard-coded users:

```yaml
# config/users.yaml
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
  - username: bob
    display:
      description: A demo user called Bob
    email: bob@example.com
    given_name: Bob
    family_name: Johnson
    custom_claims:
      roles:
        - user
      department: marketing
      employee_id: "67890"
```

## Running GhostIdP

Create a Docker Compose file:

```yaml
# docker-compose.yaml
services:
  sqlite:
    image: busybox
    volumes:
      - hydra-sqlite:/mnt/sqlite
    command: "chmod -R 777 /mnt/sqlite"
  hydra:
    container_name: ghostidp_hydra
    image: oryd/hydra:v2.3.0
    ports:
      - "4444:4444" # Public port
    command: serve -c /etc/config/hydra/hydra.yml all --dev
    volumes:
      - hydra-sqlite:/mnt/sqlite:rw
      - type: bind
        source: ./config/hydra/
        target: /etc/config/hydra
    pull_policy: missing
    environment:
      - DSN=sqlite:///mnt/sqlite/db.sqlite?_fk=true&mode=rwc
    restart: unless-stopped
    depends_on:
      - hydra-migrate
      - sqlite
  hydra-migrate:
    image: oryd/hydra:v2.3.0
    environment:
      - DSN=sqlite:///mnt/sqlite/db.sqlite?_fk=true&mode=rwc
    command: migrate -c /etc/config/hydra/hydra.yml sql up -e --yes
    pull_policy: missing
    volumes:
      - hydra-sqlite:/mnt/sqlite:rw
      - type: bind
        source: ./config/hydra/
        target: /etc/config/hydra
    restart: on-failure
    depends_on:
      - sqlite
  ghostidp:
    container_name: ghostidp
    image: ghcr.io/guillaumebour/ghostidp:latest
    volumes:
      - ./config/users.yaml:/users.yaml
    environment:
      HYDRA_ADMIN_URL: http://ghostidp_hydra:4445/admin
      USERS_FILE: users.yaml
    ports:
      - "8080:8080"
    depends_on:
      - hydra
volumes:
  hydra-sqlite:
```

Start Hydra and _ghostidp_ with `docker compose up -d`. After a few seconds, both Hydra and _ghostidp_ should be ready.

Create an OAuth2 Client for your application, here to perform an Authorization Code Flow (see [Hydra's Documentation](https://www.ory.sh/docs/hydra/cli/hydra-create-client) for the full reference of Hydra's CLI).

```bash
docker exec ghostidp_hydra hydra create client \
  --name "Demo client" \
  --endpoint http://127.0.0.1:4445 \
  --grant-type authorization_code,refresh_token \
  --response-type code,id_token \
  --format json \
  --scope openid --scope offline \
  --redirect-uri http://127.0.0.1:5050/callback \
  --skip-consent \               # Whether you trust the client and want to skip the consent page
  --id "$YOUR_CLIENT_ID" \       # Omit to let Hydra create it for you
  --secret "$YOUR_CLIENT_SECRET" # Omit to let Hydra create it for you
```

Use the resulting `client_id` and `client_secret` in your application.

The URLs are:
- Auth URL: [http://127.0.0.1:4444/oauth2/auth](http://127.0.0.1:4444/oauth2/auth).
- Token URL: [http://127.0.0.1:4444/oauth2/auth](http://127.0.0.1:4444/oauth2/auth).