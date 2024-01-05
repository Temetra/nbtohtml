# nbtohtml - Jupyter notebook HTML converter

## Usage

```powershell
# Powershell example
Get-Content source.ipynb | nbtohtml.exe
```

## Gitea installation

Build

```powershell
# Powershell example
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -o output
```

Copy binary to Gitea data volume

```
/media/usb0/gitea/data/tools/nbtohtml
```

Config Gitea (`gitea\conf\app.ini`)

```ini
[markup.jupyter]
ENABLED = true
FILE_EXTENSIONS = .ipynb
RENDER_COMMAND = "/data/tools/nbtohtml"
IS_INPUT_FILE = false
```

Generate CSS files from Pygments (in your Jupyter environment) (`gitea\public\css`)

```bash
pygmentize -S staroffice -f html -a ".markup.jupyter pre" > jupyter-light.css
pygmentize -S lightbulb -f html -a ".markup.jupyter pre" > jupyter-dark.css
```

Wrap CSS with following code, setting `prefers-color-scheme` to `light` or `dark` as appropriate

```css
@media (prefers-color-scheme: light) {
.markup.jupyter pre code { white-space: pre; }
.markup.jupyter pre .ln { padding-right:16px; }
/* Contents of CSS here */
}
```

Create custom header template (`gitea\templates\custom\header.tmpl`)

```html
<link rel="stylesheet" href="{{AppSubUrl}}/assets/css/jupyter-light.css" />
<link rel="stylesheet" href="{{AppSubUrl}}/assets/css/jupyter-dark.css" />
```
