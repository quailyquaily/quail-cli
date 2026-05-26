package comments

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
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
		Use:   "comments",
		Short: "Manage comments on your lists",
	}

	cmd.AddCommand(newLatestCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newOperateCmd("approve"))
	cmd.AddCommand(newOperateCmd("reject"))
	cmd.AddCommand(newOperateCmd("spam"))
	cmd.AddCommand(newOperateCmd("delete"))

	return cmd
}

func newLatestCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "latest",
		Short: "List latest comments across your lists",
		Run: func(cmd *cobra.Command, args []string) {
			if limit <= 0 {
				limit = 50
			}
			cl := cmd.Context().Value(common.CTX_CLIENT{}).(*client.Client)
			format := cmd.Context().Value(common.CTX_FORMAT{}).(string)

			me, err := cl.GetMe()
			if err != nil {
				slog.Error("failed to get current user", "error", err)
				return
			}
			lists, err := cl.GetUserLists(me.Data.ID)
			if err != nil {
				slog.Error("failed to get lists", "error", err)
				return
			}

			items := make([]client.Comment, 0)
			for _, list := range lists {
				resp, err := cl.GetCommentsByList(strconv.FormatUint(list.ID, 10), 0, limit)
				if err != nil {
					slog.Warn("failed to get list comments", "list_id", list.ID, "error", err)
					continue
				}
				items = append(items, resp.Data.Items...)
			}
			sort.Slice(items, func(i, j int) bool {
				return items[i].CreatedAt.After(items[j].CreatedAt)
			})
			if len(items) > limit {
				items = items[:limit]
			}

			if format == common.FORMAT_JSON {
				client.PrettyPrintJSON(map[string]any{"data": items})
				return
			}
			printComments(items)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 50, "Comment list limit")
	return cmd
}

func newListCmd() *cobra.Command {
	var list string
	var offset int
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List comments for a list",
		Run: func(cmd *cobra.Command, args []string) {
			if list == "" {
				cmd.Help()
				return
			}
			if limit <= 0 {
				limit = 20
			}

			cl := cmd.Context().Value(common.CTX_CLIENT{}).(*client.Client)
			format := cmd.Context().Value(common.CTX_FORMAT{}).(string)
			resp, err := cl.GetCommentsByList(list, offset, limit)
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
	cmd.Flags().StringVar(&list, "list", "", "List id or slug")
	cmd.Flags().IntVar(&offset, "offset", 0, "Comment list offset")
	cmd.Flags().IntVar(&limit, "limit", 20, "Comment list limit")
	return cmd
}

func newOperateCmd(op string) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s <comment_id>", op),
		Short: fmt.Sprintf("%s a comment", title(op)),
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			commentID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil || commentID == 0 {
				slog.Error("invalid comment id", "comment_id", args[0])
				return
			}

			cl := cmd.Context().Value(common.CTX_CLIENT{}).(*client.Client)
			if err := cl.OperateComment(commentID, op); err != nil {
				slog.Error("failed to operate comment", "op", op, "error", err)
				return
			}

			format := cmd.Context().Value(common.CTX_FORMAT{}).(string)
			if format == common.FORMAT_JSON {
				client.PrettyPrintJSON(map[string]any{
					"comment_id": commentID,
					"operation":  op,
					"ok":         true,
				})
				return
			}
			fmt.Printf("comment %d %s ok\n", commentID, op)
		},
	}
}

func title(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func printComments(items []client.Comment) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tLIST\tPOST\tAUTHOR\tSTATUS\tCREATED_AT\tCONTENT")
	for _, item := range items {
		author := strconv.FormatUint(item.AuthorID, 10)
		if item.Author != nil && item.Author.Name != "" {
			author = item.Author.Name
		}
		fmt.Fprintf(w, "%d\t%d\t%d\t%s\t%d\t%s\t%s\n",
			item.ID,
			item.ListID,
			item.PostID,
			author,
			item.Status,
			formatTime(item.CreatedAt),
			strings.ReplaceAll(item.Content, "\n", " "),
		)
	}
	w.Flush()
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
