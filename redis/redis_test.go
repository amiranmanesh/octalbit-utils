package redis

import (
	"context"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock/v3"
)

func TestIncr(t *testing.T) {

	mockConn := redigomock.NewConn()
	mockRP := &Redis{
		Pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockConn, nil
			},
			MaxIdle:     1,
			IdleTimeout: time.Second,
		},
	}

	expectIterate := 10
	incrCmd := mockConn.Command("INCR", "test")
	for i := 1; i <= expectIterate; i++ {
		incrCmd.Expect(int64(i))
	}

	for i := 1; i <= expectIterate; i++ {
		counter, err := mockRP.Incr(context.Background(), "test")
		if err != nil {
			t.Error(err)
		}
		if counter != i {
			t.Error("Incr error: wrong counter!")
		}
	}

	if mockConn.Stats(incrCmd) != expectIterate {
		t.Error("Incr error: mock cmd not used !")
	}
}

func TestGetString(t *testing.T) {

	mockConn := redigomock.NewConn()
	mockRP := &Redis{
		Pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockConn, nil
			},
			MaxIdle:     1,
			IdleTimeout: time.Second,
		},
	}

	expectResult := "gg"
	getCmd := mockConn.Command("GET", "test").Expect(expectResult)

	val, err := mockRP.GetString(context.Background(), "test")
	if err != nil {
		t.Error(err)
	}
	if val != expectResult {
		t.Error("GetString error: wrong value!")
	}

	if mockConn.Stats(getCmd) != 1 {
		t.Error("GetString error: getCmd not used")
	}
}

func TestGetInt(t *testing.T) {

	mockConn := redigomock.NewConn()
	mockRP := &Redis{
		Pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockConn, nil
			},
			MaxIdle:     1,
			IdleTimeout: time.Second,
		},
	}

	expectResult := 13
	getCmd := mockConn.Command("GET", "test").Expect(int64(expectResult))

	val, err := mockRP.GetInt(context.Background(), "test")
	if err != nil {
		t.Error(err)
	}
	if val != expectResult {
		t.Error("GetInt error: wrong value!")
	}

	if mockConn.Stats(getCmd) != 1 {
		t.Error("GetInt error: getCmd not used")
	}
}

func TestSet(t *testing.T) {

	mockConn := redigomock.NewConn()
	mockRP := &Redis{
		Pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockConn, nil
			},
			MaxIdle:     1,
			IdleTimeout: time.Second,
		},
	}

	expectResult := 13
	cmd := mockConn.Command("SET", "test", expectResult)
	err := mockRP.Set(context.Background(), "test", expectResult)
	if err != nil {
		t.Error(err)
	}

	if !cmd.Called() {
		t.Error("Set error: cmd not used")
	}
}

func TestIncrWithTTl(t *testing.T) {

	mockConn := redigomock.NewConn()
	mockRP := &Redis{
		Pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockConn, nil
			},
			MaxIdle:     1,
			IdleTimeout: time.Second,
		},
	}

	var ttl = time.Second * 3
	expectIterate := 10
	expCmd := mockConn.Command("EXPIRE", "test", uint(ttl.Seconds()))
	incrCmd := mockConn.Command("INCR", "test")
	for i := 1; i <= expectIterate; i++ {
		incrCmd.Expect(int64(i))
	}

	for i := 1; i <= expectIterate; i++ {
		counter, err := mockRP.IncrWithTTl(context.Background(), "test", ttl)
		if err != nil {
			t.Error(err)
		}
		if counter != i {
			t.Error("IncrWithTTl error: wrong counter!")
		}
	}

	if mockConn.Stats(incrCmd) != expectIterate ||
		mockConn.Stats(expCmd) != expectIterate {
		t.Error("IncrWithTTl error: mock cmd not used !")
	}
}

func TestSetTTl(t *testing.T) {

	mockConn := redigomock.NewConn()
	mockRP := &Redis{
		Pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockConn, nil
			},
			MaxIdle:     1,
			IdleTimeout: time.Second,
		},
	}

	var ttl = time.Second * 13
	cmd := mockConn.Command("EXPIRE", "test", uint(ttl.Seconds()))

	err := mockRP.SetTTl(context.Background(), "test", ttl)
	if err != nil {
		t.Error(err)
	}

	if !cmd.Called() {
		t.Error("SetTTl error: mock cmd not used !")
	}
}

func TestSetWithTTl(t *testing.T) {

	mockConn := redigomock.NewConn()
	mockRP := &Redis{
		Pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockConn, nil
			},
			MaxIdle:     1,
			IdleTimeout: time.Second,
		},
	}

	var ttl = time.Second * 13
	value := 13
	cmd := mockConn.Command("SETEX", "test", uint(ttl.Seconds()), value)

	err := mockRP.SetWithTTl(context.Background(), "test", value, ttl)
	if err != nil {
		t.Error(err)
	}

	if !cmd.Called() {
		t.Error("SetWithTTl error: mock cmd not used !")
	}
}

func TestDelete(t *testing.T) {

	mockConn := redigomock.NewConn()
	mockRP := &Redis{
		Pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockConn, nil
			},
			MaxIdle:     1,
			IdleTimeout: time.Second,
		},
	}

	cmd1 := mockConn.Command("DEL", "test")
	cmd2 := mockConn.Command("DEL", "test1")
	cmd3 := mockConn.Command("DEL", "test2")
	err := mockRP.Delete(context.Background(), "test", "test1", "test2")
	if err != nil {
		t.Error(err)
	}

	cmd4 := mockConn.Command("DEL", "test3").Expect(0)
	err = mockRP.Delete(context.Background(), "test3")
	if err == nil {
		t.Error("Delete error: expected err but got nil")
	}

	if !cmd1.Called() || !cmd2.Called() || !cmd3.Called() || !cmd4.Called() {
		t.Error("Delete error: mock cmd not used !")
	}
}

func TestDeleteByPattern(t *testing.T) {

	mockConn := redigomock.NewConn()
	mockRP := &Redis{
		Pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockConn, nil
			},
			MaxIdle:     1,
			IdleTimeout: time.Second,
		},
	}

	pattern := "gg*"
	expected := []string{"Penn", "Teller"}
	cmdKeys := mockConn.Command("KEYS", pattern).ExpectStringSlice(expected...)
	cmdDel1 := mockConn.Command("DEL", expected[0])
	cmdDel2 := mockConn.Command("DEL", expected[1])

	err := mockRP.DeleteByPattern(context.Background(), pattern)
	if err != nil {
		t.Error(err)
	}

	if !cmdKeys.Called() || !cmdDel1.Called() || !cmdDel2.Called() {
		t.Error("DeleteByPattern error: mock cmds not used !")
	}
}

func TestAddToSet(t *testing.T) {

	mockConn := redigomock.NewConn()
	mockRP := &Redis{
		Pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockConn, nil
			},
			MaxIdle:     1,
			IdleTimeout: time.Second,
		},
	}

	cmd := mockConn.Command("SADD", "test", "some value")
	cmd1 := mockConn.Command("SCARD", "test").Expect(int64(1))
	resp, err := mockRP.AddToSet(context.Background(), "test", "some value")
	if err != nil {
		t.Error(err)
	}

	if resp != 1 {
		t.Error("AddToSet error: expected 1")
	}

	if !cmd.Called() || !cmd1.Called() {
		t.Error("AddToSet error: cmd not used")
	}
}

func TestAddToSetWithTTl(t *testing.T) {

	mockConn := redigomock.NewConn()
	mockRP := &Redis{
		Pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockConn, nil
			},
			MaxIdle:     1,
			IdleTimeout: time.Second,
		},
	}

	ttl := time.Second * 3
	cmd1 := mockConn.Command("SADD", "test", "some value")
	cmd2 := mockConn.Command("SCARD", "test").Expect(int64(1))
	cmd3 := mockConn.Command("EXPIRE", "test", uint(ttl.Seconds()))
	resp, err := mockRP.AddToSetWithTTl(context.Background(), "test", "some value", ttl)
	if err != nil {
		t.Error(err)
	}

	if resp != 1 {
		t.Error("AddToSetWithTTl error: expected 1")
	}

	if !cmd1.Called() || !cmd2.Called() || !cmd3.Called() {
		t.Error("AddToSetWithTTl error: cmds not used")
	}
}
