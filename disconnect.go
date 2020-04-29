package mongohelper

import "context"

// Disconnect closes the client connection with database
func (l *Link) Disconnect() {
	if l.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), l.connTimeout())

		defer cancel()

		if err := l.client.Disconnect(ctx); err != nil {
			l.log("Disconnect", err.Error())
		}
	}
}
