# ghostidp
A mock Identity Provider to support development.

Generating the client in hydra:

```bash
hydra create client --name "My test client" --endpoint http://127.0.0.1:4445 --grant-type authorization_code,refresh_token --response-type code,id_token --format json --scope openid --scope offline --redirect-uri http://127.0.0.1:5050/callback --id "[id]" --secret "[secret]"
```