package main

// need cdp package, install by
//  go get -u github.com/mafredri/cdp
import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	//"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
)

func main() {
	// command line options
	chromePath := flag.String("p", "google-chrome-stable", "the name OR the whole path of google chrome in your system")
	inputHtml := flag.String("i", "", "the input html file")
	outputPDF := flag.String("o", "output.pdf", "the output pdf file")
	footerLeft := flag.String("l", "", "string for the left side of the footer")
	footerRight := flag.String("r", "<span class='pageNumber'></span>/<span class='totalPages'></span>", "string for the right side of the footer, default is page number")
	footerMiddle := flag.String("m", "", "string for the center of the footer")
	headerTemplate := flag.String("H", "<div class='...'></div>", "header template")
	footerTemplate := flag.String("F", "", "footer template")
	footerWidth := flag.String("w", "90%", "footer width, can give 100%, or 6in or 12cm etc in quote")
	flag.Parse()
	if *inputHtml == "" {
		fmt.Println("Please give -i input html file!")
		os.Exit(2)
	}
	if *footerTemplate == "" {
		*footerTemplate = `<div style="font-size:10px; font-family:'Times New Roman',serif; width:` + *footerWidth + `; margin:0 auto;">
		<p style="float:left; text-align:left; width:25%;">` + *footerLeft + `</p>
		<p style="float:left; text-align:center; width:50%;">` + *footerMiddle + `</p>
		<p style="float:left; text-align:right; width:25%;">` + *footerRight + `</p>
		</div>`
	}
	// get absolute path of inputHtml
	abs, err := filepath.Abs(*inputHtml)
	if err != nil {
		log.Fatal(err)
	}

	// start chrome headless mode
	cmd := exec.Command(*chromePath, "--headless", "--remote-debugging-port=9222", "https://www.google.com/")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// wait for chrome to launch
	duration := 1 * time.Second
	time.Sleep(duration)
	// save to pdf
	err = run(5*time.Second, abs, *outputPDF, *headerTemplate, *footerTemplate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PDF created.")
	// Kill headless mode:
	if err := cmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill process: ", err)
	}
}

func run(timeout time.Duration, input string, outpdf string, headerTemplate string, footerTemplate string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devt := devtool.New("http://127.0.0.1:9222")
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		if err != nil {
			return err
		}
	}

	// Initiate a new RPC connection to the Chrome DevTools Protocol target.
	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		return err
	}
	defer conn.Close() // Leaving connections open will leak memory.

	c := cdp.NewClient(conn)

	// Open a DOMContentEventFired client to buffer this event.
	domContent, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		return err
	}
	defer domContent.Close()

	// Enable events on the Page domain, it's often preferrable to create
	// event clients before enabling events so that we don't miss any.
	if err = c.Page.Enable(ctx); err != nil {
		return err
	}

	// Create the Navigate arguments with the optional Referrer field set.
	//navArgs := page.NewNavigateArgs("https://www.google.com").
	//	SetReferrer("https://duckduckgo.com")
	navArgs := page.NewNavigateArgs("file://" + input)
	nav, err := c.Page.Navigate(ctx, navArgs)
	if err != nil {
		return err
	}

	// Wait until we have a DOMContentEventFired event.
	if _, err = domContent.Recv(); err != nil {
		return err
	}

	fmt.Printf("Page loaded with frame ID: %s\n", nav.FrameID)

	// Print to PDF
	//outpdf := "test.pdf"
	printToPDFArgs := page.NewPrintToPDFArgs().
		// SetPrintBackground(true).
		// SetMarginTop(0).
		// SetMarginBottom(0).
		// SetMarginLeft(0).
		// SetMarginRight(0).
		SetDisplayHeaderFooter(true).
		SetHeaderTemplate(headerTemplate).
		SetFooterTemplate(footerTemplate)
	print, err := c.Page.PrintToPDF(ctx, printToPDFArgs)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(outpdf, print.Data, 0644); err != nil {
		return err
	}

	fmt.Printf("Saved pdf: %s\n", outpdf)

	return nil
}
