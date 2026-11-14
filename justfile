# The recipes of modules' CI (R03, R10, R11). They run from the repository's root, in the toolchain
# stage of .devcontainer/Dockerfile, where module-kit's CLI is the binary that stage pins, at an
# exact tag and sum: ci.yml runs them there by spec's tools/toolchain/toolchain.sh, and a person runs
# them the same way in the devcontainer. Every recipe is bash, and its tests are in test/, run by
# `just test`, with two tests of content among them: test/templates.test.mjs builds olympiad-math
# and checks each template against the task it came from (R22), and test/lesson-pictures.test.mjs
# checks the size of each picture of a lesson against its file (amendment T41 to R12).
#
# The calls of module-kit follow the forms that module-kit's docs/architecture.md gives at v0.3.0
# (R102 spec) — check, build and index in section 3.6, diff in section 3.5:
#
#   module-kit check <module>
#   module-kit build <module> <folder>
#   module-kit index (--previous <module-index.json> | --first) --tag vX.Y.Z
#       (--withdrawals <withdrawals.yaml> | --no-withdrawals) --out <file> <package folder>...
#   module-kit diff [--json] (--index <module-index.json> | --first) [--published <package folder>]... <module>
#
# where build's folder does not exist or is empty; index's --out does not exist yet, its flags come
# before its packages, folders that build wrote, with device/, and of each pair exactly one is given;
# diff's <module> is a module's folder or a package that build wrote, its --index an index that
# release-verify has accepted, and each --published a package as files, without device/, which diff
# finds by the digest of its inventory.json and checks by the chain of digests. Any other call is a
# usage error with the status 2. Every command exits with 0 when it accepts, 1 when the format
# refuses its input, with the ID of each rule broken, and 2 when it could not do its work.
# test/modules-check.test.mjs and test/published.test.mjs hold the calls of check, build and diff to
# those forms, and test/release-files.test.mjs the call of index.

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
# release yet. The newest release and its index are _newest-release's to find. Each archive comes
# from the public address the entry gives (architecture 3.3), and _unpack takes only the three files
# of the package out of it.
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
    releases=https://github.com/altnovatrail/modules/releases/download
    tag=$({{quote(just_executable())}} _newest-release "${dir}/module-index.json")
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

# The tag of the newest release of modules, printed, with its module-index.json written into the
# file index, which must not be there yet; nothing printed and no file when modules has no release
# yet. The releases are the tags vX.Y.Z that git ls-remote lists, read with no token and none of the
# caller's git settings; tags v* only a release sets, and a ruleset keeps them (R31 spec; R11). The
# newest of them whose module-index.json is there is the newest release: a tag whose release did not
# come out — its index answers 404 — is named on stderr and passed over, since such a tag can be
# neither moved nor deleted, and its versions were never published. Any other failure to read a
# file gives the status 2, rather than pass for no release. The index is not verified here: whoever
# reads more of it than that it is there runs release-verify on it, as _published does.
[private]
_newest-release index:
    #!/usr/bin/env bash
    set -euo pipefail
    index={{quote(index)}}
    if [ -e "$index" ] || [ -L "$index" ]; then
        echo "_newest-release: ${index} is there already; the index goes into a new file" >&2
        exit 2
    fi
    if ! lines=$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 GIT_TERMINAL_PROMPT=0 \
        git ls-remote --tags --refs https://github.com/altnovatrail/modules 'refs/tags/v*'); then
        echo "_newest-release: the tags of modules cannot be read" >&2
        exit 2
    fi
    releases=https://github.com/altnovatrail/modules/releases/download
    while read -r candidate; do
        if [ -z "$candidate" ]; then continue; fi
        status=0
        {{quote(just_executable())}} _fetch "${releases}/${candidate}/module-index.json" "$index" || status=$?
        if [ "$status" -eq 0 ]; then
            echo "$candidate"
            exit 0
        fi
        if [ "$status" -eq 1 ]; then
            echo "_newest-release: the tag ${candidate} has no module-index.json: its release did not come out, and the tag below it is taken" >&2
            continue
        fi
        exit 2
    done <<< "$({{quote(just_executable())}} _release-tags <<< "$lines")"

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

# A module's eval set is its folder eval/, which its package does not carry and which modules
# publishes with the repository, for good (R04; R62 spec). Its layout is flat: each sample is the
# file eval/<id>.json, whose name is the sample's ID — the field id of its one JSON object, 1 to 64
# lowercase Latin letters, digits and hyphens (contracts/README.md in spec, section 11) —, and eval/
# holds nothing else: no folder, no link, no other file, and nothing named with a dot at any depth,
# which the eval set would leave out and the repository would still publish unchecked (amendments
# 2026-11-06 and 2026-11-09 to R04). jq reads the ID alone, and holds it to the name itself, so that
# no byte of it is lost on the way. The rules of the sample file, and those of an eval set that need
# neither the module's package nor a key, are spec's: sample-check --eval-set novatrail/<folder> on
# the module's eval set as section 13 there has it — every regular file under eval/ whose name
# ends in .json, at any depth, but under a name with a dot —, with no link, which the layout
# refuses and which would stop sample-check. A module with no such file is not checked, since a
# check of no file passes none (R93 spec; amendment 2026-11-09 to R25). sample-check gets all of a
# module's files in one call, in their order: since T116 spec it holds whole only the file it reads,
# so that one call checks hundreds of large samples with no batches (contracts/README.md in spec,
# section 11), and it holds the one rule across files, that no two samples share an ID, itself, as
# the layout does too (amendment R32 to R25). Only the folders of
# modules/ are read: a file, a link or a name with a dot there is modules-check's to refuse (R03).
# Every module is checked, and a breach in one does not stop the rest; the status is the highest
# of what failed: 1 for a breach of the layout, a sample that sample-check refuses and modules that
# is no folder, and 2 when jq, spec's tools or sample-check could not do their work. A folder whose
# name gives no module ID stops sample-check with 2, and modules-check refuses it in the same job.
# It needs no key and writes nothing in the tree; before the first eval set, spec's _tools installs
# spec/tools from its lock file, as docs-check-repo does — the only write, and the only use of the
# network, which npm's registry gets —, so a pull request from a fork runs it too.
# Check the layout of each module's eval/, and each module's eval set with spec's sample-check.
eval-check:
    #!/usr/bin/env bash
    set -euo pipefail
    export LC_ALL=C
    shopt -s nullglob dotglob
    if [ ! -e modules ] && [ ! -L modules ]; then
        echo "modules/: no modules yet"
        exit 0
    fi
    if [ -L modules ] || [ ! -d modules ]; then
        echo "modules: not a folder; it holds the folders of the modules (R03)" >&2
        exit 1
    fi
    entries=(modules/*)
    if [ "${#entries[@]}" -eq 0 ]; then
        echo "modules/: no modules yet"
        exit 0
    fi
    worst=0
    fail() {
        if [ "$1" -gt "$worst" ]; then worst=$1; fi
    }
    breach() {
        echo "$1" >&2
        fail 1
    }
    jq_runs=""
    tools=""
    held=""
    for path in "${entries[@]}"; do
        module=${path#modules/}
        if [[ "$module" == .* ]] || [ -L "$path" ] || [ ! -d "$path" ]; then
            continue
        fi
        eval="${path}/eval"
        if [ ! -e "$eval" ] && [ ! -L "$eval" ]; then
            continue
        fi
        if [ -L "$eval" ] || [ ! -d "$eval" ]; then
            breach "${eval}: not a folder; eval/ is the folder of the module's eval set (R04)"
            continue
        fi
        while IFS= read -r -d '' file; do
            breach "${file}: named with a dot; eval/ holds the files eval/<id>.json alone, and nothing named with a dot (R04)"
        done < <(find "$eval" -mindepth 1 -name '.*' -print0 | sort -z)
        for file in "$eval"/*; do
            name=${file##*/}
            if [[ "$name" == .* ]]; then
                continue
            fi
            if [ -L "$file" ] || [ ! -f "$file" ]; then
                breach "${file}: a folder, a link or a special file; eval/ holds the files eval/<id>.json alone (R04)"
                continue
            fi
            if ! [[ "$name" =~ ^[a-z0-9-]{1,64}\.json$ ]]; then
                breach "${file}: not eval/<id>.json, whose <id> is 1 to 64 lowercase Latin letters, digits and hyphens (R04)"
                continue
            fi
            if [ -z "$jq_runs" ]; then
                if jq --version > /dev/null 2>&1; then
                    jq_runs=yes
                else
                    echo "eval-check: jq cannot be run; the toolchain stage has it" >&2
                    jq_runs=no
                fi
            fi
            if [ "$jq_runs" = no ]; then
                fail 2
                continue
            fi
            # The ID comes out as JSON, quotes and escapes and all, and is compared in jq, so that a
            # newline or a NUL at its end, which the shell would drop, still makes it another ID.
            if ! id=$(jq -e -r -s 'if length == 1 and (.[0] | type) == "object" and (.[0].id | type) == "string" then .[0].id | @json else empty end' "$file" 2> /dev/null); then
                breach "${file}: holds no one JSON object with a string id, so its name cannot be held to its sample's ID (R04)"
            elif ! jq -e -s --arg name "${name%.json}" '.[0].id == $name' "$file" > /dev/null 2>&1; then
                breach "${file}: its sample's ID is ${id}, and a sample lies in eval/<its ID>.json (R04)"
            fi
        done
        files=()
        while IFS= read -r -d '' file; do
            files+=("$file")
        done < <(find "$eval" -mindepth 1 -name '.*' -prune -o -name '*.json' -type f -print0 | sort -z)
        if [ "${#files[@]}" -eq 0 ]; then
            continue
        fi
        held=yes
        # spec's tools are installed first, by sample-check's own dependency _tools, so that npm,
        # which fails with 1 as a refusal does, cannot pass for one: after it, sample-check gives 1
        # when a file breaks a rule and 2 when it cannot check, and any status above 2 is that too.
        if [ -z "$tools" ]; then
            if {{quote(just_executable())}} -f spec/justfile _tools; then
                tools=yes
            else
                echo "eval-check: the tools of spec's sample-check cannot be installed, so no eval set is checked" >&2
                tools=no
            fi
        fi
        if [ "$tools" = no ]; then
            fail 2
            continue
        fi
        status=0
        {{quote(just_executable())}} -f spec/justfile sample-check --eval-set "novatrail/${module}" "${files[@]}" || status=$?
        if [ "$status" -gt 2 ]; then status=2; fi
        fail "$status"
    done
    if [ "$worst" -eq 0 ] && [ -z "$held" ]; then
        echo "eval-check: every eval/ is laid out as eval/<id>.json, and no module holds a sample"
    elif [ "$worst" -eq 0 ]; then
        echo "eval-check: every eval/ is laid out as eval/<id>.json, and sample-check holds the eval set of every module that has one"
    fi
    exit "$worst"

# Where the toolchain stage installs canvas's generator, synthetic-module_<CANVAS_VERSION>.mjs; the
# tests of the recipes put a stand-in elsewhere by CANVAS_LIB.
canvas_lib := env("CANVAS_LIB", "/usr/local/lib/canvas")

# A synthetic sample of a module's eval set is one of a task of that module, made on the package
# of the commit it lies in (R110 spec; R29). The samples name the tag of the newest release of
# modules, which _synthetic takes from _newest-release once nothing local stands in the way. It is a
# step a person runs in the devcontainer, which writes into the module's eval/ and nothing else.
# Write synthetic samples of a module's eval set, <task>:<seed> each, with canvas's generator.
[positional-arguments]
synthetic module *calls:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ "$#" -lt 2 ]; then
        echo "synthetic: name the module and the samples to write, <task>:<seed> each" >&2
        exit 2
    fi
    module=$1
    shift
    {{quote(just_executable())}} _synthetic "$module" "" "$@"

# The synthetic samples <task>:<seed> of modules/<module>/, written into its eval/ by canvas's
# generator, on the package that module-kit build makes of the module into a temporary folder, for
# the release tag, or, when the tag is empty, the newest release's, which _newest-release finds
# after every check that needs no network (R110 spec; R29). The generator is the file of the release of canvas that the
# Dockerfile pins, CANVAS_VERSION, held to CANVAS_SYNTHETIC_SHA256, as the module-kit on the PATH is
# held to MODULE_KIT_VERSION, so that a container not rebuilt since a tag was raised writes no
# sample CI would not make. The calls go to the generator as they are, in its form (R68 canvas):
#
#   node synthetic-module_<tag>.mjs --module <package> <tag> <folder> <task>:<seed>...
#
# It checks every call before it writes any file, and writes nothing when one fails; it never
# writes over a person's sample. Its lines, a sample's ID, task and what it shows, and its reasons
# for each call refused are the recipe's, and so is its status: 0 when it wrote them all, 1 when it
# refused a call — a template, a task it does not draw, a seed named twice — and 2 when the call was
# out of form. The recipe's own status is 1 when the package's ID is not novatrail/<module> (R03),
# module-kit's when it refuses the module, and 2 when it cannot do its work.
[private]
[positional-arguments]
_synthetic module tag *calls:
    #!/usr/bin/env bash
    set -euo pipefail
    module=$1 tag=$2
    shift 2
    if [ "$#" -eq 0 ]; then
        echo "_synthetic: name the samples to write, <task>:<seed> each" >&2
        exit 2
    fi
    path="modules/${module}"
    if [ -z "$module" ] || [[ "$module" == */* || "$module" == .* ]] || [ -L "$path" ] || [ ! -d "$path" ]; then
        echo "_synthetic: ${path} is no folder of a module in modules/ (R03)" >&2
        exit 2
    fi
    # The generator writes into eval/ and makes it when it is not there; a link or a file there would
    # take the samples elsewhere, which eval-check would refuse only once they were written.
    if [ -L "${path}/eval" ] || { [ -e "${path}/eval" ] && [ ! -d "${path}/eval" ]; }; then
        echo "_synthetic: ${path}/eval is a link or no folder; eval/ is the folder of the module's eval set (R04)" >&2
        exit 2
    fi
    if [ -n "$tag" ] && ! [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "_synthetic: ${tag} is no tag vX.Y.Z of a release of modules" >&2
        exit 2
    fi
    {{quote(just_executable())}} _pinned-module-kit
    generator=$({{quote(just_executable())}} _generator)
    package=$(mktemp -d)
    trap 'rm -rf "$package"' EXIT
    if [ -z "$tag" ]; then
        tag=$({{quote(just_executable())}} _newest-release "${package}/module-index.json")
        if [ -z "$tag" ]; then
            echo "_synthetic: modules has no release yet, and a synthetic sample names the release its package is of" >&2
            exit 2
        fi
    fi
    status=0
    module-kit build "$path" "${package}/package" || status=$?
    if [ "$status" -ne 0 ]; then
        exit "$status"
    fi
    if ! id=$(jq -er '.id | strings' "${package}/package/inventory.json"); then
        echo "_synthetic: the package module-kit built names no ID in its inventory.json" >&2
        exit 2
    fi
    if [ "$id" != "novatrail/${module}" ]; then
        echo "_synthetic: the module's ID is ${id}, and the module in modules/${module}/ is novatrail/${module} (R03)" >&2
        exit 1
    fi
    node "$generator" --module "${package}/package" "$tag" "${path}/eval" "$@"

# The module-kit on the PATH is the one .devcontainer/Dockerfile pins, so that a container not
# rebuilt since the tag was raised neither writes nor checks a sample with another CLI than CI's;
# the status is 2 otherwise.
[private]
_pinned-module-kit:
    #!/usr/bin/env bash
    set -euo pipefail
    pinned=$({{quote(just_executable())}} _dockerfile-arg .devcontainer/Dockerfile MODULE_KIT_VERSION)
    if ! version=$(module-kit --version | sed -n 1p); then
        echo "_pinned-module-kit: module-kit cannot be run; the toolchain stage installs it" >&2
        exit 2
    fi
    if [ "$version" != "module-kit ${pinned}" ]; then
        echo "_pinned-module-kit: the module-kit on the PATH says \"${version}\", and .devcontainer/Dockerfile pins module-kit ${pinned}: rebuild the container" >&2
        exit 2
    fi

# The path of canvas's generator, printed once it is the file the Dockerfile pins: the stage's
# synthetic-module_<CANVAS_VERSION>.mjs, a regular file whose sum is CANVAS_SYNTHETIC_SHA256 (R29).
# The status is 2 otherwise: a container not rebuilt since the tag was raised would draw samples
# CI does not.
[private]
_generator:
    #!/usr/bin/env bash
    set -euo pipefail
    canvas=$({{quote(just_executable())}} _dockerfile-arg .devcontainer/Dockerfile CANVAS_VERSION)
    wanted=$({{quote(just_executable())}} _dockerfile-arg .devcontainer/Dockerfile CANVAS_SYNTHETIC_SHA256)
    generator={{quote(canvas_lib)}}"/synthetic-module_${canvas}.mjs"
    if [ -L "$generator" ] || [ ! -f "$generator" ]; then
        echo "_generator: ${generator} is not there; the toolchain stage installs canvas's generator ${canvas}: rebuild the container" >&2
        exit 2
    fi
    sum=$(sha256sum < "$generator" | cut -d ' ' -f 1)
    if [ "$sum" != "$wanted" ]; then
        echo "_generator: ${generator} is ${sum}, and .devcontainer/Dockerfile pins ${wanted}: rebuild the container" >&2
        exit 2
    fi
    echo "$generator"

# Every sample of a module's eval set is one of the package of the commit it lies in: a pull request
# that changes the open view of a task re-makes its synthetic samples and removes the samples of
# people in the same pull request (amendment 2026-11-06 to R04; contracts/README.md in spec, section
# 13; amendment T59 to R25). The eval set is the files eval-check gives sample-check, and of each
# jq reads its origin, the source and ID of its task, its seed and the digest of its task's open
# view, and nothing else: a person's sample is never printed. A module with no sample is not built;
# the others are, by the pinned module-kit, and their samples are held to the package by
# _match-synthetic and _match-people. A file outside the layout eval/<id>.json, of which those
# fields cannot be read, or that names another module, is refused here too, as eval-check refuses
# it, so that no sample passes unchecked. Every module is checked, a refusal in one does not stop
# the rest, and the status is the highest of what failed: 1 for a sample that is not one of the
# package, 2 for what could not do its work. No network and no key: a pull request from a fork runs
# it too. The rows between the recipes are a line a sample, its fields split by U+001F, which no
# field holds: a tab would let read merge an empty seed with the next field.
# Hold each sample of every module's eval set to the package of this commit.
eval-match:
    #!/usr/bin/env bash
    set -euo pipefail
    export LC_ALL=C
    shopt -s nullglob dotglob
    if [ -L modules ] || [ ! -d modules ]; then
        echo "eval-match: no modules/, so no sample"
        exit 0
    fi
    work=$(mktemp -d)
    trap 'rm -rf "$work"' EXIT
    us=$'\x1f'
    worst=0
    fail() {
        if [ "$1" -gt "$worst" ]; then worst=$1; fi
    }
    breach() {
        echo "$1" >&2
        fail 1
    }
    pinned=""
    jq_runs=""
    held=""
    for path in modules/*; do
        module=${path#modules/}
        if [[ "$module" == .* ]] || [ -L "$path" ] || [ ! -d "$path" ] || [ -L "${path}/eval" ] || [ ! -d "${path}/eval" ]; then
            continue
        fi
        files=()
        while IFS= read -r -d '' file; do
            files+=("$file")
        done < <(find "${path}/eval" -mindepth 1 -name '.*' -prune -o -name '*.json' -type f -print0 | sort -z)
        if [ "${#files[@]}" -eq 0 ]; then
            continue
        fi
        held=yes
        if [ -z "$jq_runs" ]; then
            if jq --version > /dev/null 2>&1; then jq_runs=yes; else jq_runs=no; fi
        fi
        if [ "$jq_runs" = no ]; then
            echo "eval-match: jq cannot be run; the toolchain stage has it" >&2
            fail 2
            continue
        fi
        rows="${work}/${module}"
        mkdir "$rows"
        : > "${rows}/synthetic"
        : > "${rows}/people"
        for file in "${files[@]}"; do
            if ! [[ "$file" =~ ^modules/[^/]+/eval/[a-z0-9-]{1,64}\.json$ ]]; then
                breach "${file}: outside the layout eval/<id>.json, which eval-check refuses, so it is not held to the package"
                continue
            fi
            if ! row=$(jq -e -r -s '
                if length == 1 and (.[0] | type) == "object" and (.[0].task | type) == "object" and (.[0].task.modules | type) == "object" then
                    .[0] as $s | [
                        (if ($s.synthetic | type) == "object" then "synthetic" elif ($s.person | type) == "object" then "person" else "none" end),
                        $s.task.modules.module, $s.task.task,
                        (if ($s.synthetic | type) == "object" then $s.synthetic.seed else $s.task.seed end),
                        $s.task.modules.tag, $s.task.open_view_digest
                    ] | map(if type == "string" and (test("[\u0000-\u001f]") | not) then . else "" end) | join("\u001f")
                else empty end' "$file" 2> /dev/null); then
                breach "${file}: holds no sample whose task and origin can be read, so it cannot be held to the package"
                continue
            fi
            IFS=$us read -r origin of task seed tag digest <<< "$row"
            if [ "$of" != "novatrail/${module}" ]; then
                breach "${file}: a sample of the eval set of novatrail/${module} is of a task of that module"
                continue
            fi
            if ! [[ "$task" =~ ^[a-z0-9][a-z0-9-]*$ && "$seed" =~ ^(0|[1-9][0-9]*)?$ ]]; then
                breach "${file}: names no task ID and seed that can be held to the package"
                continue
            fi
            case "$origin" in
                synthetic)
                    if [ -z "$seed" ] || ! [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
                        breach "${file}: a synthetic sample with no seed, or no tag vX.Y.Z of modules, cannot be drawn again"
                        continue
                    fi
                    printf '%s\n' "${file}${us}${task}${us}${seed}${us}${tag}" >> "${rows}/synthetic"
                    ;;
                person)
                    # A Claude Code session may read eval/synthetic-*.json alone (.claude/settings.json),
                    # so a person's sample never takes that name (amendment T61 to R04).
                    if [[ "${file##*/}" == synthetic-* ]]; then
                        breach "${file}: the ID of a person's sample begins with synthetic-, which only a synthetic sample's does, and such a file a Claude Code session may read (amendment T61 to R04)"
                        continue
                    fi
                    printf '%s\n' "${file}${us}${task}${us}${seed}${us}${digest}" >> "${rows}/people"
                    ;;
                *)
                    breach "${file}: neither a person's sample nor a synthetic one, so it cannot be held to the package"
                    ;;
            esac
        done
        if [ ! -s "${rows}/synthetic" ] && [ ! -s "${rows}/people" ]; then
            continue
        fi
        if [ -z "$pinned" ]; then
            if {{quote(just_executable())}} _pinned-module-kit; then pinned=yes; else pinned=no; fi
        fi
        if [ "$pinned" = no ]; then
            fail 2
            continue
        fi
        status=0
        module-kit build "$path" "${rows}/package" > /dev/null || status=$?
        if [ "$status" -ne 0 ]; then
            echo "eval-match: module-kit cannot build ${path}, so its samples are not held to its package" >&2
            fail 2
            continue
        fi
        if [ -s "${rows}/synthetic" ]; then
            status=0
            {{quote(just_executable())}} _match-synthetic "$module" "${rows}/package" "${rows}/synthetic" || status=$?
            if [ "$status" -gt 2 ]; then status=2; fi
            fail "$status"
        fi
        if [ -s "${rows}/people" ]; then
            status=0
            {{quote(just_executable())}} _match-people "novatrail/${module}" "${rows}/package" "${rows}/people" || status=$?
            if [ "$status" -gt 2 ]; then status=2; fi
            fail "$status"
        fi
    done
    if [ "$worst" -eq 0 ] && [ -z "$held" ]; then
        echo "eval-match: no module holds a sample"
    elif [ "$worst" -eq 0 ]; then
        echo "eval-match: every sample of every eval set is one of the package of this commit"
    fi
    exit "$worst"

# The synthetic samples of modules/<module>/ in rows — the file, the task, the seed and the tag each
# names —, drawn again by canvas's generator on the package, with the tag of each file, not the
# newest: the tag says only what a sample was made on (contracts/README.md in spec, section 13).
# Each is then compared with its file in all but task.modules.version and
# task.modules.package_digest, which that section bars from the comparison, so that a change of
# another task of the module, which changes the package's version and digest, re-makes none of
# them; a change of the task's open view, of its picture or of the generator does. The generator is
# called on 50 samples at a time, of one tag, into a new folder: it holds every sample of a call in
# memory, and writes none when one call fails. A batch it refuses is then drawn again a sample at a
# time, so that each sample is compared or refused with the generator's own reason, which it prints
# on stderr. The status is 1 when a sample is not drawn again the same or the generator refuses it,
# and 2 when it cannot do its work.
[private]
_match-synthetic module package rows:
    #!/usr/bin/env bash
    set -euo pipefail
    module={{quote(module)}}
    package={{quote(package)}}
    rows={{quote(rows)}}
    us=$'\x1f'
    generator=$({{quote(just_executable())}} _generator)
    work=$(mktemp -d)
    trap 'rm -rf "$work"' EXIT
    batch=50
    worst=0
    fail() {
        if [ "$1" -gt "$worst" ]; then worst=$1; fi
    }
    same=0
    n=0
    # Compare the file of calls[i] with the drawing of it in the folder given.
    compare() {
        local i=$1 out=$2 file=${files[$1]} left right
        local again="${out}/${file##*/}" remake="re-make it: just synthetic ${module} ${calls[$1]}"
        if [ ! -f "$again" ]; then
            echo "${file}: canvas's generator draws ${calls[$i]} as another file, so this one is no sample it draws; ${remake}" >&2
            fail 1
        elif ! left=$(jq -cS 'del(.task.modules.version, .task.modules.package_digest)' "$file") \
            || ! right=$(jq -cS 'del(.task.modules.version, .task.modules.package_digest)' "$again"); then
            echo "_match-synthetic: jq cannot read ${file} or its drawing, so they are not compared" >&2
            fail 2
        elif [ "$left" != "$right" ]; then
            echo "${file}: is not what canvas's generator draws of ${calls[$i]} on this commit's package, with the tag ${tag}; ${remake}" >&2
            fail 1
        else
            same=$((same + 1))
        fi
    }
    while IFS= read -r tag; do
        files=()
        calls=()
        while IFS=$us read -r file task seed _; do
            files+=("$file")
            calls+=("${task}:${seed}")
        done < <(awk -F "$us" -v tag="$tag" '$4 == tag' "$rows")
        for ((start = 0; start < ${#calls[@]}; start += batch)); do
            n=$((n + 1))
            out="${work}/${n}"
            status=0
            node "$generator" --module "$package" "$tag" "$out" "${calls[@]:start:batch}" > /dev/null 2>&1 || status=$?
            if [ "$status" -eq 0 ]; then
                for ((i = start; i < start + batch && i < ${#calls[@]}; i++)); do compare "$i" "$out"; done
                rm -rf "$out"
                continue
            fi
            for ((i = start; i < start + batch && i < ${#calls[@]}; i++)); do
                out="${work}/${n}-${i}"
                status=0
                node "$generator" --module "$package" "$tag" "$out" "${calls[$i]}" > /dev/null || status=$?
                if [ "$status" -eq 0 ]; then
                    compare "$i" "$out"
                else
                    if [ "$status" -gt 2 ]; then status=2; fi
                    echo "${files[$i]}: canvas's generator draws ${calls[$i]} no more on this commit's package of novatrail/${module}, with the tag ${tag}, for the reason above; re-make or remove it" >&2
                    fail "$status"
                fi
                rm -rf "$out"
            done
        done
    done < <(cut -d "$us" -f 4 "$rows" | sort -u)
    echo "_match-synthetic: ${same} synthetic samples of novatrail/${module} are drawn again the same on this commit's package"
    exit "$worst"

# The samples of people of a module in rows — the file, the task, the seed of its instance when the
# task is a template, and the digest of the open view each names —, held to the open views of the
# package. module-kit open, over the package and an index of it alone that module-kit index builds,
# with the seeds of the samples — 0 when none has one, since open asks for a seed —, gives the
# digest of each task's open view, and of an instance of a template by its seed, as the runner will
# read the task (architecture 3.6 module-kit; R93 spec). A sample whose task is not in the package,
# or not with that seed, or whose open view has another digest, is refused with 1: it was made on
# another open view, and goes (amendment 2026-11-06 to R04). Nothing of a sample but its file, task
# and seed is printed. The status is 2 when module-kit cannot open the package.
[private]
_match-people module package rows:
    #!/usr/bin/env bash
    set -euo pipefail
    module={{quote(module)}}
    package={{quote(package)}}
    rows={{quote(rows)}}
    us=$'\x1f'
    work=$(mktemp -d)
    trap 'rm -rf "$work"' EXIT
    seeds=()
    while IFS= read -r seed; do
        seeds+=(--seed "$seed")
    done < <(cut -d "$us" -f 3 "$rows" | { grep . || true; } | sort -u)
    if [ "${#seeds[@]}" -eq 0 ]; then seeds=(--seed 0); fi
    if ! module-kit index --first --tag v0.0.0 --no-withdrawals --out "${work}/module-index.json" "$package" > /dev/null \
        || ! module-kit open --index "${work}/module-index.json" --latest "${work}/module-index.json" "${seeds[@]}" \
            --out "${work}/open" "$package" > /dev/null; then
        echo "_match-people: module-kit cannot open the package of ${module}, so its samples of people are not held to it" >&2
        exit 2
    fi
    # The views by task and seed, as strings: a seed is up to 2^64, past what awk holds as a number.
    declare -A views=()
    if ! listed=$(jq -r --arg id "$module" '.modules[] | select(.id == $id) | .tasks[] | [.task.id, (.seed // ""), .open_view_digest] | join("\u001f")' \
        "${work}/open/open-views.json"); then
        echo "_match-people: the open views module-kit wrote of ${module} cannot be read" >&2
        exit 2
    fi
    while IFS=$us read -r id seed digest; do
        if [ -n "$id" ]; then views["${id}${us}${seed}"]=$digest; fi
    done <<< "$listed"
    worst=0
    held=0
    while IFS=$us read -r file task seed digest; do
        instance=$task
        if [ -n "$seed" ]; then instance="${task} with the seed ${seed}"; fi
        open=${views["${task}${us}${seed}"]:-}
        if [ -z "$open" ]; then
            echo "${file}: the package of this commit has no open view of ${instance}; remove the sample (amendment 2026-11-06 to R04)" >&2
            worst=1
        elif [ "$open" != "$digest" ]; then
            echo "${file}: the open view of ${instance} is not the one the sample was made on; remove the sample (amendment 2026-11-06 to R04)" >&2
            worst=1
        else
            held=$((held + 1))
        fi
    done < "$rows"
    echo "_match-people: ${held} samples of people of ${module} are of open views of this commit's package"
    exit "$worst"

# What .github's shared release.yml runs to build a release's files, in the toolchain stage of the
# commit it releases, on that commit's tree from git archive: no git, no submodule — spec/ is an
# empty folder —, no token and no network (architecture 3.3 .github; R03; R17). version is the
# release's tag, vX.Y.Z; previous the path of module-index.json of the newest release, which the
# shared release downloaded and release-verify accepted before the build, read-only, or an empty
# string for the first release, which the recipe tells module-kit index by --first; and no
# withdrawals.yaml is told by --no-withdrawals: nothing is left unsaid (R63 spec). A relative
# previous is from where the recipe is called. Every module is checked by _modules-check — the
# tree, the pinned CLI, check and its ID — and built — the recipe cannot tell which changed, with
# no git —, and index keeps each entry of the
# previous index as it was: a package of a version it names keeps its release, and another digest
# for it is refused by module-index.one-package. So an archive goes out only for a package whose
# entry names this release: sha256-<hex>.tar.gz, named by its package_digest, which is the sha256 of
# the bytes of its inventory.json, the package's files at the archive's root, with no top folder, no
# ./ and nothing else. The archive is the same bytes on every runner: names sorted, no owner, the
# same mode and time for every entry, and no name or time in the gzip header. dist/ is made afresh
# and holds module-index.json, LICENSE, the root's Apache-2.0, and the archives; a failed run leaves
# no dist/. The release then signs each file.
# Build a release's files into dist/, for a version such as v0.1.0 and the previous index or "".
release-files version previous:
    #!/usr/bin/env bash
    set -euo pipefail
    rm -rf dist
    version={{quote(version)}}
    previous={{quote(previous)}}
    if ! [[ "$version" =~ ^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
        echo "release-files: the version is a tag vX.Y.Z, as v0.1.0, not ${version}" >&2
        exit 2
    fi
    if [ -n "$previous" ] && [[ "$previous" != /* ]]; then
        previous={{quote(invocation_directory())}}/$previous
    fi
    if [ -z "$previous" ]; then
        against=(--first)
    elif [ -f "$previous" ] && [ -r "$previous" ]; then
        against=(--previous "$previous")
    else
        echo "release-files: ${previous} is no file that can be read, and the first release says so by an empty string" >&2
        exit 2
    fi
    if [ -L withdrawals.yaml ] || { [ -e withdrawals.yaml ] && [ ! -f withdrawals.yaml ]; }; then
        echo "release-files: withdrawals.yaml is no regular file" >&2
        exit 2
    elif [ -f withdrawals.yaml ]; then
        withdrawals=(--withdrawals withdrawals.yaml)
    else
        withdrawals=(--no-withdrawals)
    fi
    shopt -s nullglob dotglob
    entries=(modules/*)
    if [ -L modules ] || [ ! -d modules ] || [ "${#entries[@]}" -eq 0 ]; then
        echo "release-files: modules/ holds no module, and a release has nothing to index" >&2
        exit 2
    fi
    work=$(mktemp -d)
    trap 'rm -rf "$work"' EXIT
    # The tree, the pinned CLI, check and the ID of each module, by the one copy of those rules:
    # _modules-check, against no published version. The comparison with them is CI's, on this same
    # commit; what a release needs of it, index holds: one-package and version-above.
    mkdir "${work}/published"
    echo "release-files checks the modules; the published versions are CI's to compare with" > "${work}/published/first"
    {{quote(just_executable())}} _modules-check "${work}/published"
    mkdir "${work}/packages" "${work}/dist"
    packages=()
    for path in "${entries[@]}"; do
        package="${work}/packages/${path#modules/}"
        module-kit build "$path" "$package"
        packages+=("$package")
    done
    index="${work}/dist/module-index.json"
    module-kit index "${against[@]}" --tag "$version" "${withdrawals[@]}" --out "$index" "${packages[@]}"
    for package in "${packages[@]}"; do
        id=$(jq -er '.id | strings' "${package}/inventory.json")
        module=$(jq -er '.version | strings' "${package}/inventory.json")
        if ! entry=$(jq -ec --arg id "$id" --arg v "$module" '.modules[$id][$v] | objects' "$index"); then
            echo "release-files: the index module-kit built names no entry of ${id} ${module}" >&2
            exit 2
        fi
        named=$(jq -er '.release | strings' <<< "$entry")
        digest=$(jq -er '.package_digest | strings' <<< "$entry")
        if [ "$named" != "$version" ]; then
            echo "release-files: ${id} ${module} stays in ${named}, which published it"
            continue
        fi
        own="sha256:$(sha256sum < "${package}/inventory.json" | cut -d ' ' -f 1)"
        if [ "$own" != "$digest" ]; then
            echo "release-files: ${id} ${module}: the index names the package ${digest}, and its inventory.json is ${own}" >&2
            exit 2
        fi
        found=()
        for member in "$package"/*; do
            found+=("${member##*/}")
        done
        if [ "${found[*]}" != "device device.json inventory.json server.json" ] || [ -L "${package}/device" ] || [ ! -d "${package}/device" ]; then
            echo "release-files: ${id} ${module}: the package holds ${found[*]}, not inventory.json, server.json, device.json and the folder device/ alone (architecture 3.3)" >&2
            exit 2
        fi
        archive="${work}/dist/${digest/:/-}.tar.gz"
        tar --create --format=gnu --sort=name --owner=0 --group=0 --numeric-owner \
            --mode=u=rwX,go=rX --mtime=@0 --directory="$package" --file=- \
            -- inventory.json server.json device.json device \
            | gzip -n -9 > "$archive"
        echo "release-files: ${id} ${module}, ${digest}, in ${archive##*/}"
    done
    cp LICENSE "${work}/dist/LICENSE"
    mv "${work}/dist" dist
    echo "release-files: the files of ${version} are in dist/"

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

# The generator of canvas that the toolchain stage installs is held by its sum in
# .devcontainer/Dockerfile, as module-kit's binary is (R29). Here the release's SHA256SUMS and its
# bundle are read from the release's public address, with no token, the signature is checked by
# spec's release-verify (R77 spec), and the sum of synthetic-module_<tag>.mjs is held to the
# Dockerfile's. CI runs it on every pull request and push to main.
# Verify the signature of the canvas release whose generator the Dockerfile pins, and its sum.
canvas-verify:
    #!/usr/bin/env bash
    set -euo pipefail
    tag=$({{quote(just_executable())}} _dockerfile-arg .devcontainer/Dockerfile CANVAS_VERSION)
    if ! [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "canvas-verify: CANVAS_VERSION is ${tag}, not a tag vX.Y.Z" >&2
        exit 1
    fi
    release=$(mktemp -d)
    trap 'rm -rf "$release"' EXIT
    for file in SHA256SUMS SHA256SUMS.sigstore.json; do
        {{quote(just_executable())}} _fetch "https://github.com/altnovatrail/canvas/releases/download/${tag}/${file}" "${release}/${file}"
    done
    {{quote(just_executable())}} -f spec/justfile release-verify canvas "$tag" "${release}/SHA256SUMS"
    {{quote(just_executable())}} _canvas-sum "${release}/SHA256SUMS" .devcontainer/Dockerfile

# The sum of synthetic-module_<tag>.mjs in a release's SHA256SUMS, listed once, against
# CANVAS_SYNTHETIC_SHA256, for the tag CANVAS_VERSION, of a Dockerfile. Apart from the network, so
# that its tests need none.
[private]
_canvas-sum sums dockerfile:
    #!/usr/bin/env bash
    set -euo pipefail
    sums={{quote(sums)}}
    dockerfile={{quote(dockerfile)}}
    tag=$({{quote(just_executable())}} _dockerfile-arg "$dockerfile" CANVAS_VERSION)
    wanted=$({{quote(just_executable())}} _dockerfile-arg "$dockerfile" CANVAS_SYNTHETIC_SHA256)
    name="synthetic-module_${tag}.mjs"
    if ! [[ "$wanted" =~ ^[0-9a-f]{64}$ ]]; then
        echo "_canvas-sum: CANVAS_SYNTHETIC_SHA256 is not a sha256 of 64 lowercase hexadecimal digits" >&2
        exit 1
    fi
    listed=$(awk -v name="$name" '$2 == name { print $1 }' "$sums")
    count=$(grep -c . <<< "$listed" || true)
    if [ "$count" -ne 1 ]; then
        echo "_canvas-sum: ${sums} lists ${name} not once but ${count} times" >&2
        exit 1
    fi
    if [ "$listed" != "$wanted" ]; then
        echo "_canvas-sum: ${name} is ${listed} in SHA256SUMS, and CANVAS_SYNTHETIC_SHA256 is ${wanted}" >&2
        exit 1
    fi
    echo "canvas ${tag}: ${name} ${listed}"

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

# The tests run this same just, and the module-kit on the PATH: the toolchain stage's. Those of
# eval-check run spec's sample-check too, from the spec/ submodule, which installs its tools from
# their lock file through npm once, as docs-check-repo does.
# Run the tests of the recipes.
test:
    JUST={{quote(just_executable())}} node --test 'test/*.test.mjs'

# The Go code of this repository is tools/port, the port of standalone's tasks (R12), one module
# with its go.mod at the root; spec/ is a module of its own, which ./... does not reach. The ci-
# recipes are what .github's shared go.yml runs, one job each, in the toolchain stage, with no spec
# submodule: "Go / lint", "Go / test", "Go / vuln" and "Go / licenses" (R16 .github; R48 spec). The
# tools only these checks need run straight from their modules at exact versions, as spec's and
# module-kit's justfiles run them; Dependabot does not see them, so they are raised by hand.
GOVULNCHECK := "golang.org/x/vuln/cmd/govulncheck@v1.9.0"
GO_LICENSES := "github.com/google/go-licenses/v2@v2.0.1"

# Fail on Go code gofmt would rewrite, and on what go vet reports.
ci-lint:
    #!/usr/bin/env bash
    set -euo pipefail
    unformatted="$(git ls-files -z --cached --others --exclude-standard -- '*.go' | xargs -0 -r gofmt -l)"
    if [ -n "$unformatted" ]; then
        echo "ci-lint: gofmt would rewrite these files; run gofmt -w on them:" >&2
        echo "$unformatted" >&2
        exit 1
    fi
    go vet ./...

# The profile goes into coverage/go.out, which the job coverage of go.yml sends to Codecov, and
# codecov/patch holds the changed lines to 90 % (.github/codecov.yml; R48 spec).
# Run the Go tests, with the race detector, and write their coverage.
ci-test:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p coverage
    go test -race -count=1 -covermode=atomic -coverprofile=coverage/go.out ./...

# Fail on a known vulnerability in code the module reaches, in the standard library of its Go too.
ci-vuln:
    go run {{GOVULNCHECK}} ./...

# The list is ALLOWED_LICENSES, SPDX identifiers separated by commas, which go.yml sets; without it,
# the recipe takes spec's list, as `just -f spec/justfile licenses-allowed` prints it (R12 spec).
# Fail if a package the module is built from carries a license outside that list.
ci-licenses:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -n "${ALLOWED_LICENSES:-}" ]; then
        allowed=$ALLOWED_LICENSES
    else
        allowed=$({{quote(just_executable())}} -f spec/justfile licenses-allowed | paste -sd , -)
    fi
    go run {{GO_LICENSES}} check ./... --allowed_licenses="$allowed"

# The tasks of a batch of tools/port/batches/ come from standalone at the commit the batch names: a
# clone of it in a temporary folder, never in this repository (RULES.md, Citations), read with no
# token and none of the caller's git settings. Standalone's bench proves the answer of each task
# first, and the port writes nothing unless each passes (R05; R12). It is a step of the batch's
# task, run by a person in the devcontainer, not of CI: the bench fetches standalone's Go modules.
# After the module come the options: --batch <file> in place of the first batch, whose overlays are
# then tools/port/overlay/<the batch's name>/ (T22), --overlay <folder> in place of those, and --check
# to compare the module with what the port writes.
# Write the tasks of a batch of standalone into a module, or with --check compare them.
[positional-arguments]
port module *options:
    #!/usr/bin/env bash
    set -euo pipefail
    module=$1
    shift
    batch=tools/port/batches/first.yaml
    overlay=""
    rest=()
    while [ "$#" -gt 0 ]; do
        case "$1" in
            --batch) batch=$2; shift 2 ;;
            --overlay) overlay=$2; shift 2 ;;
            *) rest+=("$1"); shift ;;
        esac
    done
    if [ -z "$overlay" ]; then
        overlay="tools/port/overlay/$(basename "$batch" .yaml)"
    fi
    commit=$(sed -n 's/^standalone: \([0-9a-f]\{40\}\)$/\1/p' "$batch")
    if [ -z "$commit" ]; then
        echo "port: ${batch} has no line standalone: with the full SHA of a commit" >&2
        exit 2
    fi
    clone=$(mktemp -d)
    trap 'rm -rf "$clone"' EXIT
    export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 GIT_TERMINAL_PROMPT=0
    git clone --quiet --filter=blob:none --no-checkout https://github.com/mathtrail/mathtrail-standalone "$clone"
    git -C "$clone" checkout --quiet --detach "$commit"
    go run ./tools/port --standalone "$clone" --batch "$batch" --overlay "$overlay" --module "$module" --prove "${rest[@]}"

# The pictures of the lessons of a module are tools/port/pictures/<the module's folder>/<name>.json,
# in the schema of standalone's pictures, which the port draws into media/lessons/<name>.png of the
# module, with no clone, as it draws a task's picture (amendment T41 to R12). CI runs it with
# --check, which writes nothing and names every picture that differs from what the port draws.
# Draw the pictures of the lessons into a module, or with --check compare them.
[positional-arguments]
lesson-pictures module *options:
    go run ./tools/port --pictures "tools/port/pictures/$(basename "$1")" --module "$@"

# The pictures of the instances of a module's templates: a template whose picture path holds
# placeholders has a source, tools/port/template-pictures/<the module's folder>/<the path without
# its placeholders>.json, in the schema of standalone's pictures with placeholders among its
# strings, and the port draws it for every instance into media/pictures/ of the module, as the
# format asks (amendment T52 to R12). CI runs it with --check.
# Draw the pictures of the instances of a module's templates, or with --check compare them.
[positional-arguments]
template-pictures module *options:
    go run ./tools/port --template-pictures "tools/port/template-pictures/$(basename "$1")" --module "$@"

# A survey runs the port over every task of standalone that no batch of tools/port/batches/ names
# yet, at a commit, 42b4408 unless another is given, before E3's batches are written (R20): the
# clone as port makes it, then _port-survey with standalone's bench. The module is a folder that is
# not there yet, which the survey makes and which never enters the repository. It is a step of E3's
# tasks, run by a person in the devcontainer, not of CI.
# Survey the tasks of standalone that no batch names, into a new module of the format's examples.
[positional-arguments]
port-survey module commit="42b4408225e4163d8efbdfb93c299f3fc4b70afa":
    #!/usr/bin/env bash
    set -euo pipefail
    clone=$(mktemp -d)
    trap 'rm -rf "$clone"' EXIT
    export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 GIT_TERMINAL_PROMPT=0
    git clone --quiet --filter=blob:none --no-checkout https://github.com/mathtrail/mathtrail-standalone "$clone"
    git -C "$clone" checkout --quiet --detach "$2"
    {{quote(just_executable())}} _port-survey "$clone" "$2" "$1" --prove

# The module of a survey is the format's examples of olympiad-math — manifest, skill graph,
# techniques and solving steps — with the traps of modules/olympiad-math, whose IDs the port writes,
# one of them its own until the example has it (amendment T50 to R23), the root's LICENSE, as in
# test/port.test.mjs, and the tasks the port writes into it. The example's graph has the 17 topics
# of standalone at c6cbeb2, so the survey adds the two topics of 42b4408 where the clone has them
# and the graph does not, with no prerequisites, as their builds_on is empty, and the topics' names
# for titles (docs/architecture.md, section 3.5; R21). They go at the graph's end, under skills:,
# which must be its last key. module-kit check then reads it: a breach of task.no-leaks is what E3's
# overlays are to mend, and the recipe counts them; any other breach, and any warning, fails it.
# The exit status is the port's — 1 when it refused a task, 2 when it could not run — or 1 when the
# check found more than task.no-leaks. The port is built first, rather than run by go run, which
# makes every status other than 0 into 1.
# Survey a clone of standalone into a new module, and check the module.
[private]
[positional-arguments]
_port-survey clone commit module *options:
    #!/usr/bin/env bash
    set -euo pipefail
    clone=$1 commit=$2 module=$3
    shift 3
    if [ -e "$module" ]; then
        echo "port-survey: ${module} is there already: a survey makes its module anew" >&2
        exit 2
    fi
    mkdir -p "$module"
    for part in manifest skill-graph techniques solving-steps; do
        cp "spec/format/examples/${part}/valid/olympiad-math.yaml" "${module}/${part}.yaml"
    done
    cp modules/olympiad-math/traps.yaml "${module}/traps.yaml"
    cp LICENSE "${module}/LICENSE"
    graph="${module}/skill-graph.yaml"
    if [ "$(awk '/^[a-z_]+:/ { key = $1 } END { print key }' "$graph")" != "skills:" ]; then
        echo "port-survey: the format's example graph does not end with skills:, so no skill can be added to it" >&2
        exit 2
    fi
    while IFS='|' read -r topic skill title; do
        if jq -e --arg id "$topic" 'any(.[]; .id == $id)' "${clone}/content/catalogs/topics.json" > /dev/null &&
            ! grep -q "^  ${skill}:$" "$graph"; then
            printf '  %s:\n    title: %s\n    assessability: verdict\n' "$skill" "$title" >> "$graph"
        fi
    done <<< $'counting.figures|counting-figures|Counting figures\ngeometry.cubes_paper|geometry-cubes-paper|Cubes and paper'
    bin=$(mktemp -d)
    trap 'rm -rf "$bin"' EXIT
    go build -o "${bin}/port" ./tools/port
    status=0
    "${bin}/port" --standalone "$clone" --survey --commit "$commit" --except tools/port/batches --module "$module" "$@" || status=$?
    if [ "$status" -eq 2 ]; then
        exit 2
    fi
    checked=$(module-kit check "$module" 2>&1) || true
    leaks=$(grep -c ': task\.no-leaks: ' <<< "$checked" || true)
    other=$(grep -v -e ': task\.no-leaks: ' -E -e ': the format [0-9.]* refuses the module: [0-9]* breach(es)?$' \
        -e ': the module breaks no rule of the format [0-9.]*$' <<< "$checked" || true)
    echo "port-survey: module-kit check ${module}: task.no-leaks ${leaks}, for the overlays to mend"
    if [ -n "$other" ]; then
        echo "port-survey: module-kit check found more than task.no-leaks:" >&2
        echo "$other" >&2
        status=1
    fi
    exit "$status"
