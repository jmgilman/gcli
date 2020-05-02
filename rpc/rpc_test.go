package rpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestDial(t *testing.T) {
	conn, err := Dial("fakehost", true)

	assert.Nil(t, err)
	assert.Equal(t, conn.Target(), "fakehost")
}