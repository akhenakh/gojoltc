#!/usr/bin/env bash
# Generate sbom.json (SPDX-lite) from the module and vendored submodules.
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

module="$(awk '$1 == "module" { print $2; exit }' go.mod)"
gover="$(awk '$1 == "go" { print $2; exit }' go.mod)"

joltc_sha="$(git -C third_party/joltc rev-parse HEAD)"
joltc_url="$(git -C third_party/joltc remote get-url origin)"
jolt_sha="$(git -C third_party/JoltPhysics rev-parse HEAD)"
jolt_ver="$(git -C third_party/JoltPhysics describe --tags 2>/dev/null || echo "$jolt_sha")"
jolt_url="$(git -C third_party/JoltPhysics remote get-url origin)"

# Normalize SSH remotes to https:// URLs.
normalize_url() {
  local u="$1"
  u="${u#git@}"
  u="${u/:/\/}"
  u="${u%.git}"
  [[ "$u" == http* ]] || u="https://$u"
  printf '%s' "$u"
}
joltc_url="$(normalize_url "$joltc_url")"
jolt_url="$(normalize_url "$jolt_url")"

cat > sbom.json <<EOF
{
  "spdxVersion": "SPDX-2.3",
  "dataLicense": "CC0-1.0",
  "name": "${module}",
  "documentNamespace": "https://github.com/akhenakh/gojoltc/sbom",
  "creationInfo": {
    "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "creators": ["Tool: scripts/gen_sbom.sh"]
  },
  "packages": [
    {
      "name": "${module}",
      "versionInfo": "go${gover}",
      "downloadLocation": "NOASSERTION",
      "licenseConcluded": "MIT",
      "primaryPackagePurpose": "APPLICATION"
    },
    {
      "name": "joltc",
      "versionInfo": "${joltc_sha}",
      "downloadLocation": "${joltc_url}",
      "licenseConcluded": "MIT",
      "comment": "C API for JoltPhysics, vendored as a git submodule (third_party/joltc)."
    },
    {
      "name": "JoltPhysics",
      "versionInfo": "${jolt_ver}",
      "downloadLocation": "${jolt_url}",
      "licenseConcluded": "MIT",
      "comment": "Vendored as a git submodule (third_party/JoltPhysics). Linked statically; no runtime shared object."
    }
  ]
}
EOF

echo "wrote sbom.json"
