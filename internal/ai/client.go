package ai

import (
	"context"

	openai "github.com/sashabaranov/go-openai"

	"aibuilder/internal/config"
	"aibuilder/internal/logger"

	"github.com/sirupsen/logrus"
)

type Client interface {
	SendMessage(messages []openai.ChatCompletionMessage) (string, error)
}

type aiClient struct {
	api    *openai.Client
	cfg    *config.Config
	logger *logger.Logger
}

func NewClient(cfg *config.Config, logger *logger.Logger) Client {
	return &aiClient{
		api:    openai.NewClient(cfg.APIKey),
		cfg:    cfg,
		logger: logger,
	}
}

func (c *aiClient) SendMessage(messages []openai.ChatCompletionMessage) (string, error) {
	if c.logger.Level >= logrus.DebugLevel {
		c.logger.Debug("Sending messages to AI:")
		for _, msg := range messages {
			c.logger.Debugf("%s: %s", msg.Role, msg.Content)
		}
	}

	resp, err := c.api.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    c.cfg.Model,
			Messages: messages,
		},
	)
	if err != nil {
		return "", err
	}

	if c.logger.Level >= logrus.DebugLevel {
		c.logger.Debugf("Received AI response: %s", resp.Choices[0].Message.Content)
	}

	return resp.Choices[0].Message.Content, nil
}
