# 读者阅读

## 背景

读者最常见的需求是看自己订阅了哪些 list、哪些文章有更新、以及能否直接读文章内容。`quaily-server` 已有大部分接口，`quail-cli` 还没有命令行入口。

## 目标

用最少命令把“订阅列表、订阅文章、正文阅读、评论”跑通。

## 第一版命令

```bash
quail-cli reader subscriptions
quail-cli reader posts --limit 20 --offset 0
quail-cli reader read <url>
quail-cli reader read --list <list_id_or_slug> --post <post_id_or_slug>
quail-cli reader comments --post <post_id> --limit 20 --offset 0
quail-cli reader comment --post <post_id> --content "..."
```

`reader read <url>` 第一版只支持 `https://quaily.com/{list_slug}/{post_slug}` 这种标准文章 URL。CLI 解析出 list slug 和 post slug 后复用同一套读取逻辑。不支持自定义域名。

`reader read` 默认输出标题、摘要、正文和付费正文访问结果。没有付费权限时，应明确显示无权限，而不是只打印原始 401。

`reader comment` 第一版只支持 `--content` 参数，不做编辑器、stdin 多行输入和引用评论。

## 可复用接口

- `GET /subscriptions/`：读取当前用户订阅。
- `GET /posts/subscribed?offset=&limit=`：读取订阅 list 的最新已发布文章。
- `GET /lists/{listIDOrSlug}/posts/{postIDOrSlug}`：读取文章元数据。
- `GET /lists/{listIDOrSlug}/posts/{postIDOrSlug}/content`：读取正文和付费正文，server 会检查权限。
- `GET /comments?post_id=&offset=&limit=`：读取文章评论。
- `POST /comments`：发表评论。

## 输出

`human` 模式建议输出：

- 订阅：list 名称、slug、订阅类型、是否开启邮件、付费到期时间。
- 文章列表：发布时间、list、标题、slug、是否付费内容。
- 正文：标题、URL、免费正文、可访问的付费正文。
- 评论：评论 ID、作者、状态、时间、内容摘要。

`json` 模式保持原始字段，方便脚本处理。

## Server 缺口

`GET /comments` 是公开读取接口，可以直接复用。

`POST /comments` 已有 handler，但 CLI 的 OAuth token 和 API Key 都需要 server 放行这条路由：

- OAuth：把 `POST /comments` 加进现有 CLI 使用的 scope，不新增单独的 `comment.write` scope。
- API Key：把 `POST /comments` 加进 API Key 鉴权匹配范围。

## 不做

- 第一版不做反应、已读状态。
- 第一版不做按 URL 直接发评论；用户先用 `reader read <url>` 取得 post ID，再用 `reader comment --post`。
- 第一版不支持自定义域名文章 URL。

## 验收

- API Key 和 OAuth token 都能读取订阅文章列表。
- 免费文章可直接读正文。
- 付费文章在有权限时返回正文，无权限时显示清楚的错误信息。
- 可以读取文章评论。
- 登录用户可以发表评论。
- `--format json` 的输出可以被 `jq` 直接处理。
