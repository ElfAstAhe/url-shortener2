# Iteration 17

## Результаты оптимизации по памяти на основе данных бенчмарков

### Оптимизация путём создания исходящих срезов с заданным capacity на основе кол-ва входящих данных

```
elf@elf-desktop:~/work/study/go/go-course/url-shortener2/internal/ep/mapper$ go tool pprof -top -diff_base=base_mem.pprof result_mem.pprof
```

```
File: mapper.test
Build ID: 6684f98fa03a40e2053629bdbb1208d597a17012
Type: alloc_space
Time: 2026-01-23 18:12:50 MSK
Showing nodes accounting for -2.46GB, 36.57% of 6.73GB total
Dropped 45 nodes (cum <= 0.03GB)
flat  flat%   sum%        cum   cum%
-1.25GB 18.53% 18.53%    -0.93GB 13.84%  github.com/ElfAstAhe/url-shortener2/internal/ep/mapper.UserShortensFromModel (inline)
-0.80GB 11.91% 30.45%    -0.80GB 11.91%  github.com/ElfAstAhe/url-shortener2/internal/ep/mapper.ShortenBatchResponseFromKeys (inline)
0.32GB  4.69% 25.75%     0.32GB  4.69%  github.com/ElfAstAhe/url-shortener2/internal/ep/dto.NewUserShorten (inline)
-0.26GB  3.86% 29.61%    -0.32GB  4.70%  net/url.(*URL).JoinPath
-0.20GB  2.92% 32.53%    -0.20GB  2.92%  net/url.parse
-0.16GB  2.38% 34.91%    -0.73GB 10.81%  github.com/ElfAstAhe/url-shortener2/internal/ep/mapper.ShortenBatchResponseFromEntity
-0.05GB  0.81% 35.72%    -0.05GB  0.81%  strings.(*Builder).grow
-0.02GB  0.36% 36.09%    -0.06GB  0.84%  path.Join
-0.02GB  0.24% 36.33%    -0.02GB  0.24%  path.(*lazybuf).append (inline)
-0.02GB  0.24% 36.57%    -0.02GB  0.24%  path.(*lazybuf).string (inline)
0     0% 36.57%    -0.73GB 10.81%  github.com/ElfAstAhe/url-shortener2/internal/ep/mapper.BenchmarkShortenBatchResponseFromEntity.func1
0     0% 36.57%    -0.80GB 11.91%  github.com/ElfAstAhe/url-shortener2/internal/ep/mapper.BenchmarkShortenBatchResponseFromKeys.func1
0     0% 36.57%    -0.93GB 13.84%  github.com/ElfAstAhe/url-shortener2/internal/ep/mapper.BenchmarkUserShortensFromModel.func1
0     0% 36.57%    -0.57GB  8.42%  github.com/ElfAstAhe/url-shortener2/internal/utils.BuildNewURI
0     0% 36.57%    -0.05GB  0.81%  net/url.(*URL).String
0     0% 36.57%    -0.20GB  2.92%  net/url.Parse
0     0% 36.57%    -0.03GB  0.48%  path.Clean
0     0% 36.57%    -0.05GB  0.81%  strings.(*Builder).Grow
```