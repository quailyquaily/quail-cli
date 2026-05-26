# 作者评论管理

## 背景

作者需要看到各个 list 的最新评论。server 已有评论查询和处理接口，但 CLI 还没有入口。

## 目标

先让作者能在 CLI 里看到最新评论，再处理评论。第一版不做复杂筛选。

## 第一版命令

```bash
quail-cli comments latest --limit 50
quail-cli comments list --list <list_id_or_slug> --limit 20 --offset 0
quail-cli comments approve <comment_id>
quail-cli comments reject <comment_id>
quail-cli comments spam <comment_id>
quail-cli comments delete <comment_id>
```

`comments latest` 第一版在 client 侧实现：先读取当前用户的 lists，再逐个拉取每个 list 的评论，最后按 `created_at` 倒序合并。

## 可复用接口

- `GET /users/me`：得到当前用户 ID。
- `GET /users/{userID}/lists`：得到作者的 lists。
- `GET /dashboard/lists/{listID}/comments?offset=&limit=`：读取某个 list 的评论。现有 dashboard route 只按数字 ID 工作。
- `PUT /comments/{commentID}/approve`：批准评论。
- `PUT /comments/{commentID}/reject`：改回待审核。
- `PUT /comments/{commentID}/spam`：标记垃圾评论。server 已有这条 route，handler 会更新为 `CommentStatusSpam`。
- `DELETE /comments/{commentID}`：删除评论。

## 输出

`human` 模式建议输出这些字段：

- comment ID
- list ID 和 list 标题
- post ID 和 post 标题
- 作者名
- 状态
- 创建时间
- 内容摘要

## Server 缺口

`/dashboard/lists/{listID}/comments` 当前不在 API Key middleware 的匹配范围内。若 CLI 要支持 API Key 调用作者评论功能，需要做一个选择：

新增 `GET /lists/{listIDOrSlug}/comments`，复用当前 handler，并让它走 `ListCooperatorRequired`。这个 route 需要同时支持 list ID 和 slug，和其它 list/post API 保持一致。

这样 CLI 不依赖 dashboard 路由，也只新增一个 server route。

评论处理接口已经存在，但 CLI 使用 OAuth token 或 API Key 时，还需要 server 放行这些路由：

- `PUT /comments/{commentID}/approve`
- `PUT /comments/{commentID}/reject`
- `PUT /comments/{commentID}/spam`
- `DELETE /comments/{commentID}`

OAuth 侧把这些路由加进现有 CLI 使用的 scope，不新增单独的 comment scope。API Key 侧只放行这些 comment 路由，不扩大到 dashboard 或其它管理接口。

## 不做

- 第一版不做跨 list 聚合 server API。
- 第一版不做按状态、关键字、时间段筛选。
- 第一版不做评论通知。

## 验收

- 作者能看到所有自己 list 的最新评论。
- cooperator 能看到自己有权限 list 的评论。
- 非作者不能读取或操作别人的 list 评论。
- API Key 配置下也能使用评论查询和处理命令。
