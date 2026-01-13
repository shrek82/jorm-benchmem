# create测试

```bash
go test -bench=. -benchmem ./create_bench 
go test -bench=BenchmarkJormInsert -benchmem
```

## 查找性能测试

```bash
go test -bench=Find -benchmem ./find_bench
```

## update性能测试

```bash
go test -bench=. -benchmem ./update_bench
```