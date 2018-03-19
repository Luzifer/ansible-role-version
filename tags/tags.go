package tags

import (
	"errors"
	"io"
	"sort"
	"time"

	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type Tag struct {
	Name string
	When time.Time
}

var (
	ErrNoTagsFound = errors.New("No tags found")
)

// GetLatestTag clones a Git repository into memory and resolves latest
// leightweight or annotated tag from it
func GetLatestTag(repoURL string, includeLightweight bool) (*Tag, error) {
	fs := memfs.New()
	// Git objects storer based on memory
	storer := memory.NewStorage()

	// Clones the repository into the worktree (fs) and storer all the .git
	// content into the storer
	r, _ := git.Clone(storer, fs, &git.CloneOptions{
		URL: repoURL,
	})

	tags := []Tag{}

	// Get reference iterator for all tags
	it, err := r.Tags()
	if err != nil {
		return nil, err
	}
	defer it.Close()

	for {
		ref, err := it.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if !ref.Name().IsTag() {
			continue
		}

		var when time.Time

		// Check whether the hash is resolvable to a tag
		if t := tryTag(r, ref.Hash()); t != nil {
			when = *t
		}

		// If it wasn't and we may include leightweight tags check for commit
		if when.IsZero() && includeLightweight {
			if t := tryCommit(r, ref.Hash()); t != nil {
				when = *t
			}
		}

		if when.IsZero() {
			// We've resolved no tag
			continue
		}

		tags = append(tags, Tag{
			Name: ref.Name().Short(),
			When: when,
		})
	}

	if len(tags) == 0 {
		return nil, ErrNoTagsFound
	}

	// Tags may be unsorted by design so we need to sort them
	sort.SliceStable(tags, func(i, j int) bool {
		return tags[i].When.Before(tags[j].When)
	})

	return &tags[len(tags)-1], nil
}

func tryTag(r *git.Repository, hash plumbing.Hash) *time.Time {
	tag, err := r.TagObject(hash)
	if err != nil {
		return nil
	}

	return &tag.Tagger.When
}

func tryCommit(r *git.Repository, hash plumbing.Hash) *time.Time {
	commit, err := r.CommitObject(hash)
	if err != nil {
		return nil
	}

	return &commit.Committer.When
}
