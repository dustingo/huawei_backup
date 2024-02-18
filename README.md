#### 华为云数据库备份工具
> Usage

```shell
huawei cloud database backup

Usage:
  huawei_backup [command]

Available Commands:
  backup      do database backup
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
      --config string   config file
  -h, --help            help for huawei_backup
  -p, --push            push to prometheus pushgateway

Use "huawei_backup [command] --help" for more information about a command.
```
> config

```text
# config file;也可以单独将mongo、rds、redis的配置写入到各自的配置文件中
global:
  region: cn-east-3
  ak:
  sk:
  projectId:
  pushgateway:

mongo:
  instanceId:
    - 63379ff0e1654845a34bff3252adf

rds:
  instanceId:
    - c402c7b5b59d4dc4974985c8ba512

redis:
  format: rdb
  instanceId:
    - 2aec819d-5c1e-4991-853e-863
```