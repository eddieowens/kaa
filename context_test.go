package kaa

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"testing"
)

type myCmdPayload struct {
	Username    string   `arg:"0"`
	Id          uint     `arg:"1"`
	Height      *int     `arg:"2,optional"`
	Env         string   `flag:"env"`
	StringSlice []string `flag:"string_slice"`
	Boolean     bool     `flag:"boolean"`
}

type ContextTest struct {
	suite.Suite
}

func baseCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().String("env", "", "")
	cmd.Flags().Bool("boolean", false, "")
	cmd.Flags().StringSlice("string_slice", nil, "")
	return cmd
}

func (c *ContextTest) TestNewContextFromData() {
	p := new(myCmdPayload)

	args := []string{
		"hello", "2", "3",
	}

	cmd := baseCmd()
	_ = cmd.Flags().Set("env", "dev")
	_ = cmd.Flags().Set("boolean", "true")
	_ = cmd.Flags().Set("string_slice", "1,2,3")

	expectedPayload := &myCmdPayload{
		Username:    "hello",
		Id:          2,
		Height:      nil,
		Env:         "dev",
		Boolean:     true,
		StringSlice: []string{"1", "2", "3"},
	}

	ctx := NewContext(cmd, args)

	if c.NotNil(ctx) {
		if c.NoError(ctx.Bind(p)) {
			c.Equal(expectedPayload, p)
		}
	}
}

func (c *ContextTest) TestNewContextFromDataWithOptional() {
	p := new(myCmdPayload)

	args := []string{
		"hello", "2",
	}

	cmd := baseCmd()
	_ = cmd.Flags().Set("env", "dev")
	_ = cmd.Flags().Set("boolean", "true")
	_ = cmd.Flags().Set("string_slice", "1,2,3")

	expectedPayload := &myCmdPayload{
		Username:    "hello",
		Id:          2,
		Height:      nil,
		Env:         "dev",
		Boolean:     true,
		StringSlice: []string{"1", "2", "3"},
	}

	ctx := NewContext(cmd, args)

	if c.NotNil(ctx) {
		if c.NoError(ctx.Bind(p)) {
			c.Equal(expectedPayload, p)
		}
	}
}

func TestContextTest(t *testing.T) {
	suite.Run(t, new(ContextTest))
}
