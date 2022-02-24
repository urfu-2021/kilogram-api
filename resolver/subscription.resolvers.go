package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"kilogram-api/model"
	"kilogram-api/server"
	"math/rand"
)

func (r *subscriptionResolver) NewMessage(ctx context.Context, chatID string) (<-chan *model.Message, error) {
	var login string

	user := GetCurrentUserFrom(ctx)

	if user != nil {
		login = user.Login
	}

	r.ChatsMu.RLock()
	chat, ok := r.ChatsByID[chatID]
	r.ChatsMu.RUnlock()

	if !ok {
		return nil, ErrChatDoesnotExists
	}

	observerID := string(rand.Int31()) // nolint: gosec
	events := make(chan *model.Message, 1)

	go func() {
		<-ctx.Done()

		chat.M.Lock()
		defer chat.M.Unlock()

		delete(chat.Observers, observerID)
	}()

	chat.M.Lock()
	chat.Observers[observerID] = &model.ChatObserver{Login: login, Message: events}
	chat.M.Unlock()

	return events, nil
}

// Subscription returns server.SubscriptionResolver implementation.
func (r *Resolver) Subscription() server.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
