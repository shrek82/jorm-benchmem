# create测试

```bash
cd /Users/up/projects/jorm-gin && go test -bench=. -benchmem ./create_bench 
```

## 查找性能测试

```bash
cd /Users/up/projects/jorm-gin && go test -bench=Find -benchmem ./find_bench
```

## update性能测试

```bash
cd /Users/up/projects/jorm-gin && go test -bench=. -benchmem ./update_bench
```