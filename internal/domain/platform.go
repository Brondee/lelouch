package domain

import "errors"

type Platform string

var ErrInvalidPlatform = errors.New("invalid platform")

const (
	PlatformVinted  Platform = "vinted"
	PlatformMercari Platform = "mercari"
	PlatformAvito   Platform = "avito"
)

var validPlatforms = map[Platform]struct{}{
	PlatformVinted:  {},
	PlatformMercari: {},
	PlatformAvito:   {},
}

func (p Platform) Validate() error {
	_, ok := validPlatforms[p]

	if !ok {
		return ErrInvalidPlatform
	}

	return nil
}
