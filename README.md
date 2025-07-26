# bradio

cli tool for searching radio stations on radio-browser by name or tag. Sorting by popularity.

```
bradio --name 'Milano Lounge'

bradio --tag 'ambient'

bradio --tag 'chillout' --limit 30

```
## Building

```
go build -o bin/bradio .

```


## If you use mpv for listening to radio

```

bradio --tag 'lounge' | fzf --nth 1 --preview="echo {}" --preview-window=bottom:2:nohidden | awk -F $' - ' '{print $4}' | mpv --playlist=-

```

or try function

```
function br() { 
bradio "$@" | fzf --nth 1 --preview='echo {}' --preview-window=bottom:2:nohidden | awk -F $'; ' '{print $4}' | mpv --playlist=- 
}

br --tag 'lounge' --limit 24

```# Bradio
