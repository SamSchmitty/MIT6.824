# Build Code
There is an issue with building and debugging in VSCode. Run the following command to build and allow debugging:
> go build -gcflags="all=-N -l" -buildmode=plugin ../mrapps/wc.go