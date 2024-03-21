# rdr

`rdr` is a command-line tool that converts web pages to Markdown format and displays the content in a pager or terminal. It uses the `go-readability` library to extract the main content from a web page and the `html-to-markdown` library to convert the HTML to Markdown.

## Usage

```
rdr [url]
```

If a URL is provided as an argument, `rdr` will fetch the content from that URL. If no URL is provided, `rdr` will read the input from stdin.

## Options

- `--no-pager`: Don't pipe the output to a pager (e.g., `less`). Instead, print the output directly to the terminal.
- `--no-links`: Don't display any links in the output. Links will be replaced with their text content.
- `--raw`: Print the raw Markdown output without any formatting or rendering.

## Examples

1. Convert a web page to Markdown and display it in a pager:

```
rdr https://example.com
```

2. Convert a local HTML file to Markdown and display it in the terminal:

```
cat index.html | rdr -
```

3. Convert a web page to Markdown, remove all links, and display it in the terminal:

```
rdr --no-links --no-pager https://example.com
```

4. Convert a web page to raw Markdown and print it to stdout:

```
rdr --raw https://example.com
```

## Dependencies

- `github.com/JohannesKaufmann/html-to-markdown`: Library for converting HTML to Markdown.
- `github.com/PuerkitoBio/goquery`: Library for parsing HTML.
- `github.com/charmbracelet/glamour`: Library for rendering Markdown with syntax highlighting.
- `github.com/go-shiori/go-readability`: Library for extracting the main content from a web page.
- `github.com/spf13/cobra`: Library for creating command-line applications.
