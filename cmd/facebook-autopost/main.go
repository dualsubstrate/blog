// Command facebook-autopost posts a specific day's post to a Facebook Page.
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dualsubstrate/blog/internal/cliout"
	"github.com/dualsubstrate/blog/internal/facebookautopost"
	"github.com/dualsubstrate/blog/internal/outbound"
	"github.com/dualsubstrate/blog/internal/socialautopost"
	"github.com/dualsubstrate/blog/internal/socialbundle"
	"github.com/dualsubstrate/blog/internal/substackpublishstate"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	root := flag.String("root", ".", "Repo root (content/)")
	postPath := flag.String("post", os.Getenv("SOCIAL_POST"), "Optional bundle path under content/, e.g. mind-infrastructure/2026-05-08-body-thinks-first (disambiguates multiple posts on same date); or env SOCIAL_POST")
	disableIdempotency := flag.Bool("disable-idempotency", false, "Disable URL-based duplicate check against recent Page posts")
	postDate := flag.String("date", "", "Target date in YYYY-MM-DD (required)")
	dryRun := flag.Bool("dry-run", false, "Log what would be posted; do not call Facebook")
	timeout := flag.Duration("http-timeout", 30*time.Second, "HTTP client timeout for Graph API")
	pageID := flag.String("page-id", os.Getenv("FACEBOOK_PAGE_ID"), "Facebook Page ID (or env FACEBOOK_PAGE_ID)")
	token := flag.String("token", os.Getenv("FACEBOOK_PAGE_ACCESS_TOKEN"), "Page access token (or env FACEBOOK_PAGE_ACCESS_TOKEN)")
	askEach := flag.Bool("ask", false, "Prompt publish or skip for each bundle (also SOCIAL_AUTOPOST_ASK=1)")
	noAsk := flag.Bool("no-ask", false, "Never prompt; post every bundle that passes idempotency (CI default)")
	noMarkPublished := flag.Bool("no-mark-published", false, "Do not write social-published in the bundle (also SOCIAL_AUTOPOST_NO_MARK=1)")
	httpRetries := flag.Int("http-retries", envHTTPRetries(facebookautopost.DefaultPostRetries), "Max attempts per Graph call on idempotency check and publish (also FACEBOOK_HTTP_RETRIES)")
	flag.Parse()

	promptEach := socialautopost.PromptEnabled(*askEach, *noAsk)
	noMark := *noMarkPublished || envNoMarkPublished()
	publishTarget := substackpublishstate.TargetFacebook
	if promptEach {
		log.Printf("interactive: will ask publish or skip for each bundle (stdin is a TTY; use -no-ask to disable)")
	}

	if strings.TrimSpace(*postDate) == "" {
		log.Fatal("-date is required (YYYY-MM-DD)")
	}

	if !*dryRun {
		if strings.TrimSpace(*pageID) == "" || strings.TrimSpace(*token) == "" {
			log.Fatal("page-id and token are required unless -dry-run (use flags or FACEBOOK_PAGE_ID and FACEBOOK_PAGE_ACCESS_TOKEN)")
		}
	}

	bundles, err := socialbundle.LoadBundlesForPublishDate(*root, strings.TrimSpace(*postDate), strings.TrimSpace(*postPath))
	if err != nil {
		log.Fatalf("select: %v", err)
	}
	if len(bundles) == 0 {
		cliout.PrintSocialAutopostEmpty(os.Stdout, "Facebook", strings.TrimSpace(*postDate), strings.TrimSpace(*postPath))
		return
	}

	var client *facebookautopost.Client
	if !*dryRun {
		client = facebookautopost.NewClient(*timeout)
		cliout.PrintSocialAutopostStart(os.Stdout, "Facebook", len(bundles), *httpRetries)
	}

	var bundleFailures int

	for i, b := range bundles {
		withImage := strings.TrimSpace(b.FeaturedImagePath) != ""
		if err := socialautopost.ValidateLinkedInTxt(b.Message, withImage); err != nil {
			log.Fatalf("content/%s: %v", b.RelUnderContent, err)
		}

		alreadyMarked, err := socialautopost.MarkerHasTarget(b.BundleDir, publishTarget)
		if err != nil {
			log.Fatalf("content/%s: social-published: %v", b.RelUnderContent, err)
		}
		if alreadyMarked {
			cliout.PrintSocialSkipAlreadyMarked(os.Stdout, publishTarget, b.RelUnderContent)
			continue
		}

		idempotencySkipped := false
		if !*dryRun && !*disableIdempotency {
			idempCtx, idempCancel := context.WithTimeout(context.Background(), outbound.JobBudget(*timeout, *httpRetries, 20*time.Second))
			defer idempCancel()
			already, err := client.RecentlyPostedURLWithRetry(idempCtx, *pageID, *token, b.PostURL, facebookautopost.DefaultFeedScanLimit, *httpRetries)
			if err != nil {
				log.Printf("facebook-autopost: idempotency check failed content/%s: %v", b.RelUnderContent, err)
				bundleFailures++
				continue
			}
			if already {
				log.Printf("facebook-autopost: skipping, URL already posted: %s", b.PostURL)
				idempotencySkipped = true
			}
		}

		if promptEach {
			choice, err := socialautopost.ConfirmPublishItem(socialautopost.ItemPrompt{
				Index:              i,
				Total:              len(bundles),
				RelUnderContent:    b.RelUnderContent,
				PostURL:            b.PostURL,
				Network:            "Facebook",
				WithImage:          withImage,
				DryRun:             *dryRun,
				IdempotencySkipped: idempotencySkipped,
			})
			if err != nil {
				if errors.Is(err, socialautopost.ErrAutopostQuit) {
					log.Printf("stop: %v", err)
					return
				}
				log.Fatalf("prompt: %v", err)
			}
			switch choice {
			case socialautopost.ChoicePublish:
			case socialautopost.ChoiceTagAsPublished:
				if *dryRun {
					log.Printf("dry-run: tag-as-published (no file write)")
				} else {
					recordPublished(b.BundleDir, publishTarget, false, noMark, "tag-as-published")
				}
				continue
			case socialautopost.ChoiceQuit:
				log.Printf("stop: user quit (%d bundle(s) remaining)", len(bundles)-i-1)
				return
			}
		} else if idempotencySkipped {
			if *dryRun {
				log.Printf("dry-run: already on network (idempotency); no file write")
			} else {
				recordPublished(b.BundleDir, publishTarget, false, noMark, "idempotency")
			}
			continue
		}

		if *dryRun {
			socialbundle.PrintDryRunItem(os.Stdout, i, len(bundles), socialbundle.FacebookPostMode, b.PostURL, b.FeaturedImagePath, b.Message)
			continue
		}

		log.Printf("\n\nfacebook-autopost: ------------------------------------------------------------")
		log.Printf("facebook-autopost: using %s", socialbundle.FacebookPostMode)

		pubCtx, pubCancel := context.WithTimeout(context.Background(), outbound.JobBudget(*timeout, *httpRetries, 20*time.Second))
		defer pubCancel()

		publishReq := facebookautopost.PublishRequest{
			PageID:               *pageID,
			AccessToken:          *token,
			PostURL:              b.PostURL,
			FeedLimit:            facebookautopost.DefaultFeedScanLimit,
			MaxAttempts:          *httpRetries,
			CheckFeedBeforeRetry: !*disableIdempotency,
		}
		if withImage {
			log.Printf("facebook-autopost: posting Page photo + caption (image: %s)", b.FeaturedImagePath)
			publishReq.Post = func() error {
				return client.PostPhotoFromFile(pubCtx, *pageID, *token, b.FeaturedImagePath, b.Message)
			}
		} else {
			log.Printf("facebook-autopost: posting link preview (no featured image in bundle)")
			publishReq.Post = func() error {
				return client.PostLink(pubCtx, *pageID, *token, b.Message, b.PostURL)
			}
		}
		if postErr := client.PublishWithRetry(pubCtx, publishReq); postErr != nil {
			log.Printf("facebook-autopost: failed content/%s: %v", b.RelUnderContent, postErr)
			bundleFailures++
			continue
		}
		cliout.PrintSocialPosted(os.Stdout, "Facebook", b.PostURL)
		recordPublished(b.BundleDir, publishTarget, *dryRun, noMark, "posted")
	}

	if bundleFailures > 0 {
		cliout.PrintSocialFailures(os.Stderr, "Facebook", bundleFailures)
		os.Exit(1)
	}
}

func recordPublished(bundleDir, targetKey string, dryRun, noMark bool, reason string) {
	if err := socialautopost.RecordPublished(bundleDir, targetKey, dryRun, noMark); err != nil {
		log.Fatalf("record social-published (%s): %v", reason, err)
	}
	if !dryRun && !noMark {
		cliout.PrintSocialMarkerRecorded(os.Stdout, targetKey, reason)
	}
}

func envHTTPRetries(defaultRetries int) int {
	v := strings.TrimSpace(os.Getenv("FACEBOOK_HTTP_RETRIES"))
	if v == "" {
		return defaultRetries
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return defaultRetries
	}
	return n
}

func envNoMarkPublished() bool {
	v := strings.TrimSpace(os.Getenv("SOCIAL_AUTOPOST_NO_MARK"))
	switch strings.ToLower(v) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
