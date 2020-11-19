package mocks

import "github.com/projecteru2/pistage/store"

// Mocks .
func Mocks() (*Store, func()) {
	ori := store.GetStore()

	mocked := &Store{}
	store.SetStore(mocked)

	return mocked, func() { store.SetStore(ori) }
}
