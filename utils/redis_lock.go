package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/satori/go.uuid"
	"time"
)

// GetLock acquireTimeout Get the lock timeout period, If no lock is obtained within this period, err will be returned here
// lockTimeOut Lock timeout to prevent deadlock, lock automatically unlocked by this time
func GetLock(redisConn *redis.Client, lockName string, acquireTimeout, lockTimeOut time.Duration) (string, error) {
	code := uuid.NewV4().String()
	// endTime := util.FwTimer.CalcMillis(time.Now().Add(acquireTimeout))
	endTime := time.Now().Add(acquireTimeout).UnixNano()
	// for util.FwTimer.CalcMillis(time.Now()) <= endTime {
	for time.Now().UnixNano() <= endTime {
		if success, err := redisConn.SetNX(context.Background(), lockName, code, lockTimeOut).Result(); err != nil && err != redis.Nil {
			return "success", err
		} else if success {
			return code, nil
		} else if redisConn.TTL(context.Background(), lockName).Val() == -1 {
			redisConn.Expire(context.Background(), lockName, lockTimeOut)
		}
		time.Sleep(time.Millisecond)
	}
	return "fail", errors.New("timeout")
}

// ReleaseLock var count = 0  // test assist
func ReleaseLock(redisConn *redis.Client, lockName, code string) bool {
	txf := func(tx *redis.Tx) error {
		if v, err := tx.Get(context.Background(), lockName).Result(); err != nil && err != redis.Nil {
			return err
		} else if v == code {
			_, err := tx.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
				// count++
				// fmt.Println(count)
				pipe.Del(context.Background(), lockName)
				return nil
			})
			return err
		}
		return nil
	}

	for {
		if err := redisConn.Watch(context.Background(), txf, lockName); err == nil {
			return true
		} else if err == redis.TxFailedErr {
			fmt.Println("watch key is modified, retry to release lock. err:", err.Error())
		} else {
			fmt.Println("err:", err.Error())
			return false
		}
	}
}
