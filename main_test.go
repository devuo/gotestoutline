package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateOutline(t *testing.T) {
	src := `package app_test

import (
	alias "testing"
)

func TestFirst(a *alias.T) {
	a.Run("First A", func(b *alias.T) {
		b.Run("First A 1", func(c *alias.T) {
			for i := 0; i < 10; i++ {
				c.Run("First A 1 Alpha" + i, func(d *alias.T) {

				})
			}

			c.Run("First A 1 Beta", func(d *alias.T) {

			})
		})

		b.Run("First A 2", func(c *alias.T) {

		})
	})

	a.Run("First B", func(b *alias.T) {
		b.Run("First B 1", func(c *alias.T) {

		})
	})
}

func TestSecond(a *alias.T) {
	a.Run("Second A", func(b *alias.T) {
		b.Run("Second A 1", func(c *cenas.T) {

		})
	})
}

type testing interface {
	Run(string, func (*alias.T) void)
}

func NotATest(t testing) {
	t.Run("Hello World", func (a *alias.T) {
		
	})
}
`

	tests, err := outline(src)
	assert.NoError(t, err)
	assert.Equal(t, []*Test{
		{
			Name:   "TestFirst",
			Type:   TestType,
			Path:   []string{},
			LBrace: 48,
			RBrace: 447,
		},
		{
			Name:   "First A",
			Type:   SubtestType,
			Path:   []string{"TestFirst"},
			LBrace: 83,
			RBrace: 356,
		},
		{
			Name:   "First A 1",
			Type:   SubtestType,
			Path:   []string{"TestFirst", "First A"},
			LBrace: 121,
			RBrace: 305,
		},
		{
			Name:   "",
			Type:   DynamicSubtestType,
			Path:   []string{"TestFirst", "First A", "First A 1"},
			LBrace: 192,
			RBrace: 241,
		},
		{
			Name:   "First A 1 Beta",
			Type:   SubtestType,
			Path:   []string{"TestFirst", "First A", "First A 1"},
			LBrace: 257,
			RBrace: 300,
		},
		{
			Name:   "First A 2",
			Type:   SubtestType,
			Path:   []string{"TestFirst", "First A"},
			LBrace: 315,
			RBrace: 352,
		},
		{
			Name:   "First B",
			Type:   SubtestType,
			Path:   []string{"TestFirst"},
			LBrace: 365,
			RBrace: 444,
		},
		{
			Name:   "First B 1",
			Type:   SubtestType,
			Path:   []string{"TestFirst", "First B"},
			LBrace: 403,
			RBrace: 440,
		},
		{
			Name:   "TestSecond",
			Type:   TestType,
			Path:   []string{},
			LBrace: 449,
			RBrace: 569,
		},
		{
			Name:   "Second A",
			Type:   SubtestType,
			Path:   []string{"TestSecond"},
			LBrace: 485,
			RBrace: 566,
		},
		{
			Name:   "Second A 1",
			Type:   SubtestType,
			Path:   []string{"TestSecond", "Second A"},
			LBrace: 524,
			RBrace: 562,
		},
	}, tests)
}
