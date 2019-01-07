package main

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"time"
)
import "github.com/chromedp/chromedp"
import "log"

const (
	//url = "http://orteil.dashnet.org/cookieclicker/"
	url = "http://orteil.dashnet.org/cookieclicker/beta"
)

var (
	clickBigCookie = chromedp.Click("#bigCookie")
	closeNotificationIfPresent = chromedp.ActionFunc(func(cc context.Context, ee cdp.Executor) error {
		closeButtons := []*cdp.Node{}
		chromedp.Nodes("#notes div.close", &closeButtons).Do(cc, ee)
		if len(closeButtons) > 0 {
			currentCloseButton := closeButtons[0]
			chromedp.MouseClickNode(currentCloseButton).Do(cc, ee)
		}
		return nil
	})
	purchaseUpgradeIfAvailable = chromedp.ActionFunc(func(cc context.Context, ee cdp.Executor) error {
		products := []*cdp.Node{}
		chromedp.Nodes("div.upgrade.enabled", &products).Do(cc, ee)
		if len(products) > 0 {
			currentProduct := products[0]
			chromedp.MouseClickNode(currentProduct).Do(cc, ee)
		}
		return nil
	})
	purchaseProductIfAvailable = chromedp.ActionFunc(func(cc context.Context, ee cdp.Executor) error {
		products := []*cdp.Node{}
		chromedp.Nodes("div.product.enabled", &products).Do(cc, ee)
		if len(products) > 0 {
			currentProduct := products[0]
			chromedp.MouseClickNode(currentProduct).Do(cc, ee)
		}
		return nil
	})
	individualTasks = []chromedp.Action{clickBigCookie, purchaseUpgradeIfAvailable, purchaseProductIfAvailable, closeNotificationIfPresent}
)

func main() {
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}
	// Do the bot that you do so well
	err = c.Run(ctxt, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(2 * time.Second),
		chromedp.Click("a.cc_btn_accept_all"),
	})
	for true {
		for _, actionFunction := range individualTasks {
			timeoutableContext, cancelFunc := context.WithTimeout(ctxt, 5*time.Millisecond)
			c.Run(timeoutableContext, chromedp.Tasks{
				actionFunction,
			})
			cancelFunc()
		}
	}
	// Close all the stuff
	if err != nil {
		log.Printf("%v", err)
	}
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Printf("%v", err)
	}
	err = c.Wait()
	if err != nil {
		log.Printf("%v", err)
	}
}
