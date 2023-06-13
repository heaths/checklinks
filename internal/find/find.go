// Copyright 2023 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package find

// cspell:ignore fsys

import (
	"context"
	"io/fs"
	"regexp"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/heaths/checklinks/internal/log"
)

var (
	isURL = regexp.MustCompile(`https?://[\w\.:\/%~_\-\+]+`)
)

type Match struct {
	URL  string
	Path string
}

func Find(ctx context.Context, fsys fs.FS, patterns []string) <-chan Match {
	ctx, cancel := context.WithCancel(ctx)
	matches := make(chan Match)
	wg := new(sync.WaitGroup)

	go func() {
		err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.Type().IsRegular() {
				return nil
			}

			var matched bool
			for _, pattern := range patterns {
				if matched, err = doublestar.PathMatch(pattern, p); err != nil {
					return err
				} else if matched {
					break
				}
			}

			if !matched {
				return nil
			}

			select {
			case <-ctx.Done():
				log.Debug("scanning canceled")
				return fs.SkipAll
			default:
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				log.Verbose("scanning %s", p)
				buf, err := fs.ReadFile(fsys, p)
				if err != nil {
					log.Debug("read %s: %s", p, err)
					return
				}

				for _, found := range isURL.FindAll(buf, len(buf)) {
					matches <- Match{
						URL:  string(found),
						Path: p,
					}
				}
			}()

			return nil
		})

		if err != nil {
			cancel()
		}

		wg.Wait()
		close(matches)
	}()

	return matches
}
