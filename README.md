# Provision HTTP example

This repository contains a deliberately small, framework-neutral HTTP application used to exercise [Provision](https://github.com/datashaman/provision) deployments. It is an example workload with its own release lifecycle, not part of the Provision product distribution.

The application exposes:

- `GET /live` — process liveness;
- `GET /ready` — traffic readiness;
- `GET /verify` — the value of `PROVISION_REVISION` as JSON;
- `GET /slow?seconds=N` — a bounded delayed response that identifies the revision which accepted it, for graceful-switch acceptance tests;
- `GET /` — a plain-text identification response.

It listens on `PROVISION_HTTP_LISTEN`, defaulting to `127.0.0.1:18081`. Build a reproducible native Linux bundle with:

For rollback acceptance tests, a `PROVISION_REVISION` ending in `-fail-stable` deliberately returns `503` from the three health endpoints when the request Host differs from the private `PROVISION_HTTP_LISTEN` address. Direct candidate verification therefore succeeds before a traffic switch, while verification through the stable proxy fails reproducibly. Other revisions retain the normal behavior.

```sh
./build.sh amd64
```

Releases follow semantic versioning. Each release publishes the Linux bundle and its SHA-256 digest for use by a Provision Revision.

Licensed under the [MIT License](LICENSE).
