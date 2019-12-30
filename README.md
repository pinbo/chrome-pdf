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
  -F string
    	footer template
  -H string
    	header template (default "<div class='...'></div>")
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
There more options in the CDC Page type (https://godoc.org/github.com/mafredri/cdp#Page), and they are all corresponding the original Page options: https://chromedevtools.github.io/devtools-protocol/tot/Page.

## Download binaries

You can go to ["**Releases**"](https://github.com/pinbo/chrome-pdf/releases) to download compiled binary file.

You can also install golang on your computer, then run the command below to compile it for your own computer.

`go build chrome-pdf.go`
