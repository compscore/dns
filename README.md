# DNS

## Check Format

```yaml
- name:
  release:
    org: compscore
    repo: dns
    tag: latest
  target:
  command:
  expectedOutput:
  weight:
```

## Parameters

|    parameter     |       path        |   type   | default  | required | description                                        |
| :--------------: | :---------------: | :------: | :------: | :------: | :------------------------------------------------- |
|      `name`      |      `.name`      | `string` |   `""`   |  `true`  | `name of check (must be unique)`                   |
|      `org`       |  `.release.org`   | `string` |   `""`   |  `true`  | `organization that check repository belongs to`    |
|      `repo`      |  `.release.repo`  | `string` |   `""`   |  `true`  | `repository of the check`                          |
|      `tag`       |  `.release.tag`   | `string` | `latest` | `false`  | `tagged version of check`                          |
|     `target`     |     `.target`     | `string` |   `""`   |  `true`  | `dns server to use to resolve`                     |
|    `command`     |    `.command`     | `string` |   `""`   |  `true`  | `record and domain to resolve; ex: "A google.com"` |
| `expectedOutput` | `.expectedOutput` | `string` |   `""`   | `false`  | `expected output for any of the returned results`  |
|     `weight`     |     `.weight`     |  `int`   |   `0`    |  `true`  | `amount of points a successful check is worth`     |

## Examples

```yaml
- name: google.com-dns
  release:
    org: compscore
    repo: dns
    tag: latest
  target: 8.8.8.8
  command: A google.com
  expectedOutput: 93.184.216.34
  weight: 1
```

```yaml
- name: host_a-dns
  release:
    org: compscore
    repo: dns
    tag: latest
  target: 10.{ .Team }.1.100:53
  command: A webserver.local
  expectedOutput: 10.1.1.1
  weight: 1
```
