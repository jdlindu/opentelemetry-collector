# HTTPS Provider

<!-- status autogenerated section -->
| Status        |           |
| ------------- |-----------|
| Stability     | [stable]  |
| Distributions | [core], [contrib], [k8s] |
| Issues        | [![Open issues](https://img.shields.io/github/issues-search/open-telemetry/opentelemetry-collector?query=is%3Aissue%20is%3Aopen%20label%3Aprovider%2Fhttpsprovider%20&label=open&color=orange&logo=opentelemetry)](https://github.com/open-telemetry/opentelemetry-collector/issues?q=is%3Aopen+is%3Aissue+label%3Aprovider%2Fhttpsprovider) [![Closed issues](https://img.shields.io/github/issues-search/open-telemetry/opentelemetry-collector?query=is%3Aissue%20is%3Aclosed%20label%3Aprovider%2Fhttpsprovider%20&label=closed&color=blue&logo=opentelemetry)](https://github.com/open-telemetry/opentelemetry-collector/issues?q=is%3Aclosed+is%3Aissue+label%3Aprovider%2Fhttpsprovider) |

[stable]: https://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/component-stability.md#stable
[core]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
[k8s]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-k8s
<!-- end autogenerated section -->

## Overview

The HTTPS Provider takes an HTTPS URI to a file and reads its contents as YAML
to provide configuration to the Collector. The validity of the certificate of
the HTTPS endpoint is verified when making the connection.

## Usage

The scheme for this provider is `https`. Usage looks like the following passed
to the Collector's command line invocation:

```text
--config=https://example.com/config.yaml
```

### Notes

The provider currently only supports communicating with servers whose
certificate can be verified using the root CA certificates installed in the
system. The process of adding more root CA certificates to the system is
Operating System-dependent. For Linux, please refer to the `update-ca-trust`
command.
