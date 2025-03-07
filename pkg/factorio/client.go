package factorio

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorcon/rcon"
)

var _ error = CommandError{}

type CommandError struct {
	CMD string
	ErrorMsg string
}

func (e CommandError) Error() string {
	return "error by cmd: " + e.ErrorMsg + "\ncmd: " + e.CMD
}

func NewCommandError(cmd, resp string) CommandError {
	index := strings.Index(resp, "Error: ")
	if index > -1 {
		msg := strings.TrimSpace(resp[index + 6:])
		return CommandError{
			CMD: cmd,
			ErrorMsg: msg,
		}
	}

	return CommandError{
		CMD: cmd,
		ErrorMsg: resp,
	}
}


type RCONClient struct {
	conn   *rcon.Conn
	logger *log.Logger
}

func NewRCONClient(address, password string) (*RCONClient, error) {
	logger := log.WithPrefix("rcon")
	logger.Debug("Creating new RCON client", "address", address, "password", password)

	conn, err := rcon.Dial(address, password)
	if err != nil {
		return nil, err
	}

	return &RCONClient{conn, logger}, nil
}

func (c *RCONClient) Execute(cmd string) (string, error) {
	logger := c.logger.With("reqid", uuid.NewString())
	logger.Debug("Sending command", "command", cmd)

	resp, err := c.conn.Execute(cmd)
	if err != nil {
		logger.Error("Failed to execute command", "error", err)
		return "", err
	}
	if strings.Contains(resp, "Cannot execute command. Error: ") {
		return "", CommandError{
			CMD: cmd,
			ErrorMsg: resp,
		}
	}

	logger.Debug("Received response", "response", resp)

	return resp, nil
}

func (c *RCONClient) Close() error {
	return c.conn.Close()
}
