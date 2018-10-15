// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/mattermost/mattermost-server/model"
import store "github.com/mattermost/mattermost-server/store"

// MentionStore is an autogenerated mock type for the MentionStore type
type MentionStore struct {
	mock.Mock
}

// Save provides a mock function with given fields: post, id
func (_m *MentionStore) Save(post *model.Post, id string) store.StoreChannel {
	ret := _m.Called(post, id)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(*model.Post, string) store.StoreChannel); ok {
		r0 = rf(post, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// View provides a mock function with given fields: mention
func (_m *MentionStore) View(mention *model.Mention) store.StoreChannel {
	ret := _m.Called(mention)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(*model.Mention) store.StoreChannel); ok {
		r0 = rf(mention)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}
