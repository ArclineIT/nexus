# Nexus Control Panel — TODO

**Product name:** Nexus Control Panel
**Subdomain:** `nexus.arcline.it`
**Stack:** Rails 8.1 + SQLite, OAuth2/OIDC provider (Doorkeeper)

## Done
- [x] Port from Go to Ruby on Rails
- [x] Cookie-session authentication (login, signup, password reset)
- [x] Roles & permissions (RBAC) model
- [x] Connected-app registry + admin UI
- [x] OpenID Connect provider (discovery, JWKS, userinfo, authorization code)
- [x] Audit logging for login / signup / password reset / admin events
- [x] Seed all Arcline applications as OAuth clients
- [x] OIDC end-to-end test (authorize → token → ID token → userinfo)

## Integration (the "spokes" of the hub)
- [ ] Portal — wire OIDC client, map roles/permissions
- [ ] Billing — wire OIDC client
- [ ] Git (Forgejo) — configure OAuth2 authentication source
- [ ] Monitoring / Status — wire OIDC client
- [x] Docs — OIDC client wired (sign-in, `admin` role, PKCE)
- [ ] Email, Provisioner, Uptime, Status — replace HTTP basic auth with an OIDC client
- [ ] Migrate, Vault, Audit, Check, DNS, Network — as each needs sign-in
- [ ] Rotate seeded client secrets and distribute to each app

## Identity hardening
- [ ] MFA (TOTP) for user accounts
- [ ] Email verification on signup
- [ ] Session management UI (list/revoke active devices)
- [ ] Rate limiting / anomaly detection on authorization endpoints
- [ ] Document key rotation procedure for the OIDC signing key

## UI / UX
- [ ] Per-user "my connected apps" consent view
- [ ] Audit log viewer in the admin UI
- [ ] Brand assets (logo/wordmark) and dark/light theming

## Third-party SSO (future)
- [ ] Google Workspace as an upstream IdP
- [ ] Microsoft Entra ID (Azure AD) as an upstream IdP
- [ ] Federation / account-linking abstraction

## Launch
- [ ] TLS for `nexus.arcline.it` and DNS
- [ ] Pen-test before rollout
- [ ] Internal beta, then migrate tools app by app
