package device

import (
	"context"
)

func (c *connection) Start(connCtx context.Context) error {
	c.ctx, c.ctxCancel = context.WithCancel(connCtx)

	go c.Read(c.ctx)
	go c.Write(c.ctx)

	return nil
}
