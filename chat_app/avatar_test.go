package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

import gomniauthtest "github.com/stretchr/gomniauth/test"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合は、AuthAvatar.GetAvatarURLはErrNoAvatarURLを返すべき")
	}

	testUrl := "http://url-to-avatar"
	testUser = &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return(testUrl, "")
	testChatUser = &chatUser{User: testUser}
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("値が存在する場合はErrorを返すべきではない")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURLは正しいURLを返すべきです")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)

	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではありません")
	}

	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("%sという誤ったURLを返しました", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.jpg")
	err := ioutil.WriteFile(filename, []byte{}, 0777)
	if err != nil {
		return
	}
	defer func() {
		os.Remove(filename)
	}()

	var fileSystemAvatar FileSystemAvatar

	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURLはエラーを返すべきではありません")
	}

	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURLが%sという間違った値を返しました", url)
	}

}
