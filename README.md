# kubectl-multidoc

Pipe your `kubectl get <resource> -oyaml` commands here to split the `items`
array into multidoc YAML. Then try piping it to
[yamlgrep](https://github.com/tnyeanderson/yamlgrep) to filter it!

```bash
go install github.com/tnyeanderson/kubectl-multidoc

# Then
kubectl get po -oyaml | kubectl multidoc
```

This program reads a Kubernetes API response in YAML format (usually from
`kubectl get <resource> -oyaml`), and outputs it as a YAML multidoc with
each member of the `items` array as its own document.

Importantly, this program does not understand YAML. It does not attempt to
load the YAML into memory, or to ensure it is valid. It checks for the
beginning of the "items" array definition, then goes line by line and
changes any array start token ("- ") to a multidoc separator, and unindents
the lines by one level (two spaces). This is done to make it lightning fast!

Therefore, this program is highly dependent on the formatting of its input.
It expects:

  - Two space indentation
  - Non-indented array start tokens (e.g. the hyphen starting an array item
    should be at the same level as its parent--in the case of the "items"
    array, not indented at all)

The output of this program is not guaranteed to be valid YAML.

## Performance

Usually it's easy enough to use `yq`:

```bash
yq '.items[] | split_doc'
```

But the above loads the YAML into memory (I think) and is slower (definitely).
Check out the benchmarks, this program is about 100 times faster than using
`yq`!

```bash
$ go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/tnyeanderson/kubectl-multidoc
cpu: AMD Ryzen 5 2600 Six-Core Processor            
BenchmarkSplitToMultidoc-12    	1000000000	         0.0005545 ns/op
BenchmarkYQ-12                 	1000000000	         0.05409 ns/op
PASS
ok  	github.com/tnyeanderson/kubectl-multidoc	0.430s
```
