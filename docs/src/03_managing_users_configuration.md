# Users configuration

Available users can be configured in the `users.yaml` config file provided to `ghostidp`:

```yaml
# This is an example config
users:
  - username: jeandupont
    description: |
      A demo user called Alice
    email: jean.dupont@example.com
    given_name: Jean
    family_name: Dupont
    custom_claims:
      roles:
        - admin
        - user
      department: engineering
      employee_id: "12345"
```

Where the fiels are used as follows:
- `username`: used in claims
- `description`: used in the UI
- `email`: used in claims
- `given_name`: used in claims
- `family_name`: used in claims
- `custom_claims`: used in claims

As the name indicates, `custom_claims` can be used to provide custom claims per users to fit your application needs.

The standard claims are created as follows:
- `sub`: mapped to `username`
- `email`: mapped to `email`
- `given_name`: mapped to `given_name`
- `family_name`: mapped to `family_name`
- `name`: `[given_name] [family_name]`

The custom claims are then added to the standard claims. The claims resulting from the above example are thus:
- `sub`: `jeandupont`
- `email`: `jean.dupont@example.com`
- `given_name`: `Jean`
- `family_name`: `Dupont`
- `name`: `Jean Dupont`
- `roles`: `["admin", "user"]`
- `department`: `engineering`
- `employee_id`: `12345`

As the name indicates, `custom_claims` can be used to provide custom claims per users to fit your application needs.