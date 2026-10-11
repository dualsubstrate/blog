package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/dualsubstrate/blog/internal/cliout"
	"github.com/dualsubstrate/blog/internal/substackpublishstate"
)

// publishQueueEmptyError is returned when pick-draft finds no unpublished bundles.
type publishQueueEmptyError struct {
	contentRoot string
	marker      string
	target      string
}

func (e *publishQueueEmptyError) Error() string {
	return fmt.Sprintf("publish queue empty for target %q", e.target)
}

func newPublishQueueEmptyError(contentRoot, marker, target string) error {
	return &publishQueueEmptyError{
		contentRoot: strings.TrimSpace(contentRoot),
		marker:      strings.TrimSpace(marker),
		target:      strings.TrimSpace(target),
	}
}

func substackQueueHints(target string) cliout.PublishQueueHints {
	target = strings.TrimSpace(target)
	switch target {
	case substackpublishstate.TargetSubstackES:
		return cliout.PublishQueueHints{
			LangLabel: "Spanish Substack",
			ListCmd:   "make sb-list-unpublished PUBLISH_TARGET=substack-es",
			PickCmd:   "make sb-es-pick-publish",
			PasteCmd:  "make sb-es section/slug",
			MarkHint:  "make sb-mark-published POST=section/slug PUBLISH_TARGET=substack-es",
		}
	default:
		return cliout.PublishQueueHints{
			LangLabel: "English Substack",
			ListCmd:   "make sb-list-unpublished",
			PickCmd:   "make sb-en-pick",
			PasteCmd:  "make sb-en section/slug",
			MarkHint:  "make sb-mark-published POST=section/slug PUBLISH_TARGET=substack-en",
		}
	}
}

func printPublishQueueEmpty(out io.Writer, contentRoot, target, marker string) {
	if strings.TrimSpace(contentRoot) == "" {
		contentRoot = "content"
	}
	if strings.TrimSpace(marker) == "" {
		marker = substackpublishstate.DefaultMarker
	}
	cliout.PrintPublishQueueEmpty(out, contentRoot, target, marker, substackQueueHints(target))
}

func printPublishQueueList(out io.Writer, list []substackpublishstate.BundlePostPath, target, marker string) {
	if strings.TrimSpace(marker) == "" {
		marker = substackpublishstate.DefaultMarker
	}
	paths := make([]string, len(list))
	for i, p := range list {
		paths[i] = string(p)
	}
	cliout.PrintPublishQueueList(out, paths, target, marker, substackQueueHints(target))
}

func handlePublishQueueError(out io.Writer, err error) bool {
	var empty *publishQueueEmptyError
	if errors.As(err, &empty) {
		printPublishQueueEmpty(out, empty.contentRoot, empty.target, empty.marker)
		return true
	}
	return false
}
