package laisky_blog_graphql

import (
	"context"
	"fmt"
	"time"

	utils "github.com/Laisky/go-utils"
	"github.com/Laisky/laisky-blog-graphql/blog"
	"github.com/Laisky/laisky-blog-graphql/libs"
	"github.com/Laisky/zap"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func (r *Resolver) BlogPost() BlogPostResolver {
	return &blogPostResolver{r}
}
func (r *Resolver) BlogUser() BlogUserResolver {
	return &blogUserResolver{r}
}

type blogPostResolver struct{ *Resolver }
type blogUserResolver struct{ *Resolver }

// =================
// query resolver
// =================

func (q *queryResolver) BlogPostInfo(ctx context.Context) (*blog.PostInfo, error) {
	return blogDB.LoadPostInfo()
}

func (q *queryResolver) BlogPosts(ctx context.Context, page *Pagination, tag string, categoryURL *string, length int, name string, regexp string) ([]*blog.Post, error) {
	cfg := &blog.BlogPostCfg{
		Page:        page.Page,
		Size:        page.Size,
		Length:      length,
		Tag:         tag,
		Regexp:      regexp,
		CategoryURL: categoryURL,
		Name:        name,
	}
	if results, err := blogDB.LoadPosts(cfg); err != nil {
		return nil, err
	} else {
		return results, nil
	}
}
func (q *queryResolver) BlogPostCategories(ctx context.Context) ([]*blog.Category, error) {
	return blogDB.LoadAllCategories()
}

// ----------------
// blog resolver
// ----------------

func (r *blogPostResolver) MongoID(ctx context.Context, obj *blog.Post) (string, error) {
	return obj.ID.Hex(), nil
}
func (r *blogPostResolver) CreatedAt(ctx context.Context, obj *blog.Post) (*libs.Datetime, error) {
	return libs.NewDatetimeFromTime(obj.CreatedAt), nil
}
func (r *blogPostResolver) ModifiedAt(ctx context.Context, obj *blog.Post) (*libs.Datetime, error) {
	return libs.NewDatetimeFromTime(obj.ModifiedAt), nil
}
func (r *blogPostResolver) Author(ctx context.Context, obj *blog.Post) (*blog.User, error) {
	return blogDB.LoadUserById(obj.Author)
}
func (r *blogPostResolver) Category(ctx context.Context, obj *blog.Post) (*blog.Category, error) {
	return blogDB.LoadCategoryById(obj.Category)
}
func (r *blogPostResolver) Type(ctx context.Context, obj *blog.Post) (BlogPostType, error) {
	switch obj.Type {
	case "markdown":
		return BlogPostTypeMarkdown, nil
	case "slide":
		return BlogPostTypeSlide, nil
	case "html":
		return BlogPostTypeHTML, nil
	}

	return "", fmt.Errorf("unknown blog post type: `%+v`", obj.Type)
}

func (r *blogUserResolver) ID(ctx context.Context, obj *blog.User) (string, error) {
	return obj.ID.Hex(), nil
}

// ============================
// mutations
// ============================

func (r *mutationResolver) BlogCreatePost(ctx context.Context, input NewBlogPost) (*blog.Post, error) {
	user, err := validateAndGetUser(ctx)
	if err != nil {
		libs.Logger.Debug("user invalidate", zap.Error(err))
		return nil, err
	}

	return blogDB.NewPost(user.ID, string(input.Title), input.Name, string(input.Markdown), input.Type.String())
}
func (r *mutationResolver) BlogLogin(ctx context.Context, account string, password string) (user *blog.User, err error) {
	if user, err = blogDB.ValidateLogin(account, password); err != nil {
		libs.Logger.Debug("user invalidate", zap.Error(err))
		return nil, err
	}

	uc := &blog.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.ID.Hex(),
			IssuedAt:  utils.Clock2.GetUTCNow().Unix(),
			ExpiresAt: utils.Clock.GetUTCNow().Add(7 * 24 * time.Hour).Unix(),
		},
		Username:    user.Account,
		DisplayName: user.Username,
	}

	if err = auth.SetLoginCookie(ctx, uc); err != nil {
		libs.Logger.Error("try to set cookie got error", zap.Error(err))
		return nil, errors.Wrap(err, "try to set cookies got error")
	}

	return user, nil
}
func (r *mutationResolver) BlogAmendPost(ctx context.Context, post NewBlogPost) (*blog.Post, error) {
	user, err := validateAndGetUser(ctx)
	if err != nil {
		libs.Logger.Debug("user invalidate", zap.Error(err))
		return nil, err
	}

	return blogDB.UpdatePost(user, post.Name, string(post.Title), string(post.Markdown), string(post.Type))
}
