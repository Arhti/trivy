# CentOS
Trivy supports the following scanners for OS packages.

|    Scanner    | Supported |
| :-----------: | :-------: |
|     SBOM      |     ✓     |
| Vulnerability |     ✓     |
|    License    |     ✓     |

Please see [here](index.md#supported-os) for supported versions.

The table below outlines the features offered by Trivy.

|                Feature                | Supported |
| :-----------------------------------: | :-------: |
|        Unfixed vulnerabilities        |     ✓     |
| [Dependency graph][dependench-graph] |     ✓     |

## SBOM
Same as [RHEL](rhel.md).

## Vulnerability
CentOS does not provide straightforward machine-readable security advisories.
As a result, Trivy utilizes the security advisories from [Red Hat Enterprise Linux (RHEL)](rhel.md) for detecting vulnerabilities in CentOS.
This approach might lead to situations where, even though Trivy displays a fixed version, CentOS might not have the patch available yet.
Since patches released for RHEL often become available in CentOS after some time, it's usually just a matter of waiting.

!!! note
    The case for CentOS Stream, which is not supported by Trivy, is entirely different from CentOS.

As Trivy relies on Red Hat's advisories, please refer to [Red Hat](rhel.md) for details regarding vulnerability severity and status.


## License
Same as [RHEL](rhel.md).


[dependench-graph]: ../../configuration/reporting.md#show-origins-of-vulnerable-dependencies