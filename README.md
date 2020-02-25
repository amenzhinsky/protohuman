# protohuman

Human-readable protobuf marshaller for Golang.

## Usage

The marshals (for now) protobuf structures in a human-readable format to the given writer:

```go
if err := protohuman.Marshall(os.Stderr, msg); err != nil {
	return err
}
```

Would produce something similar to: 

```
f:666, string_slice:["a", "b"], kv:map[five:5 four:4], buf:[1, 2, 3], state:UNKNOWN, oneofer:(one:"one"), inner:{enabled:true}}
```
