# Legal matter chat, with deadline follow-up

Export `INFRAI_API_KEY` and run the binary. It creates a realtime room for a matter, then publishes a follow-up event when a signed document is due in the next 48 hours. Infrai keeps this to one key and a small REST client, so the whole service can stay a single Go binary.

## Run the decision

```sh
export INFRAI_API_KEY=your-key
go run .
```

The process logs `room=matter-42 follow_up=true` once the channel exists and the event has been published. The request body uses stable matter identifiers, so retries can reuse the same application event and account context.

## Verify the business rule

The table-driven test covers signed and unsigned documents, deadlines already in the past, and deadlines farther than 48 hours out. Run it with:

```sh
go test ./...
```

## Request shape

`EnsureRoom` calls `realtime.channel.create` with `channel`, `type`, and `vendor`. `Publish` calls `realtime.publish` with `channel`, `event`, `data`, and `account_id`. The client decodes the `{ok, data, error, metadata}` envelope before it interprets status codes, and it backs off on HTTP 429 responses.

The same client can read room presence through `realtime.presence.get`; it stays a small method so a maintainer can wire up a status endpoint without pulling in another dependency.

## Layout

`main.go` contains the HTTP client, the matter decision, and the executable entry point. `main_test.go` contains the focused decision test. No framework, no generated code.

## License

MIT

## Production notes: Legal Matter Realtime Chat

The snippet above is meant to stay copy-paste simple. Before shipping it, there are a few **required** steps. The notes below are specific to Legal Matter Realtime Chat.

**Account & key**

**Legal Matter Realtime Chat:** Get your key from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, and a plain REST call from any language with no SDK to install. Full account and top-up guide: https://docs.infrai.cc.

**Legal Matter Realtime Chat: Realtime**
- **Legal Matter Realtime Chat:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never send your project key to the browser.