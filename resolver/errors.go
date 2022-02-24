package resolver

import "errors"

var (
	ErrChatDoesnotExists = errors.New("chat doesn't exists")
	ErrGroupChatSize     = errors.New("channel or group should have more members")
	ErrInvalidBase64     = errors.New("your base64 is invalid")
	ErrMembership        = errors.New("you are not a member of this chat")
	ErrNotAuthorized     = errors.New("not authorized")
	ErrPrivateChatSize   = errors.New("private chat should have 2 members")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserDoesnotExists = errors.New("user doesn't exists")
)
