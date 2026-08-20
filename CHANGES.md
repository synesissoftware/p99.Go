# p99.Go - Changes <!-- omit in toc -->


## 0.2.1 - 20th August 2026

* enforced Synesis Go import order via **gci** (**.golangci.yml**, **examples/.golangci.yml**);
* normalised example companion docs to the per-program layout (`examples/<name>/main.go` + `examples/<name>.md`);
* version string updated for the 0.2.1 release;


## 0.2.0 - 20th August 2026

* added **Version()** (replacing the **Version** constant), formed by **ver2go.CombineVersion()**;
* updated **ver2go** to 0.2.0-beta1;
* removed retired Go Report Card badge from README;
* version string updated for the 0.2.0 release;


## 0.1.0 - 20th August 2026

* CI modernisation (matrix + lint);
* CI reliability fixes (macOS test linking; golangci-lint config verification disabled in CI);
* boilerplate additions (scripts, markdown docs, project identity);
* added **examples/libver**;
* version string updated for the 0.1.0 release;


## 0.1.0-alpha1 - 6th July 2026

* initial release including `Histogram`;
* Synesis Go project boilerplate;


<!-- ########################### end of file ########################### -->
