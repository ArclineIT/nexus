# Nexus Control Panel — TODO

**Product name:** Nexus Control Panel
**Subdomain:** `nexus.arcline.it`
**Tagline:** "Your central command. One login, every tool."

## Naming & Branding
- [ ] Confirm "Nexus Control Panel" clears trademark/legal check
- [ ] Register/reserve `nexus.arcline.it` subdomain and DNS entry
- [ ] Lock tagline copy across marketing site, login page, and docs
- [ ] Design logo/wordmark reflecting the "hub" metaphor
- [ ] Define brand voice guidelines tied to the "central command" concept

## Infrastructure & Identity Core
- [ ] Stand up identity service as the single source of truth
- [ ] Provision `nexus.arcline.it` with TLS certs
- [ ] Design core data model: users, roles, sessions, connected apps
- [ ] Implement session management / token issuance (JWT or equivalent)
- [ ] Set up audit logging for login events across connected tools

## Tool Integrations (the "spokes" of the hub)
- [ ] **Portal** — connect via SSO, map roles/permissions
- [ ] **Billing** — connect via SSO, map roles/permissions
- [ ] **Git** — connect via SSO, map roles/permissions
- [ ] **Monitoring** — connect via SSO, map roles/permissions
- [ ] Define standard integration pattern/spec so future tools plug in consistently

## Third-Party SSO (future-proofing per the "Nexus" metaphor)
- [ ] Research Google Workspace SSO integration requirements
- [ ] Research Microsoft Entra ID (Azure AD) integration requirements
- [ ] Design abstraction layer so external IdPs can plug in without core rework
- [ ] Prioritize/schedule third-party SSO for a post-launch phase

## UI / UX
- [ ] Design login/landing page reflecting "central command" positioning
- [ ] Build dashboard showing all connected tools (Portal, Billing, Git, Monitoring)
- [ ] Add per-user access/permissions view
- [ ] Add admin view for managing connected apps and users

## Security
- [ ] Define password/MFA policy
- [ ] Pen-test the identity core before launch
- [ ] Set up rate limiting and anomaly detection on login endpoint
- [ ] Document token expiry/refresh and revocation flows

## Launch
- [ ] Internal beta with Portal + Billing only
- [ ] Roll out Git + Monitoring integrations
- [ ] Publish tagline/messaging externally
- [ ] Announce `nexus.arcline.it` as the official login hub

## Documentation
- [ ] Write internal integration guide (how to plug a new tool into Nexus)
- [ ] Write end-user help doc ("What is Nexus Control Panel?")
- [ ] Write admin guide for managing users/roles
