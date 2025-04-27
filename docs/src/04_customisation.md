# Customisation

ghostidp's UI can be customized to better fit your needs.

This page describes the available environment variables that can be customized.

## Login and consent pages

### Adding a badge

You can add a badge with an arbitrary text in the header of the login and consent pages.
While this badge is intended to help distinguish between environment, how you use it is up to you.

To add a badge, set the `BADGE` environment variable to e.g. `DEV ENV`.

To set the color of the badge, refer to the Accent color section below.

### Accent color

You can set the accent color by setting the `ACCENT_COLOR` environment variable to a hex color code, e.g. `#008F8C`. 

### Version

You can add version information in the footer of the login and consent pages by setting the `VERSION` environment variable to e.g. `v1.0.0-beta`.
While the intended use is to display version information, how you use it is up to you. 

### Header Logo

TODO (not possible yet)

### Header Text

TODO (not possible yet)


### Log level

Set `DEBUG` to `true` if you want ghostidp's logs to be more granular.


## Users

### Description

Users' description text can be set when [configuring the users](03_managing_users_configuration.md), by setting the `description` field:

```yaml
users:
  - username: alice
    description: |
      A demo user called Alice
```

### Avatar text

TODO (not possible yet)

### Avatar color

TODO (not possible yet)
