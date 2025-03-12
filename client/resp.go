package client

import "time"

type (
	UserResponse struct {
		Data struct {
			ID             uint64         `json:"id"`
			Name           string         `json:"name"`
			Email          string         `json:"email"`
			AvatarImageURL string         `json:"avatar_image_url"`
			Bio            string         `json:"bio"`
			Tagline        string         `json:"tagline"`
			CreatedAt      string         `json:"created_at"`
			SocialIDs      []UserSocialID `json:"social_ids"`
			Status         int            `json:"status"`
			UserOptions    struct {
				EditorLayout         string `json:"editor_layout"`
				KindLineBreakEnabled bool   `json:"kind_line_break_enabled"`
				Languages            string `json:"languages"`
			}
		} `json:"data"`
	}
	UserSocialID struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
)

type (
	Post struct {
		ID               uint64    `json:"id"`
		Slug             string    `json:"slug"`
		CoverImageURL    string    `json:"cover_image_url"`
		Title            string    `json:"title"`
		Summary          string    `json:"summary"`
		Content          string    `json:"content"`
		PaidContent      string    `json:"paid_content"`
		UserID           uint64    `json:"user_id"`
		ListID           uint64    `json:"list_id"`
		Tags             string    `json:"tags"`
		Theme            string    `json:"theme"`
		PublishedAt      time.Time `json:"published_at"`
		FirstPublishedAt time.Time `json:"first_published_at"`
		List             List      `json:"list"`
	}

	List struct {
		ID          uint64 `json:"id"`
		Slug        string `json:"slug"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Tagline     string `json:"tagline"`
	}

	PostResponse struct {
		Data Post `json:"data"`
	}

	ListsResponse struct {
		Data []List `json:"data"`
	}

	ListResponse struct {
		Data List `json:"data"`
	}

	PaginationResponse struct {
		Pagination struct {
			Current    uint64 `json:"current"`
			Offset     uint64 `json:"offset"`
			Limit      uint64 `json:"limit"`
			NextOffset uint64 `json:"next_offset"`
			Total      uint64 `json:"total"`
		} `json:"pagination"`
		Langs []string `json:"langs"`
		Items []Post   `json:"items"`
	}

	SearchResponse struct {
		Data PaginationResponse `json:"data"`
	}
)
