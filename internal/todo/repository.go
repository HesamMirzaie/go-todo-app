package todo

import (
	"encoding/binary"
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
)

const bucketName = "todos"

type Repository struct {
	db *bolt.DB
}

func OpenRepository(path string) (*Repository, error) {
	db, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	}); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) List() ([]Todo, error) {
	todos := []Todo{}
	err := r.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucketName)).ForEach(func(_, value []byte) error {
			todo := Todo{}
			if err := json.Unmarshal(value, &todo); err != nil {
				return err
			}

			todos = append(todos, todo)
			return nil
		})
	})

	return todos, err
}

func (r *Repository) Get(id int) (Todo, error) {
	todo := Todo{}
	err := r.db.View(func(tx *bolt.Tx) error {
		value := tx.Bucket([]byte(bucketName)).Get(todoKey(id))
		if value == nil {
			return ErrNotFound
		}

		return json.Unmarshal(value, &todo)
	})

	return todo, err
}

func (r *Repository) Create(input Input) (Todo, error) {
	todo := Todo{}
	err := r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}

		todo = Todo{
			ID:        int(id),
			Completed: input.Completed,
			Body:      input.Body,
		}

		return save(bucket, todo)
	})

	return todo, err
}

func (r *Repository) Update(id int, input Input) (Todo, error) {
	todo := Todo{}
	err := r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		value := bucket.Get(todoKey(id))
		if value == nil {
			return ErrNotFound
		}
		if err := json.Unmarshal(value, &todo); err != nil {
			return err
		}

		todo.Body = input.Body
		todo.Completed = input.Completed

		return save(bucket, todo)
	})

	return todo, err
}

func (r *Repository) Delete(id int) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		key := todoKey(id)
		if bucket.Get(key) == nil {
			return ErrNotFound
		}

		return bucket.Delete(key)
	})
}

func save(bucket *bolt.Bucket, todo Todo) error {
	data, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	return bucket.Put(todoKey(todo.ID), data)
}

func todoKey(id int) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(id))
	return key
}
