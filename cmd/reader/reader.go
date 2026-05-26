package reader

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/quailyquaily/quail-cli/client"
	"github.com/quailyquaily/quail-cli/cmd/common"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reader",
		Short: "Read subscribed posts and comments",
	}

	cmd.AddCommand(newSubscriptionsCmd())
	cmd.AddCommand(newPostsCmd())
	cmd.AddCommand(newReadCmd())
	cmd.AddCommand(newCommentsCmd())
	cmd.AddCommand(newCommentCmd())

	return cmd
}

func newSubscriptionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "subscriptions",
		Short: "List your subscriptions",
		Run: func(cmd *cobra.Command, args []string) {
			cl := cmd.Context().Value(common.CTX_CLIENT{}).(*client.Client)
			format := cmd.Context().Value(common.CTX_FORMAT{}).(string)

			resp, err := cl.GetSubscriptions()
			if err != nil {
				slog.Error("failed to get subscriptions", "error", err)
				return
			}
			if format == common.FORMAT_JSON {
				client.PrettyPrintJSON(resp)
				return
			}
			printSubscriptions(resp.Data)
		},
	}
}

func newPostsCmd() *cobra.Command {
	var offset int
	var limit int

	cmd := &cobra.Command{
		Use:   "posts",
		Short: "List posts from your subscriptions",
		Run: func(cmd *cobra.Command, args []string) {
			cl := cmd.Context().Value(common.CTX_CLIENT{}).(*client.Client)
			format := cmd.Context().Value(common.CTX_FORMAT{}).(string)

			resp, err := cl.GetSubscribedPosts(offset, limit)
			if err != nil {
				slog.Error("failed to get subscribed posts", "error", err)
				return
			}
			if format == common.FORMAT_JSON {
				client.PrettyPrintJSON(resp)
				return
			}
			printPosts(resp.Data.Items)
		},
	}
	cmd.Flags().IntVar(&offset, "offset", 0, "Post list offset")
	cmd.Flags().IntVar(&limit, "limit", 20, "Post list limit")
	return cmd
}

func newReadCmd() *cobra.Command {
	var list string
	var post string

	cmd := &cobra.Command{
		Use:   "read [url]",
		Short: "Read a post",
		Run: func(cmd *cobra.Command, args []string) {
			cl := cmd.Context().Value(common.CTX_CLIENT{}).(*client.Client)
			format := cmd.Context().Value(common.CTX_FORMAT{}).(string)

			listIDOrSlug := list
			postIDOrSlug := post
			if len(args) > 0 {
				var err error
				listIDOrSlug, postIDOrSlug, err = parsePostURL(args[0])
				if err != nil {
					slog.Error("failed to parse post url", "error", err)
					return
				}
			}
			if listIDOrSlug == "" || postIDOrSlug == "" {
				cmd.Help()
				return
			}

			postResp, err := cl.GetPost(listIDOrSlug, postIDOrSlug)
			if err != nil {
				slog.Error("failed to get post", "error", err)
				return
			}

			contentResp, contentErr := cl.GetPostContent(listIDOrSlug, postIDOrSlug)
			if format == common.FORMAT_JSON {
				ret := map[string]any{
					"post": postResp.Data,
				}
				if contentErr != nil {
					ret["content_error"] = readableContentError(contentErr)
				} else {
					ret["content"] = contentResp.Data
				}
				client.PrettyPrintJSON(ret)
				return
			}

			printPostWithContent(postResp.Data, contentResp, contentErr)
		},
	}
	cmd.Flags().StringVar(&list, "list", "", "List id or slug")
	cmd.Flags().StringVar(&post, "post", "", "Post id or slug")
	return cmd
}

func newCommentsCmd() *cobra.Command {
	var postID uint64
	var offset int
	var limit int

	cmd := &cobra.Command{
		Use:   "comments",
		Short: "List comments for a post",
		Run: func(cmd *cobra.Command, args []string) {
			if postID == 0 {
				cmd.Help()
				return
			}

			cl := cmd.Context().Value(common.CTX_CLIENT{}).(*client.Client)
			format := cmd.Context().Value(common.CTX_FORMAT{}).(string)
			resp, err := cl.GetCommentsByPost(postID, offset, limit)
			if err != nil {
				slog.Error("failed to get comments", "error", err)
				return
			}
			if format == common.FORMAT_JSON {
				client.PrettyPrintJSON(resp)
				return
			}
			printComments(resp.Data.Items)
		},
	}
	cmd.Flags().Uint64Var(&postID, "post", 0, "Post id")
	cmd.Flags().IntVar(&offset, "offset", 0, "Comment list offset")
	cmd.Flags().IntVar(&limit, "limit", 20, "Comment list limit")
	return cmd
}

func newCommentCmd() *cobra.Command {
	var postID uint64
	var content string

	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Create a comment for a post",
		Run: func(cmd *cobra.Command, args []string) {
			content = strings.TrimSpace(content)
			if postID == 0 || content == "" {
				cmd.Help()
				return
			}

			cl := cmd.Context().Value(common.CTX_CLIENT{}).(*client.Client)
			format := cmd.Context().Value(common.CTX_FORMAT{}).(string)
			resp, err := cl.CreateComment(postID, content)
			if err != nil {
				slog.Error("failed to create comment", "error", err)
				return
			}
			if format == common.FORMAT_JSON {
				client.PrettyPrintJSON(resp)
				return
			}
			printComments([]client.Comment{resp.Data})
		},
	}
	cmd.Flags().Uint64Var(&postID, "post", 0, "Post id")
	cmd.Flags().StringVar(&content, "content", "", "Comment content")
	return cmd
}

func parsePostURL(raw string) (string, string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", "", err
	}
	if u.Scheme != "https" || u.Host != "quaily.com" {
		return "", "", fmt.Errorf("only https://quaily.com/{list}/{post} URLs are supported")
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid quaily post url")
	}
	return parts[0], parts[1], nil
}

func printSubscriptions(items []client.Subscription) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tLIST\tSLUG\tTYPE\tEMAIL\tPAID_EXPIRY")
	for _, item := range items {
		listTitle := ""
		listSlug := strconv.FormatUint(item.ListID, 10)
		if item.List != nil {
			listTitle = item.List.Title
			listSlug = item.List.Slug
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%t\t%s\n",
			item.ID,
			listTitle,
			listSlug,
			item.Type,
			item.EmailEnabled,
			formatTimePtr(item.PaidExpiry),
		)
	}
	w.Flush()
}

func printPosts(items []client.Post) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tLIST\tSLUG\tTITLE\tPUBLISHED_AT\tPAID")
	for _, item := range items {
		listSlug := item.List.Slug
		if listSlug == "" {
			listSlug = strconv.FormatUint(item.ListID, 10)
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%t\n",
			item.ID,
			listSlug,
			item.Slug,
			item.Title,
			formatTime(item.PublishedAt),
			item.IsPaidContent,
		)
	}
	w.Flush()
}

func printPostWithContent(post client.Post, content *client.PostContentResponse, contentErr error) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "ID:\t%d\n", post.ID)
	fmt.Fprintf(w, "List:\t%d\n", post.ListID)
	fmt.Fprintf(w, "Slug:\t%s\n", post.Slug)
	fmt.Fprintf(w, "Title:\t%s\n", post.Title)
	fmt.Fprintf(w, "Summary:\t%s\n", post.Summary)
	fmt.Fprintf(w, "Published At:\t%s\n", formatTime(post.PublishedAt))
	if contentErr != nil {
		fmt.Fprintf(w, "Content:\t%s\n", readableContentError(contentErr))
		w.Flush()
		return
	}
	fmt.Fprintf(w, "Content:\n%s\n", content.Data.FreeContent)
	if content.Data.PaidContent != "" {
		fmt.Fprintf(w, "\nPaid Content:\n%s\n", content.Data.PaidContent)
	}
	w.Flush()
}

func printComments(items []client.Comment) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tPOST\tAUTHOR\tSTATUS\tCREATED_AT\tCONTENT")
	for _, item := range items {
		author := strconv.FormatUint(item.AuthorID, 10)
		if item.Author != nil && item.Author.Name != "" {
			author = item.Author.Name
		}
		fmt.Fprintf(w, "%d\t%d\t%s\t%d\t%s\t%s\n",
			item.ID,
			item.PostID,
			author,
			item.Status,
			formatTime(item.CreatedAt),
			strings.ReplaceAll(item.Content, "\n", " "),
		)
	}
	w.Flush()
}

func readableContentError(err error) string {
	if err == nil {
		return ""
	}
	if strings.Contains(err.Error(), "status code: 401") {
		return "no access to this post content"
	}
	return err.Error()
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func formatTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
