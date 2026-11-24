# modules — sources

> The facts of the outside world the repository's documents under `docs/` rest on and no run of ours checks: a law, a regulator, a standard's normative text, a provider's terms and policies, a price, a license, a page that changes in place with no version, a text outside the repositories. Here are their records — the address, the date it was read and a quote in the source's language — and [architecture.md](architecture.md) and the decision log hold the claim and a citation of its one record by its number, `sources 2.4`, not a chain. A record's heading names its topic, not a task.
>
> No record is needed for a library's documentation or code at a version or commit — a link beside the claim —, for a behaviour a test or a measurement checks — they are its evidence —, or for spec's tags and the files of NovaTrail's and MathTrail's repositories — git, at a commit.
>
> A record is a subsection "N.M" of a topic: one source, or one reading of several sources for one decision. A record's number never changes and is never given again: a new one goes at the end of its topic, a new topic at the end of the file, and a record removed leaves a gap. A source read again with the same text changes the record's date; with another text, the claims and decisions on it are rewritten in the same change and the record holds the new reading alone — where a decision waits on the author, an open question, and both readings stand until it is closed. The model is `spec/docs/sources.md`. The results of measurements are not here but in [measurements.md](measurements.md).

## 1. How to keep it

- **A record** holds the address, the date of the page where it has one, the date read and the quote, in the language of the source, and its first line says in words what it was read for. What was not read from the primary source — a copy, a pointer, a page that did not open — the record names as such.
- **A citation.** A claim of a document cites its one record by its number: `(sources 2.4)`. The archives of `docs/plan/` keep citing a removed record, and `docs-check` does not look up their citations of records.
- **What is not here:** the quotes that an entry of the decision log gives itself stay in the entry, and those that [architecture.md](architecture.md) gives as the argument of its item stay where they stand.

## 2. GitHub

GitHub's documentation at a commit of [github/docs](https://github.com/github/docs) is a link where the claim stands; a limit of the service or its terms, read from docs.github.com with no version, is a record.

### 2.1. The files of a GitHub release

For the release.

[About releases](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases), the section "Storage and bandwidth quotas", read 2026-11-02: "Up to 1000 release assets may be associated with a single release. Each file included in a release must be under 2 GiB." A release file has only a name, without folders ([Upload a release asset](https://docs.github.com/en/rest/releases/assets#upload-a-release-asset), ibid.: "GitHub renames asset filenames that have special characters, non-alphanumeric characters, and leading or trailing periods").

### 2.4. The age of a GitHub account

For the samples of the eval set.

[GitHub Terms of Service](https://docs.github.com/en/site-policy/github-terms/github-terms-of-service), B. Account Terms, 3. Account Requirements, in effect since 27 April 2026, read 2026-11-02: "You must be age 13 or older."

## 3. Standards and libraries

The documentation of git and of a library at a version is a link where the claim stands; the normative text of a standard, and a text that certifies a contribution, are records.

### 3.2. DCO 1.1

For the samples of the eval set.

[Developer Certificate of Origin 1.1](https://developercertificate.org/), read 2026-11-02:
- (a) "The contribution was created in whole or in part by me and I have the right to submit it under the open source license indicated in the file";
- (b) "The contribution is based upon previous work that, to the best of my knowledge, is covered under an appropriate open source license and I have the right under that license to submit that work with modifications…";
- (c) "The contribution was provided directly to me by some other person who certified (a), (b) or (c) and I have not modified it";
- (d) "a record of the contribution (including all personal information I submit with it, including my sign-off) is maintained indefinitely".

### 3.3. PNG chunks

For the samples of the eval set.

[Portable Network Graphics (PNG) Specification (Third Edition)](https://www.w3.org/TR/png-3/), W3C Recommendation 24 June 2025, read 2026-11-02:
- "A critical chunk is a chunk that is absolutely required in order to successfully decode a PNG image from a PNG datastream" (there, 11.2). There are four critical chunks — IHDR, PLTE, IDAT and IEND —, and the others are ancillary (there, 4.8.2);
- not every image needs PLTE: "This chunk shall appear for color type 3, and may appear for color types 2 and 6; it shall not appear for color types 0 and 4" (there, 11.2.2);
- among the ancillary ones are tEXt, iTXt and zTXt, "Textual information (which may be compressed) associated with the image", eXIf, "Exchangeable image file format metadata", and tIME, "The time when the PNG image was last modified" (there, 4.7, table 2).

## 4. NovaTrail's repositories

Their tags and files are git: a link to the repository at its tag or commit stands where the claim stands, and there is no record.

## 5. MathTrail: standalone and its prototype

Their files are git: a link to the repository at its commit stands where the claim stands, and there is no record.
