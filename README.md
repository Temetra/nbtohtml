# nbtohtml - Jupyter notebook HTML converter

## About

A simple external renderer for Jupyter notebook files in Gitea. Considerably faster than running `jupyter nbconvert` on a Raspberry Pi.

## Usage

```powershell
# Powershell example
Get-Content source.ipynb | nbtohtml.exe
```

## Gitea installation

The following assumes you are running Gitea in Docker, with the data volume mapped

```yaml
# Example path
volumes:
  - /media/usb0/gitea/data:/data
```

Build the binary for `linux/arm64` and copy the binary to the Gitea data volume (`data/tools/nbtohtml`)

```powershell
# Powershell example
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -o output
```

Configure Gitea to render `.ipynb` using the binary (`data/gitea/conf/app.ini`)

```ini
[markup.jupyter]
ENABLED = true
FILE_EXTENSIONS = .ipynb
RENDER_COMMAND = "/data/tools/nbtohtml"
IS_INPUT_FILE = false

[markup.sanitizer.jupyter.img]
ALLOW_DATA_URI_IMAGES = true
```

In your Jupyter environment, generate CSS files from Pygments and copy to the data volume (`data/gitea/public/assets/css`)

```bash
pygmentize -S staroffice -f html -a ".markup.jupyter pre" > jupyter-light.css
pygmentize -S lightbulb -f html -a ".markup.jupyter pre" > jupyter-dark.css
```

Wrap the CSS in each file with the following, setting `prefers-color-scheme` to `light` or `dark` as appropriate

```css
@media (prefers-color-scheme: light) {
.markup.jupyter pre code { white-space: pre; }
.markup.jupyter pre .ln { padding-right:16px; }
/* Contents of CSS here */
}
```

Create a custom header template (`data/gitea/templates/custom/header.tmpl`)

```html
<link rel="stylesheet" href="{{AppSubUrl}}/assets/css/jupyter-light.css" />
<link rel="stylesheet" href="{{AppSubUrl}}/assets/css/jupyter-dark.css" />
```
