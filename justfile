# The recipes of modules' CI (R03, R10). They run from the repository's root, in the toolchain stage
# of .devcontainer/Dockerfile, where module-kit's CLI is the binary that stage pins, at an exact tag
# and sum: ci.yml runs them there by spec's tools/toolchain/toolchain.sh, and a person runs them the
# same way in the devcontainer. Every recipe is bash, and its tests are in test/, run by `just test`.
#
# The calls of module-kit follow the form that module-kit's docs/architecture.md, section 3.6, gives
# for check and build at v0.2.3 (R102 spec):
#
#   module-kit check <module>
#   module-kit build <module> <folder>
#
# where build's folder does not exist or is empty, and any other call is a usage error with the
# status 2. Every command exits with 0 when it accepts, 1 when the format refuses its input, with
# the ID of each rule broken, and 2 when it could not do its work. test/modules-check.test.mjs holds
# the calls to that form.

# modules/ holds module folders alone: no other file, nothing named with a dot, no link (R03). Each
# folder is checked by module-kit check and built into a temporary folder by module-kit build, and the
# ID that its package's inventory.json names (format/README.md in spec, section 25) is
# novatrail/<folder>. Every module is checked, not only the changed ones, and a refusal of one does
# not stop the rest: the recipe fails at the end if any failed, whether module-kit refused a module
# or could not do its work, which its own stderr tells apart. Without modules/ there is nothing to
# check yet, and the recipe says so; git holds no empty folder, so that is the tree with no module.
# Before the first module, the module-kit on the PATH must be the one the Dockerfile pins, so that a
# container not rebuilt since the tag was raised does not check with another CLI than CI's.
# Check the tree of modules/ and every module in it with module-kit.
modules-check:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ ! -e modules ] && [ ! -L modules ]; then
        echo "modules/: no modules yet"
        exit 0
    fi
    if [ -L modules ] || [ ! -d modules ]; then
        echo "modules: not a folder; it holds the folders of the modules (R03)" >&2
        exit 1
    fi
    packages=$(mktemp -d)
    trap 'rm -rf "$packages"' EXIT
    shopt -s nullglob dotglob
    entries=(modules/*)
    if [ "${#entries[@]}" -eq 0 ]; then
        echo "modules/: no modules yet"
        exit 0
    fi
    pinned=$({{quote(just_executable())}} _dockerfile-arg .devcontainer/Dockerfile MODULE_KIT_VERSION)
    if ! version=$(module-kit --version | sed -n 1p); then
        echo "modules-check: module-kit cannot be run; the toolchain stage installs it" >&2
        exit 1
    fi
    if [ "$version" != "module-kit ${pinned}" ]; then
        echo "modules-check: the module-kit on the PATH says \"${version}\", and .devcontainer/Dockerfile pins module-kit ${pinned}: rebuild the container" >&2
        exit 1
    fi
    failed=0
    for path in "${entries[@]}"; do
        name=${path#modules/}
        if [[ "$name" == .* ]]; then
            echo "${path}: modules/ holds the folders of the modules alone, and nothing named with a dot (R03)" >&2
            failed=1
            continue
        fi
        if [ -L "$path" ] || [ ! -d "$path" ]; then
            echo "${path}: not a folder; modules/ holds the folders of the modules alone (R03)" >&2
            failed=1
            continue
        fi
        if ! module-kit check "$path"; then
            failed=1
            continue
        fi
        if ! module-kit build "$path" "${packages}/${name}"; then
            failed=1
            continue
        fi
        if ! id=$(jq -er '.id | strings' "${packages}/${name}/inventory.json"); then
            echo "${path}: the package module-kit built names no ID in its inventory.json" >&2
            failed=1
            continue
        fi
        if [ "$id" != "novatrail/${name}" ]; then
            echo "${path}: the module's ID is ${id}, and the module in modules/${name}/ is novatrail/${name} (R03)" >&2
            failed=1
        fi
    done
    exit "$failed"

# The binary of module-kit that the toolchain stage installs is held by its sums in
# .devcontainer/Dockerfile. Here the release's SHA256SUMS and its bundle are read from the release's
# public address, with no token, the signature is checked by spec's release-verify (R77 spec), and
# the sums of both linux binaries are held to the Dockerfile's: the chain of digests then runs from a
# verified file to the binary. CI runs it on every pull request and push to main (R10).
# Verify the signature of the module-kit release that the Dockerfile pins, and its sums.
module-kit-verify:
    #!/usr/bin/env bash
    set -euo pipefail
    tag=$({{quote(just_executable())}} _dockerfile-arg .devcontainer/Dockerfile MODULE_KIT_VERSION)
    if ! [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "module-kit-verify: MODULE_KIT_VERSION is ${tag}, not a tag vX.Y.Z" >&2
        exit 1
    fi
    release=$(mktemp -d)
    trap 'rm -rf "$release"' EXIT
    for file in SHA256SUMS SHA256SUMS.sigstore.json; do
        curl -fsSL --proto "=https" --proto-redir "=https" --retry 3 --retry-all-errors \
            -o "${release}/${file}" "https://github.com/altnovatrail/module-kit/releases/download/${tag}/${file}"
    done
    {{quote(just_executable())}} -f spec/justfile release-verify module-kit "$tag" "${release}/SHA256SUMS"
    {{quote(just_executable())}} _module-kit-sums "${release}/SHA256SUMS" .devcontainer/Dockerfile

# The sums of module-kit_<tag>_linux_amd64 and module-kit_<tag>_linux_arm64 in a release's
# SHA256SUMS, each listed once, against MODULE_KIT_SHA256_AMD64 and MODULE_KIT_SHA256_ARM64, for the
# tag MODULE_KIT_VERSION, of a Dockerfile. Apart from the network, so that its tests need none.
[private]
_module-kit-sums sums dockerfile:
    #!/usr/bin/env bash
    set -euo pipefail
    sums={{quote(sums)}}
    dockerfile={{quote(dockerfile)}}
    tag=$({{quote(just_executable())}} _dockerfile-arg "$dockerfile" MODULE_KIT_VERSION)
    failed=0
    for arch in amd64 arm64; do
        name="module-kit_${tag}_linux_${arch}"
        wanted=$({{quote(just_executable())}} _dockerfile-arg "$dockerfile" "MODULE_KIT_SHA256_${arch^^}")
        if ! [[ "$wanted" =~ ^[0-9a-f]{64}$ ]]; then
            echo "_module-kit-sums: MODULE_KIT_SHA256_${arch^^} is not a sha256 of 64 lowercase hexadecimal digits" >&2
            failed=1
            continue
        fi
        listed=$(awk -v name="$name" '$2 == name { print $1 }' "$sums")
        count=$(grep -c . <<< "$listed" || true)
        if [ "$count" -ne 1 ]; then
            echo "_module-kit-sums: ${sums} lists ${name} not once but ${count} times" >&2
            failed=1
            continue
        fi
        if [ "$listed" != "$wanted" ]; then
            echo "_module-kit-sums: ${name} is ${listed} in SHA256SUMS, and MODULE_KIT_SHA256_${arch^^} is ${wanted}" >&2
            failed=1
            continue
        fi
        echo "module-kit ${tag}: ${name} ${listed}"
    done
    exit "$failed"

# The value of `ARG <name>=<value>` in a Dockerfile, which holds it once.
[private]
_dockerfile-arg dockerfile name:
    #!/usr/bin/env bash
    set -euo pipefail
    dockerfile={{quote(dockerfile)}}
    name={{quote(name)}}
    if [ ! -f "$dockerfile" ] || [ ! -r "$dockerfile" ]; then
        echo "_dockerfile-arg: ${dockerfile} is no file that can be read" >&2
        exit 1
    fi
    count=$(grep -c "^ARG ${name}=" "$dockerfile" || true)
    if [ "$count" -ne 1 ]; then
        echo "_dockerfile-arg: ${dockerfile} sets ARG ${name} not once but ${count} times" >&2
        exit 1
    fi
    sed -n "s/^ARG ${name}=//p" "$dockerfile"

# The tests run this same just, and the module-kit on the PATH: the toolchain stage's.
# Run the tests of the recipes.
test:
    JUST={{quote(just_executable())}} node --test 'test/*.test.mjs'
