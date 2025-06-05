## bxcleaner

`bxcleaner` is a tool to delete local Git branches except the chosen branch. I created this because in every end of the sprint, i need to delete all branch except main to keep my workspace clean.

## how to install
You can run

```
go install github.com/pujidjayanto/bxcleaner@v1.0.0
```

Add $HOME/go/bin to your PATH if it isn’t already:
```
export PATH=$PATH:$(go env GOPATH)/bin
```

and run the tool
```
bxcleaner --help
```