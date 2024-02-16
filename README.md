# :white_check_mark: htmltest

[![Go Report Card](https://goreportcard.com/badge/github.com/wjdp/htmltest)](https://goreportcard.com/report/github.com/wjdp/htmltest)
[![GoDoc](https://godoc.org/github.com/wjdp/htmltest?status.svg)](https://godoc.org/github.com/wjdp/htmltest)

If you generate HTML files, [html-proofer](https://github.com/gjtorikian/html-proofer) might be the tool for you. If you can't be bothered with a Ruby environment or fancy something a bit faster, htmltest may be a better option.

:mag: htmltest runs your HTML output through a series of checks to ensure all your links, images, scripts references work, your alt tags are filled in, *et cetera*.

:horse_racing: *Faster?* Yep, quite a bit actually. On [a site](https://github.com/newtheatre/history-project) with over 2000 files htmlproofer took over [three minutes](https://travis-ci.org/newtheatre/history-project#L564), htmltest took [8.6 seconds](https://travis-ci.org/newtheatre/history-project#L538). Both tools had full valid caches.

:confused: *Why make another tool*: A mix of frustration with using htmlproofer/Ruby on large sites and needing a good project to get to grips with Go.

## :floppy_disk: Installation

### :penguin: Linux / :green_apple: macOS

This [godownloader](https://github.com/goreleaser/godownloader) script will query GitHub for the latest release and download the correct binary for your platform into the directory set with the `-b` flag.

#### System-wide Install

```bash
curl https://htmltest.wjdp.uk | sudo bash -s -- -b /usr/local/bin
```

You'll be prompted for your password. After simply do `htmltest` to run.

#### Into Current Directory

```bash
curl https://htmltest.wjdp.uk | bash
```

By default this will install `htmltest` into `./bin` of your current directory, to run do `bin/htmltest`. Rather suitable for CI environments.

#### More options

Run `curl https://htmltest.wjdp.uk | bash -s -- -h` for help text.

### 📦 Package Repos

We're available in some [package repositories](https://repology.org/project/htmltest/versions) such as Homebrew and repos for certain distros. Do note the version available may not be the latest.

[![Packaging status](https://repology.org/badge/vertical-allrepos/htmltest.svg)](https://repology.org/project/htmltest/versions)

### ![win64](https://user-images.githubusercontent.com/1690934/30242799-17a573f2-9595-11e7-9aa5-04e34b04b0cd.png) Windows

:arrow_down: Download the [latest binary release](https://github.com/wjdp/htmltest/releases/latest) and put it somewhere on your PATH.

### :whale: Docker

```docker run -v $(pwd):/test --rm wjdp/htmltest```
Mount your directory with html files into the container and test them.

If you need more arguments to the test run it like this:
```docker run -v $(pwd):/test --rm wjdp/htmltest -l 3 -s```

### 🧾 Temporary Files

We store temporary files in `tmp/.htmltest` by default. This contains a log of output and a cache of external links, speeding up subsequent runs. You probably want to ignore it in version control, and perhaps cache it in your CI system.

## :computer: Usage

```txt
htmltest - Test generated HTML for problems
           https://github.com/wjdp/htmltest

Usage:
  htmltest [options] [<path>]
  htmltest -v --version
  htmltest -h --help

Options:
  <path>                       Path to directory or file to test, if omitted we
                               attempt to read from .htmltest.yml.
  -c FILE, --conf FILE         Custom path to config file.
  -h, --help                   Show this text.
  -l LEVEL, --log-level LEVEL  Logging level, 0-3: debug, info, warning, error.
  -s, --skip-external          Skip external link checks, may shorten execution
                               time considerably.
  -v, --version                Show version and build time.
```

## :microscope: What's Tested?

Many options of the following tests can customised. Items marked :soon: are not checked yet, but will be *soon*.

- `a` `link` `img` `script`: Whether internal links work / are valid.
- `a`: Whether internal hashes work.
- `a` `link` `img` `script`: Whether external links work.
- `a`: :soon: Whether external hashes work.
- `a` `link`: Whether external links use HTTPS.
- `img`: Whether your images have valid alt attributes.
- `link`: Whether pages have a valid favicon.
- `meta`: Whether refresh tags are valid and the url works.
- `meta`: :soon: Whether images and URLs in the OpenGraph metadata are valid.
- `meta` `title`: :soon: Whether you've got the [recommended tags](https://support.google.com/webmasters/answer/79812?hl=en) in your head.
- `DOCTYPE`: Whether a doctype is correctly specified.

### What's Not

I'd like to test the following but won't be for a while.

- Whether your HTML markup is valid. htmlproofer has the ruby library [Nokogiri](http://www.nokogiri.org/tutorials/ensuring_well_formed_markup.html), I've not found one for Go yet.

## :see_no_evil: Ignoring content

Add the `data-proofer-ignore` attribute to any tag or to the class of a tag to ignore it from every check. The name of this attribute can be customised.

```html
<a href="http://notareallink" data-proofer-ignore>Not checked.</a>
```

## :bookmark_tabs: Caching

Checking external URLs can slow tests down and potentially annoy the URL's host. htmltest caches the status code of checked external URLs and stores this cache between runs. We write the cache to `tmp/.htmltest/refcache.json` and expire items after two weeks by default.

## :rainbow: Colour Output

By default, output in the TTY uses colours to indicate warnings, errors, and success. To turn off colourization, set an environment variable named `NO_COLOR`. If it is present, no colour will be used. The value is ignored. (See [no-color.org](https://no-color.org/).)

## :fax: Logging

If you've got a lot of errors, reading them off a TTY may be difficult. We write errors to `tmp/.htmltest/htmltest.log` by default. The log level is set in the config file.

## :wrench: Configuration

htmltest uses a YAML configuration file. Put `.htmltest.yml` in the same directory that you're running the tool from and you can just say `htmltest` to run your tests. You'll probably also want to cache the `tmp/.htmltest` directory.

### Basic Options

| Option | Description                                                                                                                                                                                                     | Default |
| :----- |:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------| :------ |
| `DirectoryPath` | Directory to scan for HTML files.                                                                                                                                                                               | |
| `DirectoryIndex` | The file to look for when linking to a directory.                                                                                                                                                               | `index.html` |
| `FilePath` | Single file to test within `DirectoryPath`, omit to test all.                                                                                                                                                   | |
| `FileExtension` | Extension of your HTML documents, includes the dot. If `FilePath` is set we use the extension from that.                                                                                                        | `.html` |
| `BaseURL` | Publication URL of the site, including subfolder if applicable.                                                                                                                                                 | |
| `CheckDoctype` | Enables checking the document type declaration.                                                                                                                                                                 | `true` |
| `CheckAnchors` | Enables checking `<a…` tags.                                                                                                                                                                                    | `true` |
| `CheckLinks` | Enables checking `<link…` tags.                                                                                                                                                                                 | `true` |
| `CheckImages` | Enables checking `<img…` tags                                                                                                                                                                                   | `true` |
| `CheckScripts` | Enables checking `<script…` tags.                                                                                                                                                                               | `true` |
| `CheckMeta` | Enables checking `<meta…` tags.                                                                                                                                                                                 | `true` |
| `CheckGeneric` | Enables other tags, see items marked with checkGeneric on the [tags wiki page](https://github.com/wjdp/htmltest/wiki/Tags).                                                                                     | `true` |
| `CheckExternal` | Enables external reference checking; all tag types.                                                                                                                                                             | `true` |
| `CheckInternal` | Enables internal reference checking; all tag types. When disabled will prevent internal hash checking unless the reference only contains a hash fragment (`#heading`) and therefore refers to the current page. | `true` |
| `CheckInternalHash` | Enables internal hash/fragment checking.                                                                                                                                                                        | `true` |
| `CheckMailto` | Enables–albeit quite basic–`mailto:` link checking.                                                                                                                                                             | `true` |
| `CheckTel` | Enables–albeit quite basic–`tel:` link checking.                                                                                                                                                                | `true` |
| `CheckFavicon` | Enables favicon checking, ensures every page has a favicon set.                                                                                                                                                 | `false` |
| `CheckMetaRefresh` | Enables checking meta refresh tags.                                                                                                                                                                             | `true` |
| `CheckSelfReferencesAsInternal` | Check external references starting with `BaseURL` as though they are internal references.                                                                                                                        | `false` |
| `EnforceHTML5` | Fails when the doctype isn't `<!DOCTYPE html>`.                                                                                                                                                                 | `false` |
| `EnforceHTTPS` | Fails when encountering an `http://` link. Useful to prevent mixed content errors when serving over HTTPS.                                                                                                      | `false` |
| `IgnoreURLs` | Array of regexs of URLs to ignore.                                                                                                                                                                              | empty |
| `IgnoreInternalURLs` | Array of strings of internal URLs to ignore. Exact matches only. ⚠ Likely to be deprecated, use `IgnoreURLs` instead.                                                                                           | empty |
| `IgnoreHTTPS` | Array of regexs of URLs to ignore for `EnforceHTTPS`. These URLs are still tested, unless also present in `IgnoreURLs`.                                                                                         | empty |
| `IgnoreDirs` | Array of regexs of directories to ignore when scanning for HTML files.                                                                                                                                          | empty |
| `IgnoreInternalEmptyHash` | When true prevents raising an error for links with `href="#"`.                                                                                                                                                  | `false` |
| `IgnoreEmptyHref` | When true prevents raising an error for links with `href=""`.                                                                                                                                                   | `false` |
| `IgnoreCanonicalBrokenLinks` | When true produces a warning, rather than an error, for broken canonical links. When testing a site which isn't live yet or before publishing a new page canonical links will fail.                             | `true` |
| `IgnoreExternalBrokenLinks` | When true produces a warning, rather than an error, for broken external links. Useful when testing a site having hundreds of external links.                                                                    | `false` |
| `IgnoreAltMissing` | Turns off image alt attribute checking.                                                                                                                                                                         | `false` |
| `IgnoreAltEmpty` | Allows `alt=""` for decorative images.                                                                                                                                                                          | `false` |
| `IgnoreDirectoryMissingTrailingSlash` | Turns off errors for links to directories without a trailing slash.                                                                                                                                             | `false` |
| `IgnoreSSLVerify` | Turns off x509 errors for self-signed certificates.                                                                                                                                                             | `false` |
| `IgnoreTagAttribute` | Specify the ignore attribute. All tags with this attribute or with this class will be excluded from every check.                                                                                                | `"data-proofer-ignore"` |
| `HTTPHeaders` | Dictionary of headers to include in external requests                                                                                                                                                           | `{"Range":  "bytes=0-0", "Accept": "*/*"}` |
| `TestFilesConcurrently` | :warning: :construction: *EXPERIMENTAL* Turns on [concurrent](https://github.com/wjdp/htmltest/wiki/Concurrency) checking of files.                                                                             | `false` |
| `DocumentConcurrencyLimit` | Maximum number of documents to process at once.                                                                                                                                                                 | `128` |
| `HTTPConcurrencyLimit` | Maximum number of open HTTP connections. If you raise this number ensure the `ExternalTimeout` is suitably raised.                                                                                              | `16` |
| `LogLevel` | Logging level, 0-3: debug, info, warning, error.                                                                                                                                                                | `2` |
| `LogSort` | How to sort/present issues. Can be `seq` for sequential output or `document` to group by document.                                                                                                              | `document` |
| `ExternalTimeout` | Number of seconds to wait on an HTTP connection before failing.                                                                                                                                                 | `15` |
| `RedirectLimit` | Allowed number of redirects. Use built-in behavior with negative values.                                                                                                                                        | `-1` |
| `StripQueryString` | Enables stripping of query strings from external checks.                                                                                                                                                        | `true` |
| `StripQueryExcludes` | List of URLs to disable query stripping on.                                                                                                                                                                     | `["fonts.googleapis.com"]` |
| `OutputDir` | Directory to store cache and log files in. Relative to executing directory.                                                                                                                                     | `tmp/.htmltest` |
| `OutputCacheFile` | File within `OutputDir` to store reference cache.                                                                                                                                                               | `refcache.json` |
| `OutputLogFile` | File within `OutputDir` to store last tests errors.                                                                                                                                                             | `htmltest.log` |
| `CacheExpires` | Cache validity period, accepts [go.time duration strings](https://golang.org/pkg/time/#ParseDuration) (…"m", "h").                                                                                              | `336h` (two weeks) |

### Example

```yaml
DirectoryPath: "_site"
EnforceHTTPS: true
IgnoreURLs:
- "example.com"
- "^/misc/js/script.js$"
IgnoreDirs:
- "lib"
CacheExpires: "6h"
```

## :loudspeaker: Issues? Suggestions?

[Submit an issue](https://github.com/wjdp/htmltest/issues/new).
