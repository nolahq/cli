# `nola-cli`

The Nola CLI can be used to make API requests against the Nola API.

The Nola API uses short-lived, OAuth 2.0 Bearer tokens for authentication. To make it easier to use the API, the Nola CLI can automatically generate and insert these tokens into your command-line CLI requests.

## Getting Started

### Download the CLI

TODO

### Login

The CLI supports the live, staging and development (localhost) Nola server environments. This is controlled by the `--server`/`-s` flag. e.g. `nola --server live` (default), `nola --server staging`, `nola --server dev`.

The CLI also supports the use of profiles. By default, the profile name will be determined by the server, e.g. `default-live`, `default-staging`, `default-dev`. You can override this by using the `--profile`/`-p` flag. e.g. `nola --profile my-profile`.

You will need to login to the Nola API before you can make requests. To do this, run (optionally providing the `--server` and/or `--profile` flags)

```bash
nola login
```

### Making `curl` requests
If this succeeds, you can then make requests to the Nola API (again, optionally setting the server or profile):

```bash
nola curl -i https://app.nolahq.com/api/v1/companies
```

### Getting the `Authorization` token
If for some reason you need the authorization token, you can run:

```bash
nola token
```

In `bash` or PowerShell you can do things like:

```bash
curl -H "Authorization: Bearer $(nola token)" https://...
```