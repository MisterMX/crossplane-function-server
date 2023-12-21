# Crossplane Function Server

Crossplane Function Server is a high-level runtime based on the
[Crossplane Function SDK](https://github.com/crossplane/function-sdk-go).
It aims to solve two issues with the current implementation of Go
composition functions:

* Provide a simpler and more abstract API to implement functions since the
  native Go SDK is very bare-metal and requires developers to deal with
  low-level Protobuf structs.
* Serve multiple functions at once using the same image to allow complex
  Crossplane platforms to reduce the required number of function images to a
  bare minimum. With the original SDK one image can only run a single function.

## Features

### High Level Function API

Function Server provides a high-level API that requires users to only deal with
`runtime.Object`s in order to define the desired resources to be created.
This allows developers to focus on the actual implementation of their function
logic without worrying about writing boilerplate or type conversion code.

### Serve Multiple Functions at Once

Function Server is able to serve multiple composition functions at once by
implementing his own high-level virtual function API called `ServerFunction`.
One can register as many server functions during startup as needed.

From a Crossplane perspective a Function Server still acts and looks like a
single function. Which server function should be executed is determined by
the input that is defined in every composition.

## Example

See [`examples`](./examples).
