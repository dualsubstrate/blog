package substackbrowser

import "errors"

// ErrAbortedBeforePublish means the user cancelled the CLI before Substack publish
// (Close after paste, q in plain prompts, Ctrl+C, or pick quit).
var ErrAbortedBeforePublish = errors.New("substackbrowser: aborted before publish")
