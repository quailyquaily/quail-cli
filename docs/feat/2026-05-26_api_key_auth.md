# API Key 配置

## 背景

`quail-cli` 现在把 OAuth token 当作唯一凭据。用户必须跑一次浏览器登录，CLI 再保存 `access_token`、`refresh_token` 和 `expiry`。

`quaily-server` 已支持 API Key 鉴权：请求可以带 `Authorization: Bearer QK-...`，也可以带 `X-QUAIL-KEY: QK-...`。缺口在 CLI：没有配置 API Key 的入口。

## 目标

让用户可以用 API Key 跑 CLI 和 MCP，不必依赖浏览器 OAuth 流程。

## 第一版命令

```bash
quail-cli login --api-key
quail-cli login --api-key QK-xxx
```

不带值时从隐藏输入读取，避免 key 进入 shell history。带值的形式保留给脚本环境使用。

## 配置

只加一个字段：

```yaml
app:
  api_key: ""
  access_token: ""
  refresh_token: ""
  expiry: ""
  token_type: ""
```

读取顺序：

1. 命令参数传入的 key。
2. 环境变量 `QUAIL_API_KEY`。
3. 配置文件里的 `app.api_key`。
4. OAuth token。

只要读到了 API Key，CLI 就不刷新 OAuth token。

## 实现要点

- `client.Client` 可以继续保存一个 token 字符串，不需要抽象出认证 provider。
- 发送请求时，API Key 继续用 `Authorization: Bearer QK-...`。server 已能识别 `QK-` 前缀。
- `initConfig` 需要避免在已有 API Key 或 `QUAIL_API_KEY` 时自动触发 OAuth login。
- 写入配置文件时不要打印 API Key 明文。
- 新建配置文件时尽量使用只允许当前用户读写的权限。

## 不做

- 第一版不做 `apikey create`、`apikey delete`、`apikey list`。这些属于 key 管理，不是“让 CLI 能用 API Key”的必要条件。
- 第一版不新增 server 权限模型。server 已能验证 API Key，CLI 先用现有能力。

## 验收

- 设置 `QUAIL_API_KEY` 后可以直接执行 `quail-cli me`。
- 执行 `quail-cli login --api-key` 后，后续命令不再打开 OAuth 浏览器流程。
- OAuth 配置仍然可用。
- MCP 模式复用同一套配置。
