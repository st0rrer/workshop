package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/st0rrer/workshop/client/src/config"
	"github.com/st0rrer/workshop/client/src/report"
	"log"
	"os"
	"time"
)

type Client struct {
	userId uuid.UUID
	ctx    context.Context
}

func NewClient(ctx context.Context) (*Client, error) {

	userId, err := uuid.NewRandom()

	return &Client{
		userId: userId,
		ctx:    ctx,
	}, err
}

func (c *Client) Start() error {

	logger := log.New(os.Stderr, "client#Start ", log.LstdFlags)

	logger.Println("start client with GUID ", c.userId)
	r := report.NewReport()

	for {
		select {
		case <-c.ctx.Done():
			return nil
		case <-time.Tick(time.Duration(config.GetConfig().ReportInterval) * time.Second):
			err := r.ReportRecords(c.userId)
			if err != nil {
				logger.Println("error during reporting records", err)
			}
		}
	}
}
