# Nexus Control Panel

**Your central command. One login, every tool.**

Nexus is the identity hub for the Arcline platform. It is an OpenID Connect
(OIDC) identity provider and the single source of truth for users, roles,
sessions and connected applications. Every Arcline tool — Portal, Billing, Git,
Monitoring and the rest — authenticates against Nexus instead of keeping its own
user store.

## Stack

- **Ruby on Rails 8.1** (Ruby 4.0)
- **SQLite** — single-file database, no server required
- **Doorkeeper** + **doorkeeper-openid_connect** — OAuth2 / OIDC authorization server
- **bcrypt** — password hashing (`has_secure_password`)
- **Hotwire** (Turbo + Stimulus) with **importmap** and **Propshaft**
- **Kamal** + **Thruster** for containerised deployment

## Quick start

```bash
bin/setup            # install gems, prepare the database, seed data
bin/rails server     # http://localhost:3000
```

`bin/setup` seeds an admin user (`admin@arcline.it` / `password123` by default —
override with `NEXUS_ADMIN_EMAIL` / `NEXUS_ADMIN_PASSWORD`), the default roles
and permissions, and all Arcline applications as OAuth clients.

## Configuration

Configuration is environment-driven (see `.env.example`). In production these are
injected by Kamal (`config/deploy.yml`).

| Variable                 | Default                         | Description                                             |
|--------------------------|---------------------------------|---------------------------------------------------------|
| `NEXUS_ISSUER`           | `https://nexus.arcline.it`      | Canonical public URL / OIDC issuer. Must be exact.      |
| `NEXUS_HOST`             | `nexus.arcline.it`              | Host used for documentation and discovery references.   |
| `NEXUS_OIDC_PRIVATE_KEY` | *(generated in dev)*            | RSA private key (PEM) used to sign ID tokens.           |
| `DATABASE_PATH`          | `storage/production.sqlite3`    | SQLite database file.                                   |
| `NEXUS_ADMIN_EMAIL`      | `admin@arcline.it`              | Bootstrap admin email used by `db:seed`.                |
| `NEXUS_ADMIN_PASSWORD`   | `password123` (dev/test only)   | Bootstrap admin password. Required in production.       |
| `NEXUS_MAIL_FROM`        | `Nexus <noreply@arcline.it>`    | "From" address for transactional email.                 |

### Signing key

ID tokens are signed with an RSA key. Locally a throwaway key is generated at
`config/keys/oidc.pem` (gitignored) on first boot. In production provide a key
via `NEXUS_OIDC_PRIVATE_KEY`:

```bash
ruby -ropenssl -e 'puts OpenSSL::PKey::RSA.new(2048).to_pem'
```

Rotating the key invalidates existing ID tokens but not access tokens; clients
refresh the key from the JWKS endpoint automatically.

## Connecting an application (SSO)

Nexus is a standard OpenID Connect provider, so any OAuth2/OIDC client library
works. Each app is registered as a **Connected App** (client) in the admin UI,
which issues a client ID and secret.

**Discovery document**

```
GET https://nexus.arcline.it/.well-known/openid-configuration
```

**Endpoints**

| Endpoint            | Path                                       |
|---------------------|--------------------------------------------|
| Authorization       | `/oauth/authorize`                         |
| Token               | `/oauth/token`                             |
| UserInfo            | `/oauth/userinfo`                          |
| JWKS                | `/oauth/discovery/keys`                    |
| Token introspection | `/oauth/introspect`                        |
| Token revocation    | `/oauth/revoke`                            |

**Scopes & claims**

| Scope            | Claims                                              |
|------------------|-----------------------------------------------------|
| `openid`         | `sub` (required)                                    |
| `email`          | `email`, `email_verified`                           |
| `profile`        | `name`, `preferred_username`                        |
| `roles`          | `roles` (array), `permissions` (UserInfo only)      |
| `offline_access` | refresh token                                       |

**Flow (Authorization Code).** Redirect the browser to `/oauth/authorize` with
`response_type=code`, `client_id`, `redirect_uri`, `scope` and `state`. After the
user signs in (existing Nexus sessions are reused, so they see no extra prompt),
exchange the returned `code` at `/oauth/token` for an access token, ID token and
refresh token. Verify the ID token against the JWKS and read additional claims
from `/oauth/userinfo`.

`sub` is the stable user identifier — always key users by `sub`, not by email.

## Roles & permissions

Nexus ships four roles, published to connected apps in the `roles` claim:

| Role       | Permissions                                     |
|------------|-------------------------------------------------|
| `admin`    | `apps.read`, `apps.manage`, `users.manage`, `audit.read` |
| `operator` | `apps.read`, `audit.read`                        |
| `member`   | `apps.read`                                      |
| `viewer`   | `apps.read`                                      |

Admins manage users and roles under `/admin/users`, and connected apps under
`/admin/connected_apps` (client credentials and secret rotation).

## Development

```bash
bin/rails test        # test suite
bin/rubocop           # style
bin/brakeman          # security scan
bin/ci                # full CI pipeline
```

## Project structure

```
app/
  controllers/        Session, registration, password, dashboard, admin
  controllers/concerns/authentication.rb   Cookie session auth
  models/             User, Session, Role, Permission, AuditLog, ConnectedApp
  views/              Nexus UI + Doorkeeper authorization pages
config/
  initializers/doorkeeper.rb                OAuth2 provider config
  initializers/doorkeeper_openid_connect.rb OIDC issuer, signing key, claims
db/seeds.rb           Roles, permissions, admin user, all Arcline apps
test/                 Model, controller and OIDC flow tests
```
