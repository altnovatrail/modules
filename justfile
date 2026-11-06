# The recipes of modules' CI (R03, R10, R11). They run from the repository's root, in the toolchain
# stage of .devcontainer/Dockerfile, where module-kit's CLI is the binary that stage pins, at an
# exact tag and sum: ci.yml runs them there by spec's tools/toolchain/toolchain.sh, and a person runs
# them the same way in the devcontainer. Every recipe is bash, and its tests are in test/, run by
# `just test`.
#
# The calls of module-kit follow the forms that module-kit's docs/architecture.md gives at v0.2.3
# (R102 spec) — check and build in section 3.6, diff in section 3.5:
#
#   module-kit check <module>
#   module-kit build <module> <folder>
#   module-kit diff [--json] (--index <module-index.json> | --first) [--published <package folder>]... <module>
#
# where build's folder does not exist or is empty; diff's <module> is a module's folder or a package
# that build wrote, its --index an index that release-verify has accepted, and each --published a
# package as files, without device/, which diff finds by the digest of its inventory.json and checks
# by the chain of digests. Any other call is a usage error with the status 2. Every command exits
# with 0 when it accepts, 1 when the format refuses its input, with the ID of each rule broken, and 2
# when it could not do its work. test/modules-check.test.mjs holds the calls to those forms.

# The published versions are those of the index of the newest release of modules, whose signature
# release-verify checks, and the package of each version of every module in modules/, withdrawn
# ones too, comes from the release its entry names (R03; R11). Before the first release there is no
# index, and every version is a first one: _published says so by a file of its own, and nothing
# skips the comparison unsaid. With no module there is nothing to compare, and the network is not
# asked. Any failure to read the published versions gives the status 2: the modules were not
# compared, not refused.
# Check the tree of modules/ and every module in it with module-kit, against its published versions.
modules-check:
    #!/usr/bin/env bash
    set -euo pipefail
    shopt -s nullglob dotglob
    entries=(modules/*)
    if { [ ! -e modules ] && [ ! -L modules ]; } || { [ -d modules ] && [ ! -L modules ] && [ "${#entries[@]}" -eq 0 ]; }; then
        echo "modules/: no modules yet"
        exit 0
    fi
    published=$(mktemp -d)
    trap 'rm -rf "$published"' EXIT
    if ! {{quote(just_executable())}} _published "${published}/published"; then
        echo "modules-check: the published versions of the modules cannot be read, so nothing is compared" >&2
        exit 2
    fi
    {{quote(just_executable())}} _modules-check "${published}/published"

# What modules-check compares with, written into a new folder: either module-index.json of the newest
# release of modules, accepted by release-verify, and packages/sha256-<hex>/ for the package of each
# version its entries name for the modules of modules/; or the file first, when modules has no
# release yet. The releases are the tags vX.Y.Z that git ls-remote lists, read with no token and none
# of the caller's git settings; tags v* only a release sets, and a ruleset keeps them (R31 spec;
# R11). The newest of them whose module-index.json is there is the newest release: a tag whose
# release did not come out — its index answers 404 — is named in the log and passed over, since such
# a tag can be neither moved nor deleted, and its versions were never published. Any other failure
# to read a file fails the recipe, rather than pass for no release. Each archive comes from the
# public address the entry gives (architecture 3.3), and _unpack takes only the three files of the
# package out of it.
[private]
_published dir:
    #!/usr/bin/env bash
    set -euo pipefail
    dir={{quote(dir)}}
    if [ -e "$dir" ] || [ -L "$dir" ]; then
        echo "_published: ${dir} is there already; the published versions go into a new folder" >&2
        exit 2
    fi
    mkdir "$dir"
    lines=$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 GIT_TERMINAL_PROMPT=0 \
        git ls-remote --tags --refs https://github.com/altnovatrail/modules 'refs/tags/v*')
    releases=https://github.com/altnovatrail/modules/releases/download
    tag=""
    while read -r candidate; do
        if [ -z "$candidate" ]; then continue; fi
        status=0
        {{quote(just_executable())}} _fetch "${releases}/${candidate}/module-index.json" "${dir}/module-index.json" || status=$?
        if [ "$status" -eq 0 ]; then
            tag=$candidate
            break
        fi
        if [ "$status" -eq 1 ]; then
            echo "_published: the tag ${candidate} has no module-index.json: its release did not come out, and the tag below it is taken"
            continue
        fi
        exit 2
    done <<< "$({{quote(just_executable())}} _release-tags <<< "$lines")"
    if [ -z "$tag" ]; then
        echo "modules has no release yet, so every version is a first one" > "${dir}/first"
        echo "_published: modules has no release yet; every version is a first one"
        exit 0
    fi
    {{quote(just_executable())}} _fetch "${releases}/${tag}/module-index.json.sigstore.json" "${dir}/module-index.json.sigstore.json"
    {{quote(just_executable())}} -f spec/justfile release-verify modules "$tag" "${dir}/module-index.json"
    rm "${dir}/module-index.json.sigstore.json"
    mkdir "${dir}/archives" "${dir}/packages"
    shopt -s nullglob
    for path in modules/*/; do
        name=$(basename "$path")
        if [ -L "${path%/}" ] || [[ "$name" == .* ]]; then
            continue
        fi
        entries=$({{quote(just_executable())}} _entries "${dir}/module-index.json" "novatrail/${name}")
        while read -r version release digest; do
            if [ -z "$version" ]; then continue; fi
            archive="${digest/:/-}.tar.gz"
            {{quote(just_executable())}} _fetch "${releases}/${release}/${archive}" "${dir}/archives/${archive}"
            {{quote(just_executable())}} _unpack "${dir}/archives/${archive}" "${dir}/packages/${digest/:/-}"
            echo "_published: novatrail/${name} ${version}, ${digest}, from ${release}"
        done <<< "$entries"
    done
    rm -rf "${dir}/archives"

# The entries of a module in an index, a line each: its version, the release that holds its package
# and the package's digest, withdrawn versions too. A release that is no tag vX.Y.Z or a digest that
# is no sha256 gives no address, and is refused, as is an index whose modules is no object.
[private]
_entries index id:
    #!/usr/bin/env bash
    set -euo pipefail
    index={{quote(index)}}
    id={{quote(id)}}
    if ! jq -e '.modules | type == "object"' "$index" > /dev/null; then
        echo "_entries: ${index} holds no object modules, so it is no index" >&2
        exit 2
    fi
    entries=$(jq -r --arg id "$id" \
        '.modules[$id] // {} | to_entries[] | "\(.key) \(.value.release) \(.value.package_digest)"' "$index")
    while read -r version release digest; do
        if [ -z "$version" ]; then continue; fi
        if ! [[ "$release" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ && "$digest" =~ ^sha256:[0-9a-f]{64}$ ]]; then
            echo "_entries: ${id} ${version} names the release ${release} and the package ${digest}, which give no address" >&2
            exit 2
        fi
        echo "${version} ${release} ${digest}"
    done <<< "$entries"

# The tags vX.Y.Z among the lines of git ls-remote on stdin, newest first by the numbers of their
# parts, or nothing when there is none. Other tags are not releases of modules and are left out.
[private]
_release-tags:
    #!/usr/bin/env bash
    set -euo pipefail
    sed -n 's|^[0-9a-f]*\trefs/tags/||p' \
        | { grep -E '^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$' || true; } \
        | sort -t . -k 1.2,1nr -k 2,2nr -k 3,3nr

# A file of a release, from its public address, with no token: https alone, redirects too, and every
# transfer retried, as the Dockerfile's fetch does. The status is 0 when the file came, 1 when the
# address answers 404 — the release holds no such file —, and 2 for any other failure.
[private]
_fetch url out:
    #!/usr/bin/env bash
    set -euo pipefail
    url={{quote(url)}}
    out={{quote(out)}}
    code=$(curl -sSL --proto "=https" --proto-redir "=https" --retry 3 --retry-all-errors \
        -w '%{http_code}' -o "$out" "$url") || {
        echo "_fetch: ${url} cannot be read" >&2
        rm -f "$out"
        exit 2
    }
    case "$code" in
        200) ;;
        404)
            rm -f "$out"
            echo "_fetch: ${url} answers 404" >&2
            exit 1
            ;;
        *)
            rm -f "$out"
            echo "_fetch: ${url} answers ${code}" >&2
            exit 2
            ;;
    esac

# The package of a release's archive as diff reads it, in a new folder: inventory.json, server.json
# and device.json, each once and a regular file, with or without ./ before its name, and nothing
# else — no device/, which diff does not read (architecture 3.5 module-kit). An archive that holds a
# link, a special file, an absolute path or a path through .. is refused whole, before anything
# comes out of it. The digests of the three files are diff's to check: it finds a package by the
# digest of its inventory.json, the one its entry names, and refuses a broken chain.
[private]
_unpack archive dir:
    #!/usr/bin/env bash
    set -euo pipefail
    archive={{quote(archive)}}
    dir={{quote(dir)}}
    if [ -e "$dir" ] || [ -L "$dir" ]; then
        echo "_unpack: ${dir} is there already; a package goes into a new folder" >&2
        exit 2
    fi
    if ! names=$(tar --quoting-style=escape -tzf "$archive") \
        || ! listing=$(tar --quoting-style=escape -tvzf "$archive"); then
        echo "_unpack: ${archive} is no gzip tar archive that can be read" >&2
        exit 2
    fi
    if [ -z "$names" ]; then
        names=()
        types=()
    else
        mapfile -t names <<< "$names"
        mapfile -t types < <(cut -c1 <<< "$listing")
    fi
    if [ "${#names[@]}" -ne "${#types[@]}" ]; then
        echo "_unpack: ${archive}: the listing of its entries does not add up" >&2
        exit 2
    fi
    declare -A member=()
    for i in "${!names[@]}"; do
        name=${names[$i]}
        case "${types[$i]}" in
            -|d) ;;
            *)
                echo "_unpack: ${archive} holds ${name}, a link or a special file" >&2
                exit 2
                ;;
        esac
        if [[ "$name" == /* || "/${name}/" == */../* ]]; then
            echo "_unpack: ${archive} holds ${name}, a path outside its folder" >&2
            exit 2
        fi
        file=${name#./}
        case "$file" in
            inventory.json|server.json|device.json)
                if [ -n "${member[$file]:-}" ]; then
                    echo "_unpack: ${archive} holds ${file} twice" >&2
                    exit 2
                fi
                member[$file]=$name
                ;;
        esac
    done
    for file in inventory.json server.json device.json; do
        if [ -z "${member[$file]:-}" ]; then
            echo "_unpack: ${archive} holds no ${file} at its root" >&2
            exit 2
        fi
    done
    mkdir -p "$dir"
    tar -xzf "$archive" -C "$dir" --no-same-owner --no-same-permissions -- \
        "${member[inventory.json]}" "${member[server.json]}" "${member[device.json]}"
    for file in inventory.json server.json device.json; do
        if [ -L "${dir}/${file}" ] || [ ! -f "${dir}/${file}" ]; then
            echo "_unpack: ${file} of ${archive} is no regular file" >&2
            exit 2
        fi
    done

# modules/ holds module folders alone: no other file, nothing named with a dot, no link (R03). Each
# folder is checked by module-kit check and built into a temporary folder by module-kit build, the
# ID that its package's inventory.json names (format/README.md in spec, section 25) is
# novatrail/<folder>, and module-kit diff compares that package with every published version that
# the folder published, a folder _published wrote, names: --index and the packages of the module's
# entries, or --first when it holds the file first, and exactly one of the two. Every module is
# checked, not only the changed ones, and a refusal of one does not stop the rest: the recipe exits
# at the end with the highest status of what failed, 1 for the tree, a refusal and a breach of R03,
# and 2 when module-kit could not do its work, so that diff's status is the job's (R11). Without
# modules/ there is nothing to check yet, and the recipe says so; git holds no empty folder, so that
# is the tree with no module. Before the first module, the module-kit on the PATH must be the one the
# Dockerfile pins, so that a container not rebuilt since the tag was raised does not check with
# another CLI than CI's.
[private]
_modules-check published:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ ! -d {{quote(published)}} ]; then
        echo "_modules-check: "{{quote(published)}}" is no folder of the published versions" >&2
        exit 2
    fi
    published=$(cd {{quote(published)}} && pwd)
    if [ -f "${published}/first" ] && [ ! -e "${published}/module-index.json" ]; then
        compare=(--first)
    elif [ -f "${published}/module-index.json" ] && [ ! -e "${published}/first" ]; then
        compare=(--index "${published}/module-index.json")
    else
        echo "_modules-check: ${published} holds neither module-index.json alone nor first alone: what to compare with is unsaid" >&2
        exit 2
    fi
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
        exit 2
    fi
    if [ "$version" != "module-kit ${pinned}" ]; then
        echo "modules-check: the module-kit on the PATH says \"${version}\", and .devcontainer/Dockerfile pins module-kit ${pinned}: rebuild the container" >&2
        exit 1
    fi
    worst=0
    fail() {
        if [ "$1" -gt "$worst" ]; then worst=$1; fi
    }
    for path in "${entries[@]}"; do
        name=${path#modules/}
        if [[ "$name" == .* ]]; then
            echo "${path}: modules/ holds the folders of the modules alone, and nothing named with a dot (R03)" >&2
            fail 1
            continue
        fi
        if [ -L "$path" ] || [ ! -d "$path" ]; then
            echo "${path}: not a folder; modules/ holds the folders of the modules alone (R03)" >&2
            fail 1
            continue
        fi
        status=0
        module-kit check "$path" || status=$?
        if [ "$status" -ne 0 ]; then
            fail "$status"
            continue
        fi
        module-kit build "$path" "${packages}/${name}" || status=$?
        if [ "$status" -ne 0 ]; then
            fail "$status"
            continue
        fi
        if ! id=$(jq -er '.id | strings' "${packages}/${name}/inventory.json"); then
            echo "${path}: the package module-kit built names no ID in its inventory.json" >&2
            fail 2
            continue
        fi
        if [ "$id" != "novatrail/${name}" ]; then
            echo "${path}: the module's ID is ${id}, and the module in modules/${name}/ is novatrail/${name} (R03)" >&2
            fail 1
            continue
        fi
        args=("${compare[@]}")
        if [ "${compare[0]}" = --index ]; then
            if ! entries=$({{quote(just_executable())}} _entries "${compare[1]}" "$id"); then
                exit 2
            fi
            while read -r _ _ digest; do
                folder="${published}/packages/${digest/:/-}"
                if [ -n "$digest" ] && [ -d "$folder" ]; then
                    args+=(--published "$folder")
                fi
            done <<< "$entries"
        fi
        module-kit diff "${args[@]}" "${packages}/${name}" || fail $?
    done
    exit "$worst"

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
        {{quote(just_executable())}} _fetch "https://github.com/altnovatrail/module-kit/releases/download/${tag}/${file}" "${release}/${file}"
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
