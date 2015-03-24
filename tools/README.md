# tools package

Only used for code generation.

## Configuration

If you would like to configure the generation process because of custom EAV and attribute models
please create a new file in this folder called `config_*.go` where `*` can be any name.

You must then use the `init()` function to append new values or change existing values
of the configuration variables.

All defined variables in the file `config.go` can be changed. File is documented.

### Why the `init()` function?

From effective Go:

> Finally, each source file can define its own niladic init function to set up
> whatever state is required. (Actually each file can have multiple init functions.)
> And finally means finally: init is called after all the variable declarations in the
> package have evaluated their initializers, and those are evaluated only after all
> the imported packages have been initialized.

## Testing

All test files depends on generated code. All other files of course not.
