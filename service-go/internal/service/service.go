package service

import "context"

type Service interface {
	// Start will prepare and eventually run the service that might be leading up to blocking any interaction.
	Start(ctx context.Context) error

	// Stop will immediately signaling the service to initiate a graceful shutdown. While it can takes time,
	// Stop might be leading up to another blocking any interaction, it is wiser to wait for the Stop to end
	// naturally rather than dropping of in the middle of cleanup process.
	Stop(ctx context.Context) error
}
