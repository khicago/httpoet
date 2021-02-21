package httpoet_test

import (
	"testing"

	"github.com/khicago/httpoet"
	"github.com/stretchr/testify/assert"
)

func TestNewH(t *testing.T) {
	h := httpoet.H{
		"key": "val",
	}
	assert.Equal(t, len(h), 1, "header count error")
	assert.Equal(t, h.CountOf("key"), 1, "header key counter error")
	assert.Equal(t, h["key"], "val", "header value error")
}

func TestNewHs(t *testing.T) {
	h := httpoet.Hs{
		"key": {"val", "val2"},
	}
	assert.Equal(t, 1, len(h), "header count error")
	assert.Equal(t, 2, h.CountOf("key"), "header key counter error")
	assert.Equal(t, []string{"val", "val2"}, h["key"], "header value error")
}

func TestHWithKV(t *testing.T) {
	org := httpoet.H{
		"key": "val",
	}
	abstract := org.WithKV("k1", "v1")
	assert.NotSame(t, org, abstract, "must keep immutable")
	switch h := abstract.(type) {
	case httpoet.H:
		assert.Equal(t, 2, len(h), "header count error")
		assert.Equal(t, 1, h.CountOf("k1"), "header set value failed")
		assert.Equal(t, "val", h["key"], "header set value error")
		assert.Equal(t, "v1", h["k1"], "header set value error")
	default:
		assert.Fail(t, "type error")
	}

	abstractX := abstract.WithKV("key", "val_new")
	assert.NotSame(t, abstract, abstractX, "must keep immutable")
	switch h := abstractX.(type) {
	case httpoet.H:
		assert.Equal(t, "val_new", h["key"], "header set value error")
		assert.Equal(t, "v1", h["k1"], "header set value error")
	default:
		assert.Fail(t, "type error")
	}

	abstractY := abstractX.WithKV("key")
	assert.NotSame(t, abstractX, abstractY, "must keep immutable")
	switch h := abstractY.(type) {
	case httpoet.H:
		assert.Equal(t, 1, len(h), "header remove value failed")
		assert.Equal(t, "", h["key"], "header remove value failed")
	default:
		assert.Fail(t, "type error")
	}
}

func TestHWithKVAppend(t *testing.T) {
	org := httpoet.H{
		"key": "val",
	}
	abstract := org.WithKVAppend("k1", "v1")
	assert.NotSame(t, org, abstract, "must keep immutable")
	switch h := abstract.(type) {
	case httpoet.H:
		assert.Equal(t, 2, len(h), "header count error")
		assert.Equal(t, 1, h.CountOf("k1"), "header set value failed")
		assert.Equal(t, "val", h["key"], "header set value error")
		assert.Equal(t, "v1", h["k1"], "header set value error")
	default:
		assert.Fail(t, "type error")
	}

	abstractX := abstract.WithKVAppend("key")
	assert.NotSame(t, abstract, abstractX, "must keep immutable")
	switch h := abstractX.(type) {
	case httpoet.H:
		assert.Equal(t, 2, len(h), "header append empty value failed, %v", h)
		assert.Equal(t, "val", h["key"], "header append empty value failed, %v", h)
	default:
		assert.Fail(t, "type error")
	}

	abstractY := abstractX.WithKVAppend("key", "xx")
	assert.NotSame(t, abstractX, abstractY, "must keep immutable")
	switch hs := abstractY.(type) {
	case httpoet.Hs:
		assert.Equalf(t, 2, len(hs), "header append value failed, %v", hs)
		assert.Equalf(t, 2, len(hs["key"]), "header append value failed, %v", hs)
		assert.Equal(t, []string{"val", "xx"}, hs["key"], "header append value failed, %v", hs)
	default:
		assert.Fail(t, "type error")
	}
}

func TestHsWithKV(t *testing.T) {
	org := httpoet.Hs{
		"key": {"val", "val_old"},
	}
	abstract := org.WithKV("k1", "v1")
	assert.NotSame(t, org, abstract, "must keep immutable")
	switch hs := abstract.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "header count error")
		assert.Equal(t, 1, hs.CountOf("k1"), "header set value failed")
		assert.Equal(t, org["key"], hs["key"], "header set value error")
		assert.Equal(t, []string{"v1"}, hs["k1"], "header set value error")
	default:
		assert.Fail(t, "type error")
	}

	abstractX := abstract.WithKV("key", "val_new")
	assert.NotSame(t, abstract, abstractX, "must keep immutable")
	switch hs := abstractX.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "header count error")
		assert.Equal(t, 1, hs.CountOf("key"), "header set value failed")
		assert.Equal(t, []string{"val_new"}, hs["key"], "header set value error")
		assert.Equal(t, []string{"v1"}, hs["k1"], "header set value error")
	default:
		assert.Fail(t, "type error")
	}

	abstractY := abstractX.WithKV("key")
	assert.NotSame(t, abstractX, abstractY, "must keep immutable")
	switch hs := abstractY.(type) {
	case httpoet.Hs:
		assert.Equal(t, 1, len(hs), "header remove value failed")
		v, ok := hs["key"]
		assert.Equal(t, false, ok, "header remove value failed")
		assert.Equal(t, []string(nil), v, "header remove value failed") // thus you cannot use the <nil>
	default:
		assert.Fail(t, "type error")
	}
}

func TestHsWithKVAppend(t *testing.T) {
	org := httpoet.Hs{
		"key": {"val", "val_old"},
	}
	abstract := org.WithKVAppend("k1", "v1")
	assert.NotSame(t, org, abstract, "must keep immutable")
	switch hs := abstract.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "header count error")
		assert.Equal(t, 1, hs.CountOf("k1"), "header set value failed")
		assert.Equal(t, org["key"], hs["key"], "header set value error")
		assert.Equal(t, []string{"v1"}, hs["k1"], "header set value error")
	default:
		assert.Fail(t, "type error")
	}

	abstractX := abstract.WithKVAppend("key")
	assert.NotSame(t, abstract, abstractX, "must keep immutable")
	switch hs := abstractX.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "header append empty value failed, %v", hs)
		assert.Equal(t, abstract.CountOf("key"), abstractX.CountOf("key"), "header append empty value failed, %v", hs)
	default:
		assert.Fail(t, "type error")
	}

	abstractY := abstractX.WithKVAppend("key", "xx")
	assert.NotSame(t, abstractX, abstractY, "must keep immutable")
	switch hs := abstractY.(type) {
	case httpoet.Hs:
		assert.Equalf(t, 2, len(hs), "header append value failed, %v", hs)
		assert.Equalf(t, 3, hs.CountOf("key"), "header append value failed, %v", hs)
		assert.Equal(t, []string{"val", "val_old", "xx"}, hs["key"], "header append value failed, %v", hs)
	default:
		assert.Fail(t, "type error")
	}
}

func TestHWithH(t *testing.T) {
	org := httpoet.H{
		"key": "val",
	}
	abstract := org.WithH(httpoet.H{"k1": "v1"})
	assert.NotSame(t, org, abstract, "must keep immutable")
	switch h := abstract.(type) {
	case httpoet.H:
		assert.Equal(t, 2, len(h), "header count error")
		assert.Equal(t, 1, h.CountOf("k1"), "header set h failed")
		assert.Equal(t, "val", h["key"], "header set h error")
		assert.Equal(t, "v1", h["k1"], "header set h error")
	default:
		assert.Fail(t, "type error")
	}

	abstractX := abstract.WithH(httpoet.H{"key": "val_new"})
	assert.NotSame(t, abstract, abstractX, "must keep immutable")
	switch h := abstractX.(type) {
	case httpoet.H:
		assert.Equal(t, "val_new", h["key"], "header set h error")
		assert.Equal(t, "v1", h["k1"], "header set h error")
	default:
		assert.Fail(t, "type error")
	}

	abstractY := abstract.WithH(httpoet.Hs{"key": {"val", "val_old", "val_new"}})
	assert.NotSame(t, abstractX, abstractY, "must keep immutable")
	switch hs := abstractY.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "header with hs failed")
		assert.Equal(t, 3, hs.CountOf("key"), "header with hs failed")
		assert.Equal(t, []string{"val", "val_old", "val_new"}, hs["key"], "header with hs failed")
	default:
		assert.Fail(t, "type error")
	}
}

func TestHWithHAppend(t *testing.T) {
	org := httpoet.H{
		"key": "val",
	}
	abstract := org.WithHAppend(httpoet.H{"k1": "v1"})
	assert.NotSame(t, org, abstract, "must keep immutable")
	switch h := abstract.(type) {
	case httpoet.H:
		assert.Equal(t, 2, len(h), "header count error")
		assert.Equal(t, 1, h.CountOf("k1"), "header append h failed")
		assert.Equal(t, "val", h["key"], "header append h error")
		assert.Equal(t, "v1", h["k1"], "header append h error")
	default:
		assert.Fail(t, "type error")
	}

	abstractX := abstract.WithHAppend(httpoet.H{"key": "val_old"})
	assert.NotSame(t, abstract, abstractX, "must keep immutable")
	switch hs := abstractX.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "header append with h len error, %v", hs)
		assert.Equal(t, 2, hs.CountOf("key"), "header append with h count error, %v", hs)
		assert.Equal(t, []string{"val", "val_old"}, hs["key"], "header append with h value failed, %v", hs)
	default:
		assert.Fail(t, "type error")
	}

	abstractXX := abstract.WithHAppend(httpoet.Hs{"key": {"val_old", "val_new"}})
	assert.NotSame(t, abstract, abstractXX, "must keep immutable")
	switch hs := abstractXX.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "header append with hs len error, %v", hs)
		assert.Equal(t, 3, hs.CountOf("key"), "header append with hs count error, %v", hs)
		assert.Equal(t, []string{"val", "val_old", "val_new"}, hs["key"], "header append with h value failed, %v", hs)
	default:
		assert.Fail(t, "type error")
	}
}

func TestHsWithH(t *testing.T) {
	org := httpoet.Hs{
		"key": {"val"},
	}
	abstract := org.WithH(httpoet.H{"k1": "v1"})
	assert.NotSame(t, org, abstract, "must keep immutable")
	switch hs := abstract.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "hs count error")
		assert.Equal(t, 1, hs.CountOf("k1"), "hs set h failed")
		assert.Equal(t, []string{"val"}, hs["key"], "hs set h error")
		assert.Equal(t, []string{"v1"}, hs["k1"], "hs set h error")
	default:
		assert.Fail(t, "type error")
	}

	abstractX := abstract.WithH(httpoet.H{"key": "val_old"})
	assert.NotSame(t, abstract, abstractX, "must keep immutable")
	switch hs := abstractX.(type) {
	case httpoet.Hs:
		assert.Equal(t, []string{"val_old"}, hs["key"], "hs set h error")
		assert.Equal(t, []string{"v1"}, hs["k1"], "hs set h error")
	default:
		assert.Fail(t, "type error")
	}

	abstractY := abstract.WithH(httpoet.Hs{"key": {"val", "val_old", "val_new"}})
	assert.NotSame(t, abstractX, abstractY, "must keep immutable")
	switch hs := abstractY.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "hs with hs failed")
		assert.Equal(t, 3, hs.CountOf("key"), "hs with hs failed")
		assert.Equal(t, []string{"val", "val_old", "val_new"}, hs["key"], "hs with hs failed")
	default:
		assert.Fail(t, "type error")
	}
}

func TestHWithHsAppend(t *testing.T) {
	org := httpoet.Hs{
		"key": {"val"},
	}
	abstract := org.WithHAppend(httpoet.H{"k1": "v1"})
	assert.NotSame(t, org, abstract, "must keep immutable")
	switch hs := abstract.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "header count error")
		assert.Equal(t, 1, hs.CountOf("k1"), "header append h failed")
		assert.Equal(t, org["key"], hs["key"], "header append h error")
		assert.Equal(t, []string{"v1"}, hs["k1"], "header append h error")
	default:
		assert.Fail(t, "type error")
	}

	abstractX := abstract.WithHAppend(httpoet.H{"key": "val_old"})
	assert.NotSame(t, abstract, abstractX, "must keep immutable")
	switch hs := abstractX.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "header append with h len error, %v", hs)
		assert.Equal(t, 2, hs.CountOf("key"), "header append with h count error, %v", hs)
		assert.Equal(t, []string{"val", "val_old"}, hs["key"], "header append with h value failed, %v", hs)
	default:
		assert.Fail(t, "type error")
	}

	abstractY := abstractX.WithHAppend(httpoet.Hs{"key": {"val_old", "val_new"}})
	assert.NotSame(t, abstractX, abstractY, "must keep immutable")
	switch hs := abstractY.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "header append with hs len error, %v", hs)
		assert.Equal(t, 4, hs.CountOf("key"), "header append with hs count error, %v", hs)
		assert.Equal(t, []string{"val", "val_old", "val_old", "val_new"}, hs["key"], "header append with h value failed, %v", hs)
	default:
		assert.Fail(t, "type error")
	}

	abstractXX := abstract.WithHAppend(httpoet.Hs{"key": {"val_old", "val_new"}})
	assert.NotSame(t, abstract, abstractXX, "must keep immutable")
	switch hs := abstractXX.(type) {
	case httpoet.Hs:
		assert.Equal(t, 2, len(hs), "header append with hs len error, %v", hs)
		assert.Equal(t, 3, hs.CountOf("key"), "header append with hs count error, %v", hs)
		assert.Equal(t, []string{"val", "val_old", "val_new"}, hs["key"], "header append with h value failed, %v", hs)
	default:
		assert.Fail(t, "type error")
	}
}
