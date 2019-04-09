# Self-Destructing Write.as Posts

This demo application automatically deletes a Write.as post after it's been viewed once.

It accomplishes this by polling a post you own, checking the number of views, and deleting it if `Views > 0`. Since the platform doesn't increment the view count for API calls, this works.

## How to use this

Publish an anonymous post on Write.as, and share the full URL with the person who you'd like to read the message.

Now modify `main.go` to include your login credentials and the post's ID (it's the part between `https://write.as/` and `.md` in the URL you just shared). Then start the application, and leave it running.

## Security

**This is for demonstration only; please don't use it to communicate sensitive information.** This tool is easily circumvented if the person you shared the link with accesses your post with the `.txt` extension, as this does not increment the view count, and thus their view won't be registered. They could also request the post via the API or with a `User-Agent` header that mimics a bot ([any listed here](https://github.com/writeas/web-core/blob/master/bots/bots.go)) to avoid the view count incrementing. These trivial workarounds mean that this tool makes no guarantees that your post is truly only viewed once.

## Polling interval

The Write.as API doesn't rate-limit requests on most endpoints today, but may in the future. Please be mindful of the polling interval you choose if you use this tool over a long period of time, and keep it to at least 30 seconds.
