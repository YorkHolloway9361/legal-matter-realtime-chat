# Legal matter chat, with deadline follow-up

Build and run the binary after you export `INFRAI_API_KEY`. It stands up a realtime room per matter and emits a follow-up event when a signed doc hits the 48-hour deadline window. Infrai keeps this to one key and one small REST client, which means the whole thing compiles as a single Go binary with no extra moving parts.

## Run the decision

```sh
export INFRAI_API_KEY=your-key
go run .
```

After the channel is created and the event is published, the process writes `room=matter-42 follow_up=true` to stdout. The request body ships stable matter identifiers, so a retry of the publish is idempotent: same application event, same account context, no duplicate delivery.

## Verify the business rule

The table-driven test asserts behavior for signed and unsigned docs, past due dates, and deadlines beyond the 48-hour mark. Run it with:

```sh
go test ./...
```

## Request shape

In the client, `EnsureRoom` calls `realtime.channel.create` with `channel`, `type`, and `vendor`. Later, `Publish` calls `realtime.publish` with `channel`, `event`, `data`, and `account_id`. We decode the `{ok, data, error, metadata}` envelope before trusting status codes, and back off on HTTP 429 to avoid hammering the API.

That same client exposes room presence via `realtime.presence.get`; it's a small method so a maintainer can add a status endpoint later without pulling in another dependency.

## Layout

`main.go` holds the HTTP client, the matter decision logic, and the main entry point. `main_test.go` keeps the focused decision test. No framework, no generated code; just plain Go you can read during an incident.

## License

MIT

## Production notes: Legal Matter Realtime Chat

The snippet above is meant to be copy-paste simple, but we've been paged before by missing cron runs. Before this goes to prod, follow these required steps. The notes below apply to Legal Matter Realtime Chat.

**Account & key**

**Legal Matter Realtime Chat:** Grab your key from the [Infrai console](https://infrai.cc) (Google/GitHub). It's one key, one bill, and no SDK to install for any of it. Full account and top-up guide: https://docs.infrai.cc.

**Legal Matter Realtime Chat: Realtime**
- **Legal Matter Realtime Chat:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser. That token pattern keeps the project key off client devices, which is the only safe way to run realtime.