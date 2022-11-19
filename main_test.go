package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateOutline(t *testing.T) {
	src := `package testdata

import (
	cenas "testing"
)

func TestFirst(a *cenas.T) {
	a.Run("First A", func(b *cenas.T) {
		b.Run("First A 1", func(c *cenas.T) {
			for i := 0; i < 10; i++ {
				c.Run("First A 1 Alpha" + i, func(d *cenas.T) {

				})
			}

			c.Run("First A 1 Beta", func(d *cenas.T) {

			})
		})

		b.Run("First A 2", func(c *cenas.T) {

		})
	})

	a.Run("First B", func(b *cenas.T) {
		b.Run("First B 1", func(c *cenas.T) {

		})
	})
}

func TestSecond(a *cenas.T) {
	a.Run("Second A", func(b *cenas.T) {
		b.Run("Second A 1", func(c *cenas.T) {

		})
	})
}`

	tests, err := generateOutline(src)
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
