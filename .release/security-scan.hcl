# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

binary {
    go_stdlib  = true // Scan the Go standard library used to build the binary.
    go_modules = true // Scan the Go modules included in the binary.
    osv        = true // Use the OSV vulnerability database.
    oss_index  = true // And use OSS Index vulnerability database.

    secrets {
        all = true
    }

    triage {
        suppress {
            // Add known false positive CVE IDs here as they are identified
            // during CRT security scans. Leave empty on first release.
            vulnerabilities = []
        }
    }
}

container {
    dependencies = true // Scan any installed UBI packages for vulnerabilities.
    osv          = true // Use the OSV vulnerability database.

    secrets {
        all = true
    }

    triage {
        suppress {
            // Add known false positive CVE IDs here as they are identified
            // during CRT security scans. Leave empty on first release.
            vulnerabilities = []
        }
    }
}
