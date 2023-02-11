# Build Code
There is an issue with building and debugging in VSCode. Run the following command to build and allow debugging:
> go build -gcflags="all=-N -l" -buildmode=plugin ../mrapps/wc.go

Run the build command from the lab when executing from the terminal
> go build -race -buildmode=plugin ../mrapps/wc.go