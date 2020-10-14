package cmdutil

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/Shopify/themekit/src/colors"
	"github.com/Shopify/themekit/src/file"
)

type cmdSummary struct {
	actions, download, upload, skipped, removed int32
	errors                                      []string
}

func (sum *cmdSummary) completeOp(op file.Op) {
	atomic.AddInt32(&sum.actions, 1)
	switch op {
	case file.Update:
		atomic.AddInt32(&sum.upload, 1)
	case file.Skip:
		atomic.AddInt32(&sum.skipped, 1)
	case file.Remove:
		atomic.AddInt32(&sum.removed, 1)
	case file.Get:
		atomic.AddInt32(&sum.download, 1)
	}
}

func (sum *cmdSummary) err(errStr string) {
	sum.errors = append(sum.errors, errStr)
}

func (sum *cmdSummary) display(ctx *Ctx) {
	if sum.actions == 0 {
		return
	}
	var results = []string{fmt.Sprintf("%v files", sum.actions)}
	if sum.download > 0 {
		results = append(results, fmt.Sprintf("%v: %v", colors.Blue("Downloaded"), sum.download))
	}
	if sum.upload > 0 {
		results = append(results, fmt.Sprintf("%v: %v", colors.Green("Updated"), sum.upload))
	}
	if sum.removed > 0 {
		results = append(results, fmt.Sprintf("%v: %v", colors.Yellow("Removed"), sum.removed))
	}
	if sum.skipped > 0 {
		results = append(results, fmt.Sprintf("%v: %v", colors.Cyan("No Change"), sum.skipped))
	}
	if len(sum.errors) > 0 {
		results = append(results, fmt.Sprintf("%v: %v", colors.Red("Errored"), len(sum.errors)))
	}
	ctx.Log.Printf("[%v] %v", colors.Green(ctx.Env.Name), strings.Join(results, ", "))
	if len(sum.errors) > 0 {
		ctx.ErrLog.Println(colors.Red("Errors encountered: "))
		for _, msg := range sum.errors {
			ctx.ErrLog.Printf("\t%v", msg)
		}
		ctx.ErrLog.Println(colors.Red("finished command with errors"))
	}
}
