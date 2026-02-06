package wire

import "github.com/zoobzio/sum"

// RegisterBoundaries creates and registers all wire Boundaries with the service registry.
func RegisterBoundaries(k sum.Key) error {
	if _, err := sum.NewBoundary[UserResponse](k); err != nil {
		return err
	}
	return nil
}
