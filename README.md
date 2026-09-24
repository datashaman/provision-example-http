# Provision HTTP example

This repository contains a deliberately small, framework-neutral HTTP application used to exercise [Provision](https://github.com/datashaman/provision) deployments. It is an example workload with its own release lifecycle, not part of the Provision product distribution.

The application exposes:

- `GET /live` — process liveness;
- `GET /ready` — traffic readiness;
- `GET /verify` — the value of `PROVISION_REVISION` as JSON;
- `GET /` — a plain-text identification response.

It listens on `PROVISION_HTTP_LISTEN`, defaulting to `127.0.0.1:18081`. Build a reproducible native Linux bundle with:

```sh
./build.sh amd64
```

Releases follow semantic versioning. Each release publishes the Linux bundle and its SHA-256 digest for use by a Provision Revision.

Licensed under the [MIT License](LICENSE).
