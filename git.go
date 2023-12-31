package git_ver

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io"
)

type VerInfo struct {
	Hash plumbing.Hash
	Name string
}

func GetLatestVersion(repo *git.Repository) (*VerInfo, error) {
	// tags
	tagRefs, tagErr := repo.Tags()
	if tagErr != nil {
		return nil, tagErr
	}
	tagMap := make(map[plumbing.Hash]string)
	if err := tagRefs.ForEach(func(t *plumbing.Reference) error {
		tagMap[t.Hash()] = t.Name().Short()
		return nil
	}); err != nil {
		return nil, err
	}

	// annotated tags
	tagObjRefs, tagObjErr := repo.TagObjects()
	if tagObjErr != nil {
		return nil, tagObjErr
	}
	if err := tagObjRefs.ForEach(func(t *object.Tag) error {
		if _, exists := tagMap[t.Hash]; exists { // deleted annotated tags may also have obj
			tagMap[t.Target] = t.Name
		}
		return nil
	}); err != nil {
		return nil, err
	}

	log, logErr := repo.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
	})
	if logErr != nil {
		return nil, logErr
	}
	var latestTag *VerInfo
	_ = log.ForEach(func(obj *object.Commit) error {
		if tag, exists := tagMap[obj.Hash]; exists {
			latestTag = &VerInfo{
				Hash: obj.Hash,
				Name: tag,
			}
			return io.EOF
		}
		return nil
	})
	return latestTag, nil
}
