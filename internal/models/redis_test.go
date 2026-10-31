package models

import (
	"errors"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/gomodule/redigo/redis"
)

func testRedis(t *testing.T) (*redis.Pool, string) {
	t.Helper()

	url := os.Getenv("REDIS_URL")
	if url == "" {
		t.Skip("REDIS_URL is not set")
	}

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(url)
		},
	}

	prefix := "homettp-test:" + t.Name() + ":"

	t.Cleanup(func() {
		conn := pool.Get()
		defer conn.Close()

		keys, _ := redis.Strings(conn.Do("KEYS", prefix+"*"))
		if len(keys) > 0 {
			_, _ = conn.Do("DEL", redis.Args{}.AddFlat(keys)...)
		}

		pool.Close()
	})

	return pool, prefix
}

func TestRedisUserParallelCreates(t *testing.T) {
	pool, prefix := testRedis(t)
	repository := &RedisUserRepository{RedisPool: pool, RedisKeyPrefix: prefix}

	var mu sync.Mutex
	var wg sync.WaitGroup
	created := 0

	for range 10 {
		wg.Go(func() {
			err := repository.Create(&User{Username: "bob", Email: "bob@example.com", Password: []byte("password1")})
			if err != nil && !errors.Is(err, ErrDuplicateUsername) {
				t.Error(err)
			}

			if err == nil {
				mu.Lock()
				created++
				mu.Unlock()
			}
		})
	}

	wg.Wait()

	if created != 1 {
		t.Fatalf("expected 1 created user, got %d", created)
	}
}

func TestRedisUserUpdateAfterDelete(t *testing.T) {
	pool, prefix := testRedis(t)
	repository := &RedisUserRepository{RedisPool: pool, RedisKeyPrefix: prefix}

	user := &User{Username: "alice", Email: "alice@example.com", Password: []byte("password1")}

	err := repository.Create(user)
	if err != nil {
		t.Fatal(err)
	}

	loaded := *user

	err = repository.Delete(user)
	if err != nil {
		t.Fatal(err)
	}

	err = repository.Update(&loaded, &User{Username: "alice2", Email: "alice@example.com"})
	if !errors.Is(err, ErrNoRecord) {
		t.Fatalf("expected %v, got %v", ErrNoRecord, err)
	}

	err = repository.UpdateRememberToken(&loaded, "token")
	if !errors.Is(err, ErrNoRecord) {
		t.Fatalf("expected %v, got %v", ErrNoRecord, err)
	}

	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", prefix+userKeyPrefix+strconv.Itoa(user.ID)))
	if err != nil {
		t.Fatal(err)
	}

	if exists {
		t.Fatal("expected the deleted user to stay deleted")
	}
}

func TestRedisUserConcurrentUpdates(t *testing.T) {
	pool, prefix := testRedis(t)
	repository := &RedisUserRepository{RedisPool: pool, RedisKeyPrefix: prefix}

	user := &User{Username: "alice", Email: "alice@example.com", Password: []byte("password1")}

	err := repository.Create(user)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	for i := range 5 {
		loaded := *user

		wg.Go(func() {
			name := "user" + strconv.Itoa(i)

			err := repository.Update(&loaded, &User{Username: name, Email: name + "@example.com"})
			if err != nil {
				t.Error(err)
			}
		})
	}

	wg.Wait()

	current, err := repository.Find(user.ID)
	if err != nil {
		t.Fatal(err)
	}

	conn := pool.Get()
	defer conn.Close()

	for key, value := range map[string]string{userUsernameKey: current.Username, userEmailKey: current.Email} {
		values, err := redis.Strings(conn.Do("ZRANGE", prefix+userKeyPrefix+key, 0, -1))
		if err != nil {
			t.Fatal(err)
		}

		if len(values) != 1 || values[0] != value {
			t.Fatalf("expected %s index [%s], got %v", key, value, values)
		}
	}
}

func TestRedisCommandUpdateAfterDelete(t *testing.T) {
	pool, prefix := testRedis(t)
	repository := &RedisCommandRepository{RedisPool: pool, RedisKeyPrefix: prefix}

	command := &Command{Name: "door", Image: Door, Value: "echo open"}

	err := repository.Create(command, "token")
	if err != nil {
		t.Fatal(err)
	}

	loaded := *command

	err = repository.Delete(command)
	if err != nil {
		t.Fatal(err)
	}

	err = repository.Update(&loaded, &Command{Name: "door", Image: Door, Value: "echo open"})
	if !errors.Is(err, ErrNoRecord) {
		t.Fatalf("expected %v, got %v", ErrNoRecord, err)
	}

	err = repository.UpdateToken(&loaded, "new-token")
	if !errors.Is(err, ErrNoRecord) {
		t.Fatalf("expected %v, got %v", ErrNoRecord, err)
	}

	_, err = repository.Find(command.ID)
	if !errors.Is(err, ErrNoRecord) {
		t.Fatalf("expected %v, got %v", ErrNoRecord, err)
	}
}

func TestRedisCallCreateAfterCommandDelete(t *testing.T) {
	pool, prefix := testRedis(t)
	commands := &RedisCommandRepository{RedisPool: pool, RedisKeyPrefix: prefix}
	calls := &RedisCallRepository{RedisPool: pool, RedisKeyPrefix: prefix, HistoryLimit: 10}

	command := &Command{Name: "door", Image: Door, Value: "echo open"}

	err := commands.Create(command, "token")
	if err != nil {
		t.Fatal(err)
	}

	loaded := *command

	err = commands.Delete(command)
	if err != nil {
		t.Fatal(err)
	}

	call := NewCall(&loaded)

	err = calls.Create(call)
	if !errors.Is(err, ErrNoRecord) {
		t.Fatalf("expected %v, got %v", ErrNoRecord, err)
	}

	conn := pool.Get()
	defer conn.Close()

	count, err := redis.Int(conn.Do(
		"EXISTS",
		prefix+callKeyPrefix+strconv.FormatInt(call.ID, 10),
		prefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(loaded.ID),
	))
	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatalf("expected no call hash or call list, got %d keys", count)
	}
}

func TestRedisAuthenticateReturnsEmailLookupError(t *testing.T) {
	pool, prefix := testRedis(t)
	repository := &RedisUserRepository{RedisPool: pool, RedisKeyPrefix: prefix}

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", prefix+userKeyPrefix+userEmailKey, "corrupted")
	if err != nil {
		t.Fatal(err)
	}

	_, err = repository.Authenticate("nobody", "password1")
	if err == nil || errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected the Redis error, got %v", err)
	}
}
