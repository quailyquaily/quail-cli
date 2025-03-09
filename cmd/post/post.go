package post

import (
	"fmt"
	"time"

	"github.com/quail-ink/quail-cli/client"
	"github.com/quail-ink/quail-cli/cmd/common"
	"github.com/quail-ink/quail-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	listSlug  string
	postSlug  string
	doPublish bool
)

func upsertPost(cl *client.Client, filepath string, frontMatterMapping map[string]string, format string) error {
	if filepath == "" {
		return fmt.Errorf("filepath is required")
	}

	frontMatter, content, err := util.ParseMarkdownWithFrontMatter(filepath, frontMatterMapping)
	if err != nil {
		return err
	}

	var datetime *time.Time
	if doPublish {
		datetime = frontMatter.Datetime
		if datetime == nil {
			now := time.Now()
			datetime = &now
		}
	}

	payload := map[string]any{
		"slug":               frontMatter.Slug,
		"cover_image_url":    frontMatter.CoverImageUrl,
		"title":              frontMatter.Title,
		"summary":            frontMatter.Summary,
		"content":            content,
		"datetime":           datetime,
		"first_published_at": frontMatter.Datetime,
		"tags":               frontMatter.Tags,
		"theme":              frontMatter.Theme,
	}

	result, err := cl.CreatePost(listSlug, payload)
	if err != nil {
		return err
	}

	if format == common.FORMAT_JSON {
		client.PrettyPrintJSON(result)
	} else {
		client.PrettyPrintPost(result)
	}

	return nil
}

func modPost(cmd *cobra.Command, cl *client.Client, op, format string) {
	if postSlug == "" || listSlug == "" {
		cmd.Help()
		return
	}
	result, err := cl.ModPost(listSlug, postSlug, op)
	if err != nil {
		fmt.Println(err)
		return
	}
	if format == common.FORMAT_JSON {
		client.PrettyPrintJSON(result)
	} else {
		client.PrettyPrintPost(result)
	}
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post upsert [filepath]\n\tpost <delete||publish|unpublish|deliver>",
		Short: "Manpulate posts",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			}

			format := cmd.Context().Value(common.CTX_FORMAT{}).(string)
			cl := cmd.Context().Value(common.CTX_CLIENT{}).(*client.Client)

			frontMatterMapping := viper.GetStringMapString("post.frontmatter_mapping")

			action := args[0]
			switch action {
			case "upsert":
				if len(args) < 2 {
					cmd.Help()
					return
				}

				filepath := args[1]
				if err := upsertPost(cl, filepath, frontMatterMapping, format); err != nil {
					fmt.Println(err)
					return
				}
			case "delete":
				{
					if postSlug == "" || listSlug == "" {
						cmd.Help()
						return
					}
					result, err := cl.DeletePost(listSlug, postSlug)
					if err != nil {
						fmt.Println(err)
						return
					}
					if format == common.FORMAT_JSON {
						client.PrettyPrintJSON(result)
					} else {
						client.PrettyPrintPost(result)
					}
				}
			case "publish":
				{
					modPost(cmd, cl, "publish", format)
				}
			case "unpublish":
				{
					modPost(cmd, cl, "unpublish", format)
				}
			case "deliver":
				{
					modPost(cmd, cl, "deliver", format)
				}
			default:
				cmd.Help()
			}
		},
	}

	cmd.Flags().StringVarP(&listSlug, "list", "l", "", "Channel slug")
	cmd.Flags().StringVarP(&postSlug, "post", "p", "", "Post slug")
	cmd.Flags().BoolVar(&doPublish, "publish", false, "Publish the post")

	return cmd
}
