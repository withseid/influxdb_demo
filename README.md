influxdb 调研
===

**按 1w 一个周期找到 pe 的第一条记录(期初)**
```
from(bucket: "company_pe")
//  |> range(start: 2021-03-29T00:00:00Z,stop: 2021-04-03T00:00:00Z)
  |> range(start: -1mo)
  |> filter(fn: (r) => r["_measurement"] == "stat")
  |> filter(fn: (r) => r["_field"] == "pe")
  |> window(every: 1w, period: 5d, offset: 2d)
  |> first()
```

**按 1w 一个周期找到 pe 的最后一条记录(期末)**
```
from(bucket: "company_pe")
//  |> range(start: 2021-03-29T00:00:00Z,stop: 2021-04-03T00:00:00Z)
  |> range(start: -1mo)
  |> filter(fn: (r) => r["_measurement"] == "stat")
  |> filter(fn: (r) => r["_field"] == "pe")
  |> window(every: 1w, period: 5d, offset: 2d)
  |> last()
```

**按 1w 一个周期找到 pe 的平均值**
```
from(bucket: "company_pe")
//   |> range(start: 2021-03-26T00:00:00Z)
  |> range(start: -1mo)
  |> filter(fn: (r) => r["_measurement"] == "stat")
  |> filter(fn: (r) => r["_field"] == "pe")
  |> window(every: 1w, period: 5d, offset: 2d)
  |> mean()
```

**按 1w 一个周期找到 pe 的最大值**
```
from(bucket: "company_pe")
//   |> range(start: 2021-03-26T00:00:00Z)
  |> range(start: -1mo)
  |> filter(fn: (r) => r["_measurement"] == "stat")
  |> filter(fn: (r) => r["_field"] == "pe")
  |> window(every: 1w, period: 5d, offset: 2d)
  |> max()
```

**按 1w 一个周期找到 pe 的最小值**
```
from(bucket: "company_pe")
//   |> range(start: 2021-03-26T00:00:00Z)
  |> range(start: -1mo)
  |> filter(fn: (r) => r["_measurement"] == "stat")
  |> filter(fn: (r) => r["_field"] == "pe")
  |> window(every: 1w, period: 5d, offset: 2d)
  |> min()
```

### demo 需求
- [x] 计算时间范围 PE/PB/PS 的 AVG/ 1STDB / -1STDV
  - [x] 标准差使用 stddev() 函数
  - [x] 平均值使用 mean() 函数
- [x] 计算行情周期的期末值，期初值， 平均值， 最大值， 最小值
  - [x] 期末值使用 last() 函数
  - [x] 期初值使用 first() 函数
  - [x] 平均值使用 mean() 函数
  - [x] 最大值使用 max() 函数
  - [x] 最小值使用 min() 函数
  - [x] 怎么用一个查询语句，获取 5 个统计数据
  - [x] 以周为周期，起始/截止日期必须是有数据的




### 写入性能
开 10 个线程，插入 10000 个公司 10 年期间的数据，也就是 10000 * 2600 = 2600 万条数据，花了 10 分钟
![avatar](https://i.bmp.ovh/imgs/2021/04/1fe4c266c661152b.png)
