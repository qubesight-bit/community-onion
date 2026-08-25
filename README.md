# Community Onion

Federated, JavaScript-free community platform designed for censorship-resistant publishing and preservation over Tor.

## Status

Early experimental prototype. Do not use this version for sensitive communications, whistleblowing, or high-risk communities.

## Current foundation

- Self-hosted Tor Onion Service.
- Minimal Go HTTP application.
- Server-rendered HTML without JavaScript.
- No trackers or third-party resources.
- Security-focused HTTP headers.
- Unprivileged and read-only Docker runtime.
- Host exposure restricted to `127.0.0.1`.

## Planned features

PostgreSQL, authentication, communities, moderation, federation, document verification and distributed preservation.

## Security

Never commit Tor identities, SSH private keys, credentials, database files, uploads or encrypted backups.
