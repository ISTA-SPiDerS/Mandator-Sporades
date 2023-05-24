We support two types of data stores 

(1) ```resident key value store```: a simple ```map[string][string]```

(2) ```redis```

[```Redis```](https://redis.io/topics/quickstart) should be installed with ```default options``` in each replica machines.

Both data stores are adopted from [```Rabia```](https://github.com/haochenpan/rabia) 2021 SOSP paper. 