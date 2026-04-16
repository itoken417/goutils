package chromedp

import (
	"context"
	"sync"
	"time"

	cdp "github.com/chromedp/chromedp"
)

var (
	chromeDPMu     sync.Mutex
	chromeDPCtx    context.Context
	chromeDPCancel context.CancelFunc
)

func initChromeDP() {
	chromeDPMu.Lock()
	defer chromeDPMu.Unlock()

	if chromeDPCancel != nil {
		return
	}

	chromeDPCtx, chromeDPCancel = NewChromeDPContext(context.Background(), true)
}

func destroyChromeDP() {
	chromeDPMu.Lock()
	defer chromeDPMu.Unlock()

	if chromeDPCancel == nil {
		return
	}

	chromeDPCancel()
	chromeDPCtx = nil
	chromeDPCancel = nil
}

func Init() {
	initChromeDP()
}

func Destroy() {
	destroyChromeDP()
}

func NewChromeDPContext(parent context.Context, sslCheck bool) (context.Context, context.CancelFunc) {
	opts := append(cdp.DefaultExecAllocatorOptions[:],
		cdp.ExecPath("/usr/bin/google-chrome-stable"),
		cdp.Headless,
		cdp.DisableGPU,
		cdp.NoFirstRun,
		cdp.NoDefaultBrowserCheck,
		cdp.NoSandbox,
		cdp.WindowSize(1280, 800),
		cdp.Flag("homepage", "about:blank"),
		cdp.Flag("allow-insecure-localhost", true),
		cdp.Flag("disable-popup-blocking", true),
		cdp.Flag("disable-dev-shm-usage", true),
	)

	if !sslCheck {
		opts = append(opts, cdp.Flag("ignore-certificate-errors", true))
	}

	allocatorCtx, allocatorCancel := cdp.NewExecAllocator(parent, opts...)
	ctx, cancel := cdp.NewContext(allocatorCtx)

	return ctx, func() {
		cancel()
		allocatorCancel()
	}
}

func NewChromeDPContextWithTimeout(parent context.Context, timeout time.Duration, sslCheck bool) (context.Context, context.CancelFunc) {
	ctx, cancel := NewChromeDPContext(parent, sslCheck)
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, timeout)

	return timeoutCtx, func() {
		timeoutCancel()
		cancel()
	}
}

func RunChromeDP(ctx context.Context, actions ...cdp.Action) error {
	return cdp.Run(ctx, actions...)
}

func GetChromeDPContext() context.Context {
	initChromeDP()

	chromeDPMu.Lock()
	defer chromeDPMu.Unlock()

	return chromeDPCtx
}

func GetNewPage(sslCheck bool) (context.Context, context.CancelFunc) {
	if !sslCheck {
		return NewChromeDPContext(context.Background(), false)
	}

	return cdp.NewContext(GetChromeDPContext())
}

func RunDefaultChromeDP(actions ...cdp.Action) error {
	return RunChromeDP(GetChromeDPContext(), actions...)
}
