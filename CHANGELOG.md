# [v0.5.0](https://github.com/dtan4/k8stail/releases/tag/v0.5.0) (2017-04-25)

## Features

- Watch Kubernetes events to detect Pod lifecycle correctly [#26](https://github.com/dtan4/k8stail/pull/26)
- Add more short flags [#25](https://github.com/dtan4/k8stail/pull/25) (thanks @atombender)

# [v0.4.0](https://github.com/dtan4/k8stail/releases/tag/v0.4.0) (2017-04-11)

## Features

- Use default namespace set in kubecfg [#21](https://github.com/dtan4/k8stail/pull/21)
- Add `--no-halt` flag [#18](https://github.com/dtan4/k8stail/pull/18)

## Fixed

- Detect container recreation [#20](https://github.com/dtan4/k8stail/pull/20)

# [v0.3.0](https://github.com/dtan4/k8stail/releases/tag/v0.3.0) (2016-12-12)

## Backward incompatible changes

- Deprecate `-flag` style flag, use `--flag` [#11](https://github.com/dtan4/k8stail/pull/11)

## Features

- Support context switch by `--context` flag [#13](https://github.com/dtan4/k8stail/pull/13) (Thanks @apstndb)

# [v0.2.1](https://github.com/dtan4/k8stail/releases/tag/v0.2.1) (2016-11-16)

Rebuilt binaries to be statically-linked.

# [v0.2.0](https://github.com/dtan4/k8stail/releases/tag/v0.2.0) (2016-11-16)

## Features

- Stream logs of all containers in pod [#5](https://github.com/dtan4/k8stail/pull/5)
- Get kubeconfig path from KUBECONFIG [#4](https://github.com/dtan4/k8stail/pull/4)

# [v0.1.0](https://github.com/dtan4/k8stail/releases/tag/v0.1.0) (2016-11-15)

Initial release.
