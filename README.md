# :white_check_mark: htmltest

If you generate HTML files, [html-proofer](https://github.com/gjtorikian/html-proofer) might be the tool for you. If you can't be bothered with a Ruby environment or fancy something a bit faster, htmltest may be a better option.

:mag: htmltest runs your HTML output through a series of checks to ensure all your links, images, scripts references work, your alt tags are filled in, *et cetera*.

:horse_racing: *Faster?* Yep, quite a bit actually. On [a site](https://github.com/newtheatre/history-project) with over 2000 files htmlproofer took over [three minutes](https://travis-ci.org/newtheatre/history-project#L564), htmltest took [8.6 seconds](https://travis-ci.org/newtheatre/history-project#L538). Both tools had full valid caches.

:confused: *Why make another tool*: A mix of frustration with using htmlproofer/Ruby on large sites and needing a good project to learn Go with. Yep, this is my first major Go program, it uses pretty much the same spec/tests as htmlproofer so should be accurate, but probably could do with some refactoring.

## :floppy_disk: Installation

Download the [latest binary release](https://github.com/wjdp/htmltest/releases) and stick it in your binary folder (`~/bin`).

For a CI environment–like Travis–put the desired version in the following and run: `curl -L https://github.com/wjdp/htmltest/releases/download/vX.X.X/htmltest-linux` then execute with `./htmltest-linux`

We store temporary files in `tmp/.htmltest` by default. You probably want to ignore that in your version control system.

## :computer: Usage

```
htmltest - Test generated HTML for problems
           https://github.com/wjdp/htmltest

Usage:
  htmltest
  htmltest [--log-level=LEVEL] <path>
  htmltest --conf=CFILE
  htmltest --version
  htmltest -h --help

Options:
  <path>              Path to directory or file to test, if omitted:
                      htmlproofer --conf=.htmltest.yml
  --log-level=LEVEL   Logging level, 0-3: debug, info, warning, error.
  --conf=CFILE        Custom path to config file.
  -h --help           Show this text.
```

## :microscope: What's Tested?

Many options of the following tests can customised. Items marked :soon: are not checked yet, but will be *soon*.

- `a` `link` `img` `script`: Whether internal links work / are valid.
- `a`: :soon: Whether internal hashes work.
- `a` `link` `img` `script`: Whether external links work.
- `a`: :soon: Whether external hashes work.
- `a` `link`: Whether external links use HTTPS.
- `a` `link`: Whether external links use HTTPS.
- `img`: Whether your images have valid alt attributes.
- `meta`: :soon: Whether favicons are valid.
- `meta`: :soon: Whether images and URLs in the OpenGraph metadata are valid.
- `meta` `title`: :soon: Whether you've got the [recommended tags](https://support.google.com/webmasters/answer/79812?hl=en) in your head.

### What's Not

I'd like to test the following but won't be for a while.

- Whether your HTML markup is valid. htmlproofer has the ruby library [Nokogiri](http://www.nokogiri.org/tutorials/ensuring_well_formed_markup.html), I've not found one for Go yet.

## :see_no_evil: Ignoring content

Add the `data-proofer-ignore` attribute to any tag to ignore it from every check. The name of this attribute can be customised.

```html
<a href="http://notareallink" data-proofer-ignore>Not checked.</a>
```

## :bookmark_tabs: Caching

Checking external URLs can slow tests down and potentially annoy the URL's host. htmltest caches the status code of checked external URLs and stores this cache between runs. We write the cache to `tmp/.htmltest/refcache.json` and expire items after two weeks by default.

## :fax: Logging

If you've got a lot of errors, reading them off a TTY may be difficult. We write errors to `tmp/.htmltest/htmltest.log` by default. The log level is set in the config file.

## :wrench: Configuration

htmltest uses a YAML configuration file. Put `.htmltest.yml` in the same directory that you're running the tool from and you can just say `htmltest` to run your tests. You'll probably also want to cache the `tmp/.htmltest` directory.

| Option | Description | Default |
| :----- | :---------- | :------ |
| `DirectoryPath` | Directory to scan for HTML files. | |
| `FilePath` | File to scan, omit if using `DirectoryPath`. | |
| `CheckAnchors` | Enables checking `<a…` tags. | `true` |
| `CheckLinks` | Enables checking `<link…` tags. | `true` |
| `CheckImages` | Enables checking `<img…` tags | `true` |
| `CheckScripts` | Enables checking `<script…` tags. | `true` |
| `CheckExternal` | Enables external reference checking; all tag types. | `true` |
| `CheckInternal` | Enables internal reference checking; all tag types. | `true` |
| `CheckMailto` | Enables–albeit quite basic–`mailto:` link checking. | `true` |
| `CheckTel` | Enables–albeit quite basic–`tel:` link checking. | `true` |
| `EnforceHTTPS` | Fails when encountering an `http://` link. Useful to prevent mixed content errors when serving over HTTPS. | `false` |
| `IgnoreURLs` | Array of strings or regexs of URLs to ignore. | empty |
| `IgnoreDirs` | Array of strings or regexs of directories to ignore when scanning for HTML files. | empty |
| `IgnoreCanonicalBrokenLinks` | When true produces a warning, rather than an error for broken canonical links. When testing a site which isn't live yet or before publishing a new page canonical links will fail. | `true` |
| `IgnoreAltMissing` | Turns off image alt attribute checking. | `false` |
| `IgnoreDirectoryMissingTrailingSlash` | Turns off errors for links to directories without a trailing slash. | `false` |
| `IgnoreTagAttribute` | Specify the ignore attribute. All tags with this attribute will be excluded from every check. | `"data-proofer-ignore"` |
| `TestFilesConcurrently` | :warning: :construction: *EXPERIMENTAL* Turns on concurrent checking of files. | `false` |
| `DocumentConcurrencyLimit` | Maximum number of documents to process at once. | `128` |
| `HTTPConcurrencyLimit` | Maximum number of open HTTP connections. If you raise this number ensure the `ExternalTimeout` is suitably raised. | `4` |
| `LogLevel` | Logging level, 0-3: debug, info, warning, error. | `2` |
| `DirectoryIndex` | The file to look for when linking to a directory. | `index.html` |
| `ExternalTimeout` | Number of seconds to wait on an HTTP connection before failing. | `15` |
| `StripQueryString` | Enables stripping of query strings from external checks. | `true` |
| `StripQueryExcludes` | List of URLs to disable query stripping on. | `["fonts.googleapis.com"]` |
| `ProgDir` | Directory to store cache and log files in. Relative to executing directory. | `tmp/.htmltest` |
| `CacheFile` | File within `ProgDir` to store reference cache. | `refcache.json` |
| `LogFile` | File within `ProgDir` to store last tests errors. | `htmltest.log` |
| `CacheExpires` | Cache validity period, accepts [go.time duration strings](https://golang.org/pkg/time/#ParseDuration) (…"m", "h"). | `336h` (two weeks) |

### Example

```yaml
DirectoryPath: "_site"
EnforceHTTPS: true
IgnoreURLs:
- "example.com"
IgnoreDirs:
- "lib"
CacheExpires: "6h"
```

## :loudspeaker: Issues? Suggestions?

[Submit an issue](https://github.com/wjdp/htmltest/issues/new).
