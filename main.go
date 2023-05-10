package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "salam")
	})

	e.POST("/pdf/", downloadPDF)
	e.Logger.Fatal(e.Start(":8080"))

}

//
//func downloadPDF(c echo.Context) error {
//	link := c.FormValue("link")
//
//	taskCtx, cancel := chromedp.NewContext(
//		context.Background(),
//		chromedp.WithLogf(log.Printf),
//	)
//	defer cancel()
//
//	var pdfBuffer []byte
//	if err := chromedp.Run(taskCtx, pdfGrabber(link, "body", &pdfBuffer)); err != nil {
//		return err
//	}
//
//	// Generate a unique file name for the PDF
//	fileName := fmt.Sprintf("github_%d.pdf", time.Now().Unix())
//
//	// Write buffer to file
//	err := ioutil.WriteFile(fileName, pdfBuffer, 0644)
//	if err != nil {
//		return err
//	}
//
//	// Generate a URL for the PDF file
//	hostname := c.Request().Host
//	pdfURL := fmt.Sprintf("http://%s/%s", hostname, fileName)
//
//	// Return the PDF download link as an HTML link
//	linkHTML := fmt.Sprintf("<a href=\"%s\">Download PDF</a>", pdfURL)
//	return c.HTML(http.StatusOK, linkHTML)
//}

func downloadPDF(c echo.Context) error {
	token := c.FormValue("token")
	confirm := "1223456"
	if token != confirm {
		return c.String(http.StatusBadRequest, "token is not correct")
	}

	link := c.FormValue("link")

	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	var pdfBuffer []byte
	if err := chromedp.Run(taskCtx, pdfGrabber(link, "body", &pdfBuffer)); err != nil {
		return err
	}

	// Generate a unique file name for the PDF
	fileName := fmt.Sprintf("github_%d.pdf", time.Now().Unix())

	// Write buffer to file
	err := ioutil.WriteFile(fileName, pdfBuffer, 0644)
	if err != nil {
		return err
	}

	// Generate a URL for the PDF file
	hostname := c.Request().Host
	pdfURL := fmt.Sprintf("http://%s/%s", hostname, fileName)

	// Return the PDF download link as an HTML link
	return c.HTML(http.StatusOK, pdfURL)
}

//func pdfGrabber(url string, sel string, res *[]byte) chromedp.Tasks {
//
//	start := time.Now()
//	return chromedp.Tasks{
//		emulation.SetUserAgentOverride("WebScraper 1.0"),
//		chromedp.Navigate(url),
//		chromedp.WaitVisible(`body`, chromedp.ByQuery),
//		chromedp.ActionFunc(func(ctx context.Context) error {
//			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
//			if err != nil {
//				return err
//			}
//			*res = buf
//			fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
//			return nil
//		}),
//	}
//}

func pdfGrabber(url string, sel string, res *[]byte) chromedp.Tasks {

	start := time.Now()
	return chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(1920, 1080, 1.0, false),
		emulation.SetUserAgentOverride("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"),
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
			return nil
		}),
	}
}
