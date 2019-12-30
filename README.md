# chrome-pdf
Create pdf from html using Google Chrome but with better control of header and footer.

## Introduction

Google chrome can save a html as PDF file, but there are no options to edit the header and footer. This tool can help you overcome these limitations using chrome headless mode through [Chrome Devtools Protocol](https://chromedevtools.github.io/devtools-protocol/).

I mainly need to modify the footer, so I make it easy to set left, middle and right side of the footer. Check the options below.

## Dependencies

Google chrome needs to be installed on your system.

## Usage

```sh
# Example
../chrome-pdf-Linux -i test.html -o test.pdf -w 4.875in -l "Dec. 30, 2019"  -H "<div class='text center' style=\"font-family:'Times New Roman',serif;\">qPCR Summary</div>"

# Usage of .chrome-pdf:
  -B float
    	Bottom margin in inches (default 0.5)
  -F string
    	footer template
  -H string
    	header template (default "<div class='...'></div>")
  -L float
    	Left margin in inches (default 0.75)
  -R float
    	Right margin in inches (default 0.75)
  -T float
    	Top margin in inches (default 0.5)
  -i string
    	the input html file
  -l string
    	string for the left side of the footer
  -m string
    	string for the center of the footer
  -o string
    	the output pdf file (default "output.pdf")
  -p string
    	the name OR the whole path of google chrome in your system (default "google-chrome-stable")
  -r string
    	string for the right side of the footer, default is page number (default "<span class='pageNumber'></span>/<span class='totalPages'></span>")
  -w string
    	footer width, can give 100%, or 6in or 12cm etc in quote (default "90%"). I found the footer is scaled by 75%, so if you want to set margin to 1 inch, set it to 1 * 0.75 = 0.75 inch.
  -h This help.

```

You need to give at least the input html file, the output pdf file name, and the Google Chrome path if different from the default.

There more options in the CDC Page type (https://godoc.org/github.com/mafredri/cdp#Page), and they are all corresponding the original Page options: https://chromedevtools.github.io/devtools-protocol/tot/Page.

## Changes
- v1.1: add page margin settings.

## Notes

1. Paper size is Letter (8.5 in x 11 in).

1. You can add a `@page {margin: 1.5cm 0cm 1.5cm 0cm;}` in the css file or the **style** part of your html to modify the margin. Check the test.html for an example. You can also set this in the command line now, but the `@page` overwrites the command line margin settings.

1. If you cannot the see the footers, possibly your footer margin is too small. 1.5 cm is good.

## Download binaries

You can go to ["**Releases**"](https://github.com/pinbo/chrome-pdf/releases) to download compiled binary file.

You can also install golang on your computer, then run the command below to compile it for your own computer.

`go get -u github.com/mafredri/cdp`

`go build chrome-pdf.go`

## Acknowledgements

Thanks to the **go/cdp** package (https://github.com/mafredri/cdp), I could make this possible with Golang.

Also Thanks to this article, [HTML to PDF Conversion with Headless Chrome using Go](https://medium.com/compass-true-north/go-service-to-convert-web-pages-to-pdf-using-headless-chrome-5fd9ffbae1af).
