# telegram-msgsender

Minimal Go package for sending Telegram messages from scrapers and monitors. Designed to feel similar to Discord webhooks — pass a token + chat target, build a message, send it.

Messages are delivered asynchronously via a per-bot queue. Rate limits (HTTP 429) are handled automatically using Telegram's `retry_after` value.

## Install

```
go get github.com/pasa33/telegram-msgsender
```

## Usage

```go
import tg "github.com/pasa33/telegram-msgsender"

tg.Send(botToken, chatID, threadID, tg.NewMessage().
    Line(tg.Bold("New item found")).
    Line(tg.Field("Name", itemName)).
    Line(tg.Field("Price", "$99")),
)
```

Set `threadID` to `0` if you don't use topics/threads.

### Attach an image

Prefer `ImageBytes` over `Image` when the URL might be banned by Telegram (common with scraped images).

```go
// from bytes (recommended for scraped images)
tg.Send(botToken, chatID, threadID, tg.NewMessage().
    Line(tg.Bold("Alert")).
    ImageBytes(imgBytes),
)

// from URL
tg.Send(botToken, chatID, threadID, tg.NewMessage().
    Line(tg.Bold("Alert")).
    Image("https://example.com/image.jpg"),
)
```

### Formatting helpers (MarkdownV2)

All helpers handle escaping internally — you don't need to call `Escape()` manually when using them.

```go
tg.Bold("text")            // *text*
tg.Italic("text")          // _text_
tg.Code("text")            // `text`
tg.Pre("block")            // ```\nblock\n```
tg.Link("label", url)      // [label](url)
tg.Field("Key", "value")   // *Key*: value
tg.Escape("raw text")      // escape plain text for use in raw MarkdownV2
```

Mix them inline:

```go
tg.NewMessage().
    Line(tg.Bold("Signal") + " — " + tg.Link("open chart", chartURL)).
    Line(tg.Field("Entry", "42,100")).
    Line(tg.Field("Stop", "41,800"))
```

### HTML mode

Switch to HTML if you prefer it over MarkdownV2:

```go
tg.NewMessage().
    HTML().
    Line("<b>Bold</b> and <i>italic</i>").
    Line("<code>some code</code>")
```

`Bold()`, `Italic()`, `Code()` helpers still produce MarkdownV2 syntax — in HTML mode, write tags directly.

### Message options

```go
tg.NewMessage().
    Line("sensitive content").
    Protect().    // disables forwarding and saving
    NoPreview()   // disables link preview
```

### Debug mode

Redirect all sends to a single chat — useful during development so you don't spam production channels.

```go
// override everything
tg.SetDebug(debugBotToken, debugChatID, debugThreadID)

// clear
tg.SetDebug("", "", 0)
```

## Notes

- The queue holds up to 100 pending messages per bot token. `Send` returns an error if the queue is full.
- `sendPhoto` does not support `disable_web_page_preview` — that field is silently ignored when an image is attached (Telegram limitation).
