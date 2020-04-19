# overrustle-scraper

Due to a [request](https://overrustlelogs.net/assets/twitch_email.png) from Twitch Legal overrustlelogs.net is forced to shut down and delete all saved twitch chat data.

So I decided to create this tool so people can keep their personal twitch chat archive.

This little tool allows you to download all existing chat data from all channels caputred by overrustlelogs.net.

I'm not affiliated with overrustlelogs.net or twitch.

## Usage

Install the tool with go get:

```bash
go get github.com/RustedDart/overrustle-scraper
```

This will build the binary in $GOPATH/bin so make sure that $GOPATH is in your PATH.

To run the script simply run the binary.

The script will create a folder called scrape in your current working directory. All logs will be saved in this folder with the same structure as overrustlelogs.net.

The folder will be around 300GB.


