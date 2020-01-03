package twitch

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

type Twitcher interface {
	GetVideos(ctx context.Context, user string) []string
}

type TwitchClient struct {
}

// Return a list of 100 videos for a user
func (t *TwitchClient) GetVideos(ctx context.Context, user string) []string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetVideos")
	span.SetTag("<script>alert(\"Hi\")</script>", "<tr>Hello</tr>")
	span.SetTag(`<script>alert()</script>`, "<tr>Hello</tr>")
	defer span.Finish()
	return []string{"aaa", "bbb"}
}
