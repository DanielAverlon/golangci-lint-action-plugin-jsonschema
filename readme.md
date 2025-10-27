## Build and Install Plugins

```yml
      - name: Build and install plugins
        run: |
          git clone https://github.com/golangci/example-plugin-linter.git
          cd example-plugin-linter
          go build -o '${{ github.workspace }}/.plugins/example.so' -buildmode=plugin plugin/example.go
        working-directory: ${{ runner.temp }}
        env:
          CGO_ENABLED: 1
```

## Install and Run golangci-lint

```yml
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v7
        with:
          version: v2.1.1
          # The installation mode `goinstall` always uses `CGO_ENABLED=1`.
          install-mode: goinstall
```

## Full Example

The [Workflow](https://github.com/golangci/golangci-lint-action-plugin-example/blob/main/.github/workflows/basic.yml).
