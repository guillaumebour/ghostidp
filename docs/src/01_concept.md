# Concept

_ghostidp_ is built on top of Ory Hydra as an implementation of the "Ory OAuth 2.0 login & consent flow."

Ory Hydra is an OpenID Certified OAuth 2.0 Server and OpenID Connect Provider, and as such, it doesn't contain a database with end users.
Instead, it delegates the login and consent flow to a dedicated application (in our case, _ghostidp_).

The full explanation of how this works from Hydra's perspective is available in Hydra's documentation: [User login and consent flow](https://www.ory.sh/docs/oauth2-oidc/custom-login-consent/flow).

## Sequence Diagram

Here is a sequence Diagram of what happens, adapted from the [Ory Hydra Documentation](https://www.ory.sh/docs/oauth2-oidc/custom-login-consent/flow#sequence-diagram):

```mermaid
sequenceDiagram
    OAuth2 client->>Ory OAuth2 and OpenID Connect: Initiates OAuth2 Authorize Code or Implicit Flow
    Ory OAuth2 and OpenID Connect-->>Ory OAuth2 and OpenID Connect: No end user session available (not authenticated)
    Ory OAuth2 and OpenID Connect->>Login Endpoint (ghostidp): Redirects end user with login challenge
    Login Endpoint (ghostidp)-->Ory OAuth2 and OpenID Connect: Fetches login info
    Login Endpoint (ghostidp)-->>Login Endpoint (ghostidp): Authenticates selected user
    Login Endpoint (ghostidp)-->Ory OAuth2 and OpenID Connect: Transmits login info and receives redirect url with login verifier
    Login Endpoint (ghostidp)->>Ory OAuth2 and OpenID Connect: Redirects end user to redirect url with login verifier
    Ory OAuth2 and OpenID Connect-->>Ory OAuth2 and OpenID Connect: First time that client asks user for permissions
    Ory OAuth2 and OpenID Connect->>Consent Endpoint (ghostidp): Redirects end user with consent challenge
    Consent Endpoint (ghostidp)-->Ory OAuth2 and OpenID Connect: Fetches consent info (which user, what app, what scopes)
    Consent Endpoint (ghostidp)-->>Consent Endpoint (ghostidp): Asks for end user's permission to grant application access
    Consent Endpoint (ghostidp)-->Ory OAuth2 and OpenID Connect: Transmits consent result and receives redirect url with consent verifier
    Consent Endpoint (ghostidp)->>Ory OAuth2 and OpenID Connect: Redirects to redirect url with consent verifier
    Ory OAuth2 and OpenID Connect-->>Ory OAuth2 and OpenID Connect: Verifies grant
    Ory OAuth2 and OpenID Connect->>OAuth2 client: Transmits authorization code/token
```

