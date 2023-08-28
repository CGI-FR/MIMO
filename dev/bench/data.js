window.BENCHMARK_DATA = {
  "lastUpdate": 1693259349643,
  "repoUrl": "https://github.com/CGI-FR/MIMO",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "44274230+adrienaury@users.noreply.github.com",
            "name": "Adrien Aury",
            "username": "adrienaury"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3cd6bbebabf2fe2091c801aa93e8cc25766aaa20",
          "message": "chore: store benchmark results (#19)\n\n* chore: add benchmark\r\n\r\n* chore: improve benchmark\r\n\r\n* chore: store benchmark results",
          "timestamp": "2023-08-28T23:45:58+02:00",
          "tree_id": "c48757228a68eae90d8a4a6ae539d95c4db4e2bd",
          "url": "https://github.com/CGI-FR/MIMO/commit/3cd6bbebabf2fe2091c801aa93e8cc25766aaa20"
        },
        "date": 1693259349017,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkInMemory - ns/op",
            "value": 86.02,
            "unit": "ns/op",
            "extra": "138071065 times\n2 procs"
          },
          {
            "name": "BenchmarkInMemory - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "138071065 times\n2 procs"
          },
          {
            "name": "BenchmarkInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "138071065 times\n2 procs"
          },
          {
            "name": "BenchmarkOnDisk - ns/op",
            "value": 86.29,
            "unit": "ns/op",
            "extra": "132417262 times\n2 procs"
          },
          {
            "name": "BenchmarkOnDisk - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "132417262 times\n2 procs"
          },
          {
            "name": "BenchmarkOnDisk - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "132417262 times\n2 procs"
          }
        ]
      }
    ]
  }
}