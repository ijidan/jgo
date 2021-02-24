package jchrome

import (
	"context"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/ijidan/jgo/jgo/jlogger"
	"io/ioutil"
	"math"
)

//根据元素截屏
func EleSS(url string, sel string, ssName string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, doElementSS(url, sel, &buf)); err != nil {
		jlogger.Error(err.Error())
		return err
	}
	if err := ioutil.WriteFile(ssName, buf, 0644); err != nil {
		jlogger.Error(err.Error())
		return err
	}
	return nil
}

//全屏截取
func FullSS(url string, quality int64, ssName string) error {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, doFullSS(url, 90, &buf)); err != nil {
		jlogger.Error(err.Error())
		return err
	}
	if err := ioutil.WriteFile(ssName, buf, 0644); err != nil {
		jlogger.Error(err.Error())
		return err
	}
	return nil
}

//执行根据元素截屏
func doElementSS(url string, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(sel, chromedp.ByQuery),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByQuery),
	}
}

//执行全屏截取
func doFullSS(url string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}
			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))
			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
