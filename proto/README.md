# proto

This directory defines the public api surface for temporal-shop.

## buf
[buf](https://docs.buf.build/tour/configure-and-build) is adopted for linting/generating messages.
As such, buf's [style guide](https://docs.buf.build/format/style) is used for messages. Temporal workflow/activity 
contracts should always use protobuf messages.

You can look in the SA 1Password account for the BSR token if needed under `sa_buf`.

### Generating Protobufs

The `buf.gen.yaml` file defines generation rules.

```shell
cd proto && buf generate
```

This generates the `api` root level directory and exposes our versioned, public api.

You can also just

```shell
make genapi
```

### Linting

The `buf.yaml` file defines linting rules.

```shell
cd proto && buf lint
```

### Source Control

Right now, we don't need to use the buf [GH Actions](https://docs.buf.build/ci-cd/github-actions)
since we should be generating our api and checking into source control.