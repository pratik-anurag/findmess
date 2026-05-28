# Anti-Stalking

FindMesh is designed for lost-item recovery, not person tracking.

## Mobile Protections

The app tracks unknown FindMesh-compatible advertisements observed across local time windows. If the same unknown ephemeral ID pattern appears repeatedly near the user over a meaningful duration, the app raises an alert and offers an abuse report flow.

Because rotating IDs intentionally prevent easy cross-time linkage, this repository implements a practical prototype and documents the limitation. Production should integrate with platform unwanted-tracker standards when available.

## Tag Protections

The lost-tag firmware includes separation logic. If a paired tag is away from its owner for a prolonged period, it can periodically sound and expose NFC disable instructions. Exact thresholds should be tuned by safety policy and jurisdiction.

## Abuse Workflow

Users can file abuse reports. Safety reviewers can investigate suspicious repeated reports, disable tags, disable stands, or escalate account action. All admin actions are audited.

## User Language

Product copy uses lost-item recovery language and avoids normalizing people tracking. Finder participation is opt-in and can be disabled.
