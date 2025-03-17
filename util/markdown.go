package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/quailyquaily/quail-cli/core"
)

func ParseMarkdownWithFrontMatter(filepath string, frontMatterMapping map[string]string) (*core.QuailPostFrontMatter, string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, "", fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	frontMatter := &core.QuailPostFrontMatter{}
	var content strings.Builder
	var isFrontMatter bool
	var frontMatterLines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "---" {
			if isFrontMatter {
				if err := frontMatter.LoadFromYAML(strings.Join(frontMatterLines, "\n"), frontMatterMapping); err != nil {
					return nil, "", fmt.Errorf("could not parse frontmatter: %w", err)
				}

				isFrontMatter = false
			} else {
				isFrontMatter = true
			}
			continue
		}

		if isFrontMatter {
			frontMatterLines = append(frontMatterLines, line)
		} else {
			content.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, "", fmt.Errorf("error reading file: %w", err)
	}

	return frontMatter, content.String(), nil
}
