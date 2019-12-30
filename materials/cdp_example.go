package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"os/exec"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	//"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
)

func main() {
	cmd := exec.Command("google-chrome-stable", "--headless --remote-debugging-port=9222")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err := run(5 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

func run(timeout time.Duration) error {
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
	navArgs := page.NewNavigateArgs("file:///home/junli/Dropbox/Mix/JobApplication/20190401_VT_breeder/Research-statement-JZ.html")
	nav, err := c.Page.Navigate(ctx, navArgs)
	if err != nil {
		return err
	}

	// Wait until we have a DOMContentEventFired event.
	if _, err = domContent.Recv(); err != nil {
		return err
	}

	fmt.Printf("Page loaded with frame ID: %s\n", nav.FrameID)

	// Fetch the document root node. We can pass nil here
	// since this method only takes optional arguments.
	// doc, err := c.DOM.GetDocument(ctx, nil)
	// if err != nil {
	// 	return err
	// }

	// Get the outer HTML for the page.
	// result, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
	// 	NodeID: &doc.Root.NodeID,
	// })
	// if err != nil {
	// 	return err
	// }

	//fmt.Printf("HTML: %s\n", result.OuterHTML)

	// Capture a screenshot of the current page.
	//screenshotName := "screenshot.jpg"
	// screenshotArgs := page.NewCaptureScreenshotArgs().
	// 	SetFormat("jpeg").
	// 	SetQuality(80)
	// screenshot, err := c.Page.CaptureScreenshot(ctx, screenshotArgs)
	// if err != nil {
	// 	return err
	// }
	// if err = ioutil.WriteFile(screenshotName, screenshot.Data, 0644); err != nil {
	// 	return err
	// }

	// fmt.Printf("Saved screenshot: %s\n", screenshotName)

	// // Wait for the page to finish loading
	// _, _ := loadEventFiredClient.Recv()
	
	// Print to PDF
	outpdf := "test.pdf"
	printToPDFArgs := page.NewPrintToPDFArgs().
		// SetPrintBackground(true).
		// SetMarginTop(0).
		// SetMarginBottom(0).
		// SetMarginLeft(0).
		// SetMarginRight(0).
		SetDisplayHeaderFooter(true).
		SetHeaderTemplate(`<div class='...'></div>`).
		SetFooterTemplate(`<div style="font-size:10px; font-family:'Times New Roman',serif; width:12.4cm; margin:0 auto;">
    <p style="float:left; text-align:left; width:20%;"></p>
    <p style="float:left; text-align:center; width:60%;">Junli Zhang - Dec 27, 2019</p>
    <p style="float:left; text-align:right; width:20%;"><span class='pageNumber'></span>/<span class='totalPages'></span></p>
    </div>`)
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