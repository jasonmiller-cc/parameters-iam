# parameters-iam

REST API for identity and access management across the parameters ecosystem.
Aggregates users, groups, roles, and policies from LDAP, Kerberos, and
related backends via [parameters-core](../parameters-core).

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/users | Provision user (LDAP + Kerberos) |
| GET | /api/v1/users | List users with roles |
| GET | /api/v1/users/{uid} | Get user with full profile |
| DELETE | /api/v1/users/{uid} | Deprovision user |
| POST | /api/v1/users/{uid}/roles | Assign roles to user |
| DELETE | /api/v1/users/{uid}/roles/{role} | Remove role from user |
| GET | /api/v1/roles | List roles |
| POST | /api/v1/roles | Create role |
| DELETE | /api/v1/roles/{name} | Delete role |
| GET | /api/v1/roles/{name}/permissions | List permissions for role |
| PUT | /api/v1/roles/{name}/permissions | Set permissions for role |
| POST | /api/v1/auth/token | Issue JWT (username + password) |
| POST | /api/v1/auth/refresh | Refresh JWT |
| POST | /api/v1/auth/revoke | Revoke token |
| GET | /healthz | Liveness probe |
| GET | /readyz | Readiness probe |
| GET | /version | Build version info |

## Configuration

```yaml
server:
  host: 0.0.0.0
  port: 8080

auth:
  jwt_secret: changeme
  jwt_issuer: parameters-iam

iam:
  ldap_service_url: http://parameters-ldap:8081
  kerberos_service_url: http://parameters-kerberos:8082
  ca_service_url: http://parameters-ca:8083
  default_roles:
    - viewer
  super_admin_group: iam-admins
```

Environment variable overrides use the `PARAMS_IAM_` prefix, e.g.
`PARAMS_IAM_CONFIG=/etc/iam/config.yaml`, `AUTH_JWT_SECRET=...`.

## Development

```bash
# Requires parameters-core checked out as a sibling directory.
make build
make test
make run
```
