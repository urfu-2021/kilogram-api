package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"kilogram-api/model"
	"kilogram-api/server"
	"time"
)

func (r *mutationResolver) Register(ctx context.Context, login string, password string, name string) (*model.User, error) {
	r.UsersMu.Lock()
	defer r.UsersMu.Unlock()

	if _, ok := r.UsersByLogin[login]; ok {
		return nil, ErrUserAlreadyExists
	}

	user := &model.User{Login: login, Name: name, Password: password}

	r.Users = append(r.Users, user)
	r.UsersByLogin[login] = user

	return user, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, image *string, name *string) (*model.User, error) {
	user := GetCurrentUserFrom(ctx)

	if user == nil {
		return nil, ErrNotAuthorized
	}

	if image != nil {
		user.M.Lock()
		defer user.M.Unlock()

		if err := validateBase64(*image); err != nil {
			return nil, err
		}

		user.Image = image
	}

	if name != nil {
		user.M.Lock()
		defer user.M.Unlock()

		user.Name = *name
	}

	return user, nil
}

func (r *mutationResolver) UpsertUserMeta(ctx context.Context, key string, val string) (*model.User, error) {
	user := GetCurrentUserFrom(ctx)

	if user == nil {
		return nil, ErrNotAuthorized
	}

	user.M.Lock()
	defer user.M.Unlock()

	user.Meta = appendMeta(user.Meta, key, val)

	return user, nil
}

func (r *mutationResolver) CreateChat(ctx context.Context, typeArg model.ChatType, name string, members []string) (*model.Chat, error) {
	user := GetCurrentUserFrom(ctx)

	if user == nil {
		return nil, ErrNotAuthorized
	}

	uniqMembers := map[string]*model.User{user.Login: user}

	r.UsersMu.RLock()
	for _, login := range members {
		if member, ok := r.UsersByLogin[login]; ok {
			uniqMembers[login] = member
		}
	}
	r.UsersMu.RUnlock()

	if (typeArg == model.ChatTypeChannel || typeArg == model.ChatTypeGroup) && len(uniqMembers) < 3 {
		return nil, ErrGroupChatSize
	}

	if typeArg == model.ChatTypePrivate && len(uniqMembers) != 2 {
		return nil, ErrPrivateChatSize
	}

	chatMembers := make([]*model.User, 0, len(uniqMembers))

	for _, member := range uniqMembers {
		chatMembers = append(chatMembers, member)
	}

	chat := &model.Chat{
		ID:   fmt.Sprint(len(r.Chats)),
		Type: typeArg,

		Name: name,

		AllMembers:        chatMembers,
		AllMembersByLogin: uniqMembers,

		AllMessagesByID: make(map[string]*model.Message),

		Owner:      user,
		OwnerLogin: &user.Login,

		CreatedAt: time.Now(),

		Observers: make(map[string]*model.ChatObserver),
	}

	r.ChatsMu.Lock()
	defer r.ChatsMu.Unlock()

	r.Chats = append(r.Chats, chat)
	r.ChatsByID[chat.ID] = chat

	return chat, nil
}

func (r *mutationResolver) UpdateChat(ctx context.Context, id string, image *string, name *string) (*model.Chat, error) {
	user := GetCurrentUserFrom(ctx)

	if user == nil {
		return nil, ErrNotAuthorized
	}

	r.ChatsMu.RLock()
	chat, ok := r.ChatsByID[id]
	r.ChatsMu.RUnlock()

	if !ok {
		return nil, ErrChatDoesnotExists
	}

	if chat.Owner != user {
		return nil, ErrNotAuthorized
	}

	if image != nil {
		chat.M.Lock()
		defer chat.M.Unlock()

		if err := validateBase64(*image); err != nil {
			return nil, err
		}

		chat.Image = image
	}

	if name != nil {
		chat.M.Lock()
		defer chat.M.Unlock()

		chat.Name = *name
	}

	return chat, nil
}

func (r *mutationResolver) UpsertChatMeta(ctx context.Context, id string, key string, val string) (*model.Chat, error) {
	user := GetCurrentUserFrom(ctx)

	if user == nil {
		return nil, ErrNotAuthorized
	}

	r.ChatsMu.RLock()
	chat, ok := r.ChatsByID[id]
	r.ChatsMu.RUnlock()

	if !ok {
		return nil, ErrChatDoesnotExists
	}

	if chat.Owner != user {
		return nil, ErrNotAuthorized
	}

	chat.M.Lock()
	defer chat.M.Unlock()

	chat.Meta = appendMeta(chat.Meta, key, val)

	return chat, nil
}

func (r *mutationResolver) SendMessage(ctx context.Context, chatID string, text string) (*model.Message, error) {
	user := GetCurrentUserFrom(ctx)

	r.ChatsMu.RLock()
	chat, ok := r.ChatsByID[chatID]
	r.ChatsMu.RUnlock()

	if !ok {
		return nil, ErrChatDoesnotExists
	}

	var authorLogin *string

	if user != nil {
		authorLogin = &user.Login
	}

	if user != nil && chat.ID != model.SpamChatID {
		if _, ok := chat.AllMembersByLogin[user.Login]; !ok {
			return nil, ErrMembership
		}
	}

	if chat.Type == model.ChatTypeChannel && user != chat.Owner {
		return nil, model.ErrUnauthorized
	}

	message := &model.Message{
		ID: fmt.Sprint(len(chat.AllMessages)),

		Author:      user,
		AuthorLogin: authorLogin,

		CreatedAt: time.Now(),

		Text: text,
	}

	chat.M.Lock()
	chat.AllMessages = append(chat.AllMessages, message)
	chat.AllMessagesByID[message.ID] = message
	chat.M.Unlock()

	for _, observer := range chat.Observers {
		if user == nil || message.Author == nil || observer.Login != message.Author.Login {
			observer.Message <- message
		}
	}

	return message, nil
}

// Mutation returns server.MutationResolver implementation.
func (r *Resolver) Mutation() server.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
