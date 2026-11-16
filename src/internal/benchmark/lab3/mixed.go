package lab3

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// RunMixed runs the read and write workloads concurrently. The caller must
// provide a pool that permits at least two simultaneous connections.
func RunMixed(ctx context.Context, read, write func(context.Context) error) error {
	group, ctx := errgroup.WithContext(ctx)
	start := make(chan struct{})
	group.Go(func() error {
		<-start
		if err := read(ctx); err != nil {
			return fmt.Errorf("mixed read: %w", err)
		}
		return nil
	})
	group.Go(func() error {
		<-start
		if err := write(ctx); err != nil {
			return fmt.Errorf("mixed write: %w", err)
		}
		return nil
	})
	close(start)
	return group.Wait()
}
